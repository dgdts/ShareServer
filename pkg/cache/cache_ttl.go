package cache

import (
	"math/rand/v2"
	"time"
)

type TTLStrategy interface {
	GetTTL(key string, level int) time.Duration
}

type IncrementalTTL struct {
	baseTTL   time.Duration
	increment time.Duration
}

func (t *IncrementalTTL) GetTTL(key string, level int) time.Duration {
	return t.baseTTL + time.Duration(level)*t.increment
}

type RandomTTL struct {
	baseTTL time.Duration
	jitter  float64
}

func (t *RandomTTL) GetTTL(key string, level int) time.Duration {
	jitterDuration := time.Duration(float64(t.baseTTL) * (1 + (rand.Float64()*2-1)*t.jitter))
	return jitterDuration
}
