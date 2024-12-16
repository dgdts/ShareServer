package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type CacheStore[T any] interface {
	Get(ctx context.Context, key string) (T, error)
	Set(ctx context.Context, key string, value T, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type UpdateStrategy int

const (
	DeleteOnUpdate UpdateStrategy = iota
	RefreshOnUpdate
)

type ChainCache[T any] struct {
	stores []struct {
		name   string
		store  CacheStore[T]
		getTTL TTLStrategy
	}
	updateStrategy UpdateStrategy
	group          *singleflight.Group
	updating       sync.Map
}

func NewChainCache[T any](updateStrategy UpdateStrategy) *ChainCache[T] {
	return &ChainCache[T]{
		updateStrategy: updateStrategy,
		group:          &singleflight.Group{},
	}
}

func (c *ChainCache[T]) AddStore(name string, store CacheStore[T], ttlStrategy TTLStrategy) *ChainCache[T] {
	c.stores = append(c.stores, struct {
		name   string
		store  CacheStore[T]
		getTTL TTLStrategy
	}{name, store, ttlStrategy})
	return c
}

func (c *ChainCache[T]) Get(ctx context.Context, key string) (T, error) {
	value, err, _ := c.group.Do(key, func() (interface{}, error) {
		return c.doGet(ctx, key)
	})
	if err != nil {
		var zero T
		return zero, err
	}
	return value.(T), nil
}

func (c *ChainCache[T]) doGet(ctx context.Context, key string) (T, error) {
	var lastErr error
	for i, s := range c.stores {
		value, err := s.store.Get(ctx, key)
		if err == nil {
			backfillCtx := context.Background()
			go func(level int, k string, v T) {
				c.group.Do(fmt.Sprintf("backfill-%s-%d", k, level), func() (interface{}, error) {
					return nil, c.backfill(backfillCtx, level, k, v)
				})
			}(i, key, value)
			return value, nil
		}
		lastErr = err
	}
	var zero T
	return zero, lastErr
}

func (c *ChainCache[T]) Set(ctx context.Context, key string, value T) error {
	lastStore := c.stores[len(c.stores)-1]
	if err := lastStore.store.Set(ctx, key, value, 0); err != nil {
		return err
	}

	return c.handleCacheUpdate(ctx, key, value)
}

func (c *ChainCache[T]) Update(ctx context.Context, key string) error {
	if _, updating := c.updating.LoadOrStore(key, true); updating {
		return errors.New("update in progress")
	}

	defer c.updating.Delete(key)

	value, err := c.stores[len(c.stores)-1].store.Get(ctx, key)
	if err != nil {
		return err
	}

	return c.handleCacheUpdate(ctx, key, value)
}

func (c *ChainCache[T]) handleCacheUpdate(ctx context.Context, key string, value T) error {
	switch c.updateStrategy {
	case DeleteOnUpdate:
		for i := 0; i < len(c.stores)-1; i++ {
			if err := c.stores[i].store.Delete(ctx, key); err != nil {
				return err
			}
		}
	case RefreshOnUpdate:
		for i := 0; i < len(c.stores)-1; i++ {
			if err := c.stores[i].store.Set(ctx, key, value, c.stores[i].getTTL.GetTTL(key, i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ChainCache[T]) backfill(ctx context.Context, level int, key string, value T) error {
	for i := 0; i < level; i++ {
		ttl := c.stores[i].getTTL.GetTTL(key, i)
		if ttl > 0 {
			if err := c.stores[i].store.Set(ctx, key, value, ttl); err != nil {
				return err
			}
		}
	}
	return nil
}
