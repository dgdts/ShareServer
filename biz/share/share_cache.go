package share

import (
	"time"

	"github.com/dgdts/ShareServer/pkg/cache"
	"github.com/dgdts/UniversalServer/pkg/redis"
)

var (
	shareCache = cache.NewChainCache[[]byte]()
)

func InitShareCache() {
	shareCache.AddStore("memory", cache.NewMemoryStore[[]byte](), &cache.IncrementalTTL{
		BaseTTL:   time.Minute * 10,
		Increment: time.Minute * 5,
	})
	shareCache.AddStore("redis", cache.NewRedisStore[[]byte](redis.GetConnection(), "share_cache"), &cache.RandomTTL{
		BaseTTL: time.Minute * 20,
		Jitter:  0.5,
	})
	shareCache.AddStore("mongo", NewShareMongoStore(), nil)
}
