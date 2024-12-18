package cache

import (
	"math/rand/v2"
	"time"
)

type TTLStrategy interface {
	GetTTL(key string, level int) time.Duration
}

type IncrementalTTL struct {
	BaseTTL   time.Duration
	Increment time.Duration
}

func (t *IncrementalTTL) GetTTL(key string, level int) time.Duration {
	return t.BaseTTL + time.Duration(level)*t.Increment
}

type RandomTTL struct {
	BaseTTL time.Duration
	Jitter  float64
}

func (t *RandomTTL) GetTTL(key string, level int) time.Duration {
	jitterDuration := time.Duration(float64(t.BaseTTL) * (1 + (rand.Float64()*2-1)*t.Jitter))
	return jitterDuration
}
