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
}

type ChainCache[T any] struct {
	stores []struct {
		name   string
		store  CacheStore[T]
		getTTL TTLStrategy
	}

	group    *singleflight.Group
	updating sync.Map
}

func NewChainCache[T any]() *ChainCache[T] {
	return &ChainCache[T]{
		group: &singleflight.Group{},
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
			if i > 0 {
				go c.handleCacheUpdate(ctx, key, value, i-1)
			}
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

	return c.handleCacheUpdate(ctx, key, value, len(c.stores)-2)
}

func (c *ChainCache[T]) Update(ctx context.Context, key string) error {
	value, err := c.stores[len(c.stores)-1].store.Get(ctx, key)
	if err != nil {
		return err
	}

	return c.handleCacheUpdate(ctx, key, value, len(c.stores)-2)
}

func (c *ChainCache[T]) handleCacheUpdate(ctx context.Context, key string, value T, level int) error {
	if level < 0 {
		return fmt.Errorf("level must be greater than 0, got %d", level)
	}

	if level >= len(c.stores)-1 {
		return errors.New("this is only for backfill, please use Set or Update to update the datasource")
	}

	// use the last write wins strategy to update the cache
	var updateChanAny any
	var updateChan chan T

	updateChanAny, _ = c.updating.LoadOrStore(key, make(chan T, 10000))
	updateChan = updateChanAny.(chan T)
	updateChan <- value

	defer c.updating.Delete(key)

	var err error
	for {
		v, ok := <-updateChan
		if !ok {
			break
		}
		for len(updateChan) > 0 {
			v = <-updateChan
		}

		err = c.updateCache(ctx, key, v, level)
	}

	return err
}

func (c *ChainCache[T]) updateCache(ctx context.Context, key string, value T, level int) error {
	for i := level; i >= 0; i-- {
		err := c.stores[i].store.Set(ctx, key, value, c.stores[i].getTTL.GetTTL(key, i))
		if err != nil {
			return err
		}
	}
	return nil
}
