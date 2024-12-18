package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ CacheStore[any] = (*RedisStore[any])(nil)

type RedisStore[T any] struct {
	redis  *redis.Client
	prefix string
}

func NewRedisStore[T any](redis *redis.Client, prefix string) *RedisStore[T] {
	return &RedisStore[T]{
		redis:  redis,
		prefix: prefix,
	}
}

func (r *RedisStore[T]) Get(ctx context.Context, key string) (T, error) {
	var zero T
	redisKey := fmt.Sprintf("%s:%s", r.prefix, key)
	data, err := r.redis.Get(ctx, redisKey).Result()
	if err != nil {
		return zero, err
	}
	var value T
	// data will be deep copied to the []byte
	// might consider use unsafe.Pointer to optimize
	err = json.Unmarshal([]byte(data), &value)
	return value, err
}

func (r *RedisStore[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	redisKey := fmt.Sprintf("%s:%s", r.prefix, key)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.redis.Set(ctx, redisKey, data, ttl).Err()
}
