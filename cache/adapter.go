package cache

import "time"

// Adapter defines a group of methods to manipulate cache keys
type Adapter interface {
	Set(key string, value interface{}, ttl time.Duration) (string, error)
	SetOrFail(key string, value interface{}, ttl time.Duration) (bool, error)
	Get(key string) (string, error)
	Delete(key string) error
}
