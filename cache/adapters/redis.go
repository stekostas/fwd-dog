package adapters

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

// RedisAdapter provides a concrete implementation of the `cache.Adapter` contract in order to provide an interface
// between the application and a Redis client.
type RedisAdapter struct {
	Client *redis.Client
}

// NewRedisAdapter instantiates a new RedisAdapter and returns a pointer to the object.
func NewRedisAdapter(client *redis.Client) *RedisAdapter {
	adapter := &RedisAdapter{
		Client: client,
	}
	return adapter
}

// Get returns the value of the key found or an error when it does not exist.
func (ra *RedisAdapter) Get(key string) (string, error) {
	return ra.Client.Get(key).Result()
}

// Set saves the provided value to the specified key for the TTL or returns an error on failure.
func (ra *RedisAdapter) Set(key string, value interface{}, ttl time.Duration) (string, error) {
	return ra.Client.Set(key, value, ttl).Result()
}

// SetOrFail tries to set the provided value to the specified key with the TTL or returns false when the key already
// exists in Redis.
func (ra *RedisAdapter) SetOrFail(key string, value interface{}, ttl time.Duration) (bool, error) {
	return ra.Client.SetNX(key, value, ttl).Result()
}

// Delete deletes the specified key from Redis.
func (ra *RedisAdapter) Delete(key string) error {
	deleted, err := ra.Client.Del(key).Result()

	if err != nil {
		return err
	}

	if deleted != 1 {
		return fmt.Errorf("could not delete key '%s' from cache, deleted %d", key, deleted)
	}

	return nil
}
