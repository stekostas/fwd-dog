package main

import (
	"github.com/gorilla/mux"
	"github.com/stekostas/fwd-dog/cache"
	"github.com/stekostas/fwd-dog/handlers"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	StaticFilesRoot = "public"
	AssetsRoot      = "public/assets"
)

var (
	ServerHost   = os.Getenv("APP_HOST")
	DefaultTtl   = time.Minute * 5
	RedisAddress = os.Getenv("REDIS_ADDRESS")
	TtlOptions   = []time.Duration{
		time.Minute * 5,
		time.Minute * 30,
		time.Hour,
		time.Hour * 12,
		time.Hour * 24,
	}
)

func main() {
	renderer := handlers.NewRenderer(StaticFilesRoot, AssetsRoot)
	cacheAdapter := cache.NewRedisAdapter(RedisAddress, DefaultTtl)
	context := handlers.NewContext(renderer, cacheAdapter, TtlOptions)

	handler := mux.NewRouter()
	fs := http.FileServer(http.Dir(AssetsRoot))

	handler.NotFoundHandler = handlers.NewNotFoundHandler(context)
	handler.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))
	handler.Handle("/", handlers.NewHomepageHandler(context))
	handler.Handle("/{key}", handlers.NewLinkRedirectHandler(context))

	log.Printf("[INFO] Starting web server; listening on '%s'\n", ServerHost)
	log.Fatal(http.ListenAndServe(ServerHost, handler))
}
