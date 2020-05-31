package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stekostas/fwd-dog/cache"
	"github.com/stekostas/fwd-dog/handlers"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHomepageHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	cacheAdapter := cache.NewRedisAdapter(RedisAddress, DefaultTtl)
	renderer := handlers.NewRenderer(StaticFilesRoot, AssetsRoot)
	context := handlers.NewContext(renderer, cacheAdapter, TtlOptions)
	handler := handlers.NewHomepageHandler(context)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHomepageFailedSubmission(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	cacheAdapter := cache.NewRedisAdapter(RedisAddress, DefaultTtl)
	renderer := handlers.NewRenderer(StaticFilesRoot, AssetsRoot)
	context := handlers.NewContext(renderer, cacheAdapter, TtlOptions)
	handler := handlers.NewHomepageHandler(context)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestHomepageSuccessfulSubmission(t *testing.T) {
	data := url.Values{}
	data.Add("url", "https://github.com")
	data.Add("ttl", "300")
	req, err := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	cacheAdapter := cache.NewRedisAdapter(RedisAddress, DefaultTtl)
	renderer := handlers.NewRenderer(StaticFilesRoot, AssetsRoot)
	context := handlers.NewContext(renderer, cacheAdapter, TtlOptions)
	handler := handlers.NewHomepageHandler(context)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestMissingLink(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	cacheAdapter := cache.NewRedisAdapter(RedisAddress, DefaultTtl)
	renderer := handlers.NewRenderer(StaticFilesRoot, AssetsRoot)
	context := handlers.NewContext(renderer, cacheAdapter, TtlOptions)
	handler := handlers.NewLinkRedirectHandler(context)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestCorrectLink(t *testing.T) {
	// Generate key
	data := url.Values{}
	data.Add("url", "https://github.com")
	data.Add("ttl", "300")
	req, err := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	cacheAdapter := cache.NewRedisAdapter(RedisAddress, DefaultTtl)
	renderer := handlers.NewRenderer(StaticFilesRoot, AssetsRoot)
	context := handlers.NewContext(renderer, cacheAdapter, TtlOptions)
	handler := handlers.NewHomepageHandler(context)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	key := rr.Result().Header.Get("X-Fwd-Key")

	if key == "" {
		t.Error("handler does not contain a valid key in the 'X-Fwd-Key' header")
	}

	req2, err2 := http.NewRequest("GET", fmt.Sprintf("/%s", key), nil)

	if err2 != nil {
		t.Fatal(err2)
	}

	rr2 := httptest.NewRecorder()
	r := mux.NewRouter()
	linkHandler := handlers.NewLinkRedirectHandler(context)
	r.Handle("/{key}", linkHandler)
	r.ServeHTTP(rr2, req2)

	// Check the status code is what we expect
	if status := rr2.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}
}
