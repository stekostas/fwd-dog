package cache

import (
	"github.com/go-redis/redis/v7"
	"time"
)

type RedisAdapter struct {
	Client     *redis.Client
	DefaultTtl time.Duration
}

func NewRedisAdapter(address string, defaultTtl time.Duration) *RedisAdapter {
	adapter := &RedisAdapter{Client: redis.NewClient(&redis.Options{Addr: address})}
	adapter.DefaultTtl = defaultTtl
	return adapter
}

func (ra *RedisAdapter) Get(key string) (string, error) {
	return ra.Client.Get(key).Result()
}

func (ra *RedisAdapter) Set(key string, value interface{}, ttl time.Duration) (string, error) {
	return ra.Client.Set(key, value, ttl).Result()
}

func (ra *RedisAdapter) SetOrFail(key string, value interface{}, ttl time.Duration) (bool, error) {
	return ra.Client.SetNX(key, value, ttl).Result()
}
