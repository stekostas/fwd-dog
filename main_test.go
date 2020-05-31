// +build integration_tests

package main

import (
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

const AppUrl = "http://localhost:3000"

type MainTestSuite struct {
	suite.Suite
	redisClient *redis.Client
}

func TestMainSuite(t *testing.T) {
	s := new(MainTestSuite)
	s.redisClient = redis.NewClient(&redis.Options{Addr: RedisHost})
	go main()
	defer s.redisClient.FlushAll()
	suite.Run(t, s)
}

func (s *MainTestSuite) SetupTest() {
	s.redisClient.FlushAll()
}

func (s *MainTestSuite) TestServer() {
	res, err := http.Get(AppUrl + "/")

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)
}

func (s *MainTestSuite) TestNotFoundKey() {
	res, err := http.Get(AppUrl + "/missing")

	s.Nil(err)
	s.Equal(http.StatusNotFound, res.StatusCode)
}

func (s *MainTestSuite) TestNotFoundPage() {
	res, err := http.Get(AppUrl + "/_missing")

	s.Nil(err)
	s.Equal(http.StatusNotFound, res.StatusCode)
}

func (s *MainTestSuite) TestGenerateLinkNoOptions() {
	data := url.Values{}
	data.Add("targetUrl", AppUrl)
	data.Add("ttl", "300")

	res, err := http.Post(AppUrl+"/", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	s.Nil(err)
	s.Equal(http.StatusCreated, res.StatusCode)

	key := res.Header.Get("X-Fwd-Key")
	s.NotEmpty(key)

	res, err = http.Get(AppUrl + "/" + key)

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)
}

func (s *MainTestSuite) TestGenerateLinkSingleUse() {
	data := url.Values{}
	data.Add("targetUrl", AppUrl)
	data.Add("ttl", "300")
	data.Add("single-use", "on")

	res, err := http.Post(AppUrl+"/", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	s.Nil(err)
	s.Equal(http.StatusCreated, res.StatusCode)

	key := res.Header.Get("X-Fwd-Key")
	s.NotEmpty(key)
	s.True(strings.HasPrefix(key, "."))

	res, err = http.Get(AppUrl + "/" + key)

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)

	res, err = http.Get(AppUrl + "/" + key)

	s.Nil(err)
	s.Equal(http.StatusNotFound, res.StatusCode)
}

func (s *MainTestSuite) TestGenerateLinkWithPassword() {
	data := url.Values{}
	data.Add("targetUrl", AppUrl)
	data.Add("ttl", "300")
	data.Add("password-protected", "on")
	data.Add("password", "hunter2")

	res, err := http.Post(AppUrl+"/", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	s.Nil(err)
	s.Equal(http.StatusCreated, res.StatusCode)

	key := res.Header.Get("X-Fwd-Key")
	s.NotEmpty(key)

	res, err = http.Get(AppUrl + "/" + key)

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)

	data = url.Values{}
	data.Add("password", "hunter2")

	res, err = http.Post(AppUrl+"/"+key, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)
}

func (s *MainTestSuite) TestGenerateLinkWithPasswordCannotUnlock() {
	data := url.Values{}
	data.Add("targetUrl", AppUrl)
	data.Add("ttl", "300")
	data.Add("password-protected", "on")
	data.Add("password", "hunter2")

	res, err := http.Post(AppUrl+"/", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	s.Nil(err)
	s.Equal(http.StatusCreated, res.StatusCode)

	key := res.Header.Get("X-Fwd-Key")
	s.NotEmpty(key)

	res, err = http.Get(AppUrl + "/" + key)

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)

	data = url.Values{}
	data.Add("password", "wrong")

	res, err = http.Post(AppUrl+"/"+key, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	s.Nil(err)
	s.Equal(http.StatusUnauthorized, res.StatusCode)
}

func (s *MainTestSuite) TestGenerateLinkWithPasswordCannotUnlockBadRequest() {
	data := url.Values{}
	data.Add("targetUrl", AppUrl)
	data.Add("ttl", "300")
	data.Add("password-protected", "on")
	data.Add("password", "hunter2")

	res, err := http.Post(AppUrl+"/", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	s.Nil(err)
	s.Equal(http.StatusCreated, res.StatusCode)

	key := res.Header.Get("X-Fwd-Key")
	s.NotEmpty(key)

	res, err = http.Get(AppUrl + "/" + key)

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)

	data = url.Values{}

	res, err = http.Post(AppUrl+"/"+key, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	s.Nil(err)
	s.Equal(http.StatusBadRequest, res.StatusCode)
}

func (s *MainTestSuite) TestGenerateLinkPasswordProtectedNoPassword() {
	data := url.Values{}
	data.Add("targetUrl", AppUrl)
	data.Add("ttl", "300")
	data.Add("password-protected", "on")

	res, err := http.Post(AppUrl+"/", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	s.Nil(err)
	s.Equal(http.StatusBadRequest, res.StatusCode)
}

func (s *MainTestSuite) TestGenerateLinkAllOptions() {
	data := url.Values{}
	data.Add("targetUrl", AppUrl)
	data.Add("ttl", "300")
	data.Add("single-use", "on")
	data.Add("password-protected", "on")
	data.Add("password", "hunter2")

	res, err := http.Post(AppUrl+"/", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	s.Nil(err)
	s.Equal(http.StatusCreated, res.StatusCode)

	key := res.Header.Get("X-Fwd-Key")
	s.NotEmpty(key)
	s.True(strings.HasPrefix(key, "."))

	res, err = http.Get(AppUrl + "/" + key)

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)

	data = url.Values{}
	data.Add("password", "hunter2")

	res, err = http.Post(AppUrl+"/"+key, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)

	res, err = http.Get(AppUrl + "/" + key)

	s.Nil(err)
	s.Equal(http.StatusNotFound, res.StatusCode)
}
