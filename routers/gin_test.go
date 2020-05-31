package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	RedisHost  = "test"
	TtlOptions = map[time.Duration]string{}
)

type GinTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func TestRedisAdapterSuite(t *testing.T) {
	suite.Run(t, new(GinTestSuite))
}

func (s *GinTestSuite) SetupTest() {
	s.router = SetupGinRouter("..", RedisHost, TtlOptions)
}

func (s *GinTestSuite) TestGetHomepageRoute() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.NotEmpty(w.Body.String())
}

func (s *GinTestSuite) TestPostHomepageRoute() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	s.NotEmpty(w.Body.String())
}

func (s *GinTestSuite) TestAssetRoutes() {
	w := httptest.NewRecorder()
	files := []string{
		"/assets/css/main.css",
		"/assets/js/copy-link.js",
		"/assets/js/form.js",
	}

	for _, file := range files {
		req, _ := http.NewRequest(http.MethodGet, file, nil)
		s.router.ServeHTTP(w, req)
		s.Equal(http.StatusOK, w.Code)
		s.NotEmpty(w.Body.String())
	}
}

func (s *GinTestSuite) TestPagesRoutes() {
	w := httptest.NewRecorder()
	pages := []string{
		"/_about",
		"/_credits",
	}

	for _, page := range pages {
		req, _ := http.NewRequest(http.MethodGet, page, nil)
		s.router.ServeHTTP(w, req)
		s.Equal(http.StatusOK, w.Code)
		s.NotEmpty(w.Body.String())
	}
}

func (s *GinTestSuite) TestInternalServerError() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/missing", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.NotEmpty(w.Body.String())
}
