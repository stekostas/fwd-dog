// +build integration_tests

package adapters

import (
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type RedisAdapterSuite struct {
	suite.Suite
	adapter *RedisAdapter
}

func TestRedisAdapterSuite(t *testing.T) {
	suite.Run(t, new(RedisAdapterSuite))
}

func (s *RedisAdapterSuite) SetupTest() {
	options := &redis.Options{
		Addr: "redis:6379",
	}
	client := redis.NewClient(options)
	client.FlushAll()
	s.adapter = NewRedisAdapter(client)
}

func (s *RedisAdapterSuite) TestNewRedisAdapter() {
	s.Equal("redis:6379", s.adapter.Client.Options().Addr)
}

func (s *RedisAdapterSuite) TestSuccessfulGet() {
	// Set new key
	s.adapter.Client.Set("test", "test", 0)

	key, err := s.adapter.Get("test")

	s.Equal("test", key)
	s.NoError(err)
}

func (s *RedisAdapterSuite) TestUnsuccessfulGet() {
	key, err := s.adapter.Get("test")

	s.Empty(key)
	s.NotNil(err)
}

func (s *RedisAdapterSuite) TestSuccessfulSet() {
	val, err := s.adapter.Set("test", "test", 0)

	s.Equal("OK", val)
	s.Nil(err)
}

func (s *RedisAdapterSuite) TestSuccessfulSetWithTtl() {
	val, err := s.adapter.Set("test", "test", time.Second)

	s.Equal("OK", val)
	s.Nil(err)
}

func (s *RedisAdapterSuite) TestSuccessfulSetWithDuplicate() {
	val, err := s.adapter.Set("test", "test", 0)

	s.Equal("OK", val)
	s.Nil(err)

	val, err = s.adapter.Set("test", "test2", 0)

	s.Equal("OK", val)
	s.Nil(err)

	val, err = s.adapter.Client.Get("test").Result()

	s.Equal("test2", val)
	s.Nil(err)
}

func (s *RedisAdapterSuite) TestSuccessfulSetNX() {
	ok, err := s.adapter.SetOrFail("test", "test", 0)

	s.True(ok)
	s.Nil(err)
}

func (s *RedisAdapterSuite) TestSuccessfulSetNXWithTtl() {
	ok, err := s.adapter.SetOrFail("test", "test", time.Second)

	s.True(ok)
	s.Nil(err)
}

func (s *RedisAdapterSuite) TestUnsuccessfulSetNXWithDuplicate() {
	ok, err := s.adapter.SetOrFail("test", "test", 0)

	s.True(ok)
	s.Nil(err)

	ok, err = s.adapter.SetOrFail("test", "test2", 0)

	s.False(ok)
	s.Nil(err)

	val, err := s.adapter.Client.Get("test").Result()

	s.Equal("test", val)
	s.Nil(err)
}

func (s *RedisAdapterSuite) TestSuccessfulDelete() {
	val, err := s.adapter.Client.Set("test", "test", 0).Result()

	s.Equal("OK", val)
	s.Nil(err)

	val, err = s.adapter.Client.Get("test").Result()

	s.Equal("test", val)
	s.Nil(err)

	err = s.adapter.Delete("test")

	s.Nil(err)

	val, err = s.adapter.Client.Get("test").Result()

	s.Empty(val)
	s.NotNil(err)
}

func (s *RedisAdapterSuite) TestUnsuccessfulDelete() {
	err := s.adapter.Delete("test")

	s.NotNil(err)
}
