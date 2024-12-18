package cache

import (
	"context"
	"errors"
	"time"

	"github.com/dgdts/UniversalServer/pkg/utils"
)

var _ CacheStore[any] = (*MemoryStore[any])(nil)

type MemoryStore[T any] struct {
	cache *utils.TTLMap[string, T]
}

func NewMemoryStore[T any]() *MemoryStore[T] {
	return &MemoryStore[T]{
		cache: utils.New[string, T](),
	}
}

func (m *MemoryStore[T]) Get(ctx context.Context, key string) (T, error) {
	if v, ok := m.cache.Get(key); ok {
		return v, nil
	}
	var zero T
	return zero, errors.New("key not found")
}

func (m *MemoryStore[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	m.cache.Set(key, value, ttl)
	return nil
}
