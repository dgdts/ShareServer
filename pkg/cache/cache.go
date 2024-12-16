package cache

import (
	"context"
	"time"
)

type CacheStore[T any] interface {
	Get(ctx context.Context, key string) (T, error)
	Set(ctx context.Context, key string, value T, ttl time.Duration) error
}

type TTLStrategy interface {
	GetTTL(key string, level int) time.Duration
}

type ChainCache[T any] struct {
	stores []struct {
		name  string
		store CacheStore[T]
	}
}

func NewChainCache[T any]() *ChainCache[T] {
	return &ChainCache[T]{}
}

func (c *ChainCache[T]) AddStore(name string, store CacheStore[T]) *ChainCache[T] {
	c.stores = append(c.stores, struct {
		name  string
		store CacheStore[T]
	}{name, store})
	return c
}

func (c *ChainCache[T]) Get(ctx context.Context, key string) (T, error) {
	var lastErr error
	for i, s := range c.stores {
		value, err := s.store.Get(ctx, key)
		if err == nil {
			c.backfill(ctx, i, key, value)
			return value, nil
		}
		lastErr = err
	}
	var zero T
	return zero, lastErr
}

func (c *ChainCache[T]) backfill(ctx context.Context, level int, key string, value T) {
	for i := 0; i < level; i++ {
		s := c.stores[i]
		s.store.Set(ctx, key, value, 0)
	}
}
