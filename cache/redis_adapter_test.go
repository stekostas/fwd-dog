package cache

import (
	"testing"
	"time"
)

func TestNewRedisAdapter(t *testing.T) {
	address := "addr"
	defaultTtl := time.Minute
	adapter := NewRedisAdapter(address, defaultTtl)

	if adapter.DefaultTtl != defaultTtl {
		t.Fatalf("default TTL option is invalid: expected %v, got %v", defaultTtl, adapter.DefaultTtl)
	}

	clientAddress := adapter.Client.Options().Addr

	if clientAddress != address {
		t.Fatalf("client address does not match the one specified: expected %v, got %v", address, clientAddress)
	}
}
