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
	TtlOptions   = map[time.Duration]string{
		time.Minute * 5:  "5 minutes",
		time.Minute * 15: "15 minutes",
		time.Minute * 30: "30 minutes",
		time.Hour:        "1 hour",
		time.Hour * 3:    "3 hours",
		time.Hour * 6:    "6 hours",
		time.Hour * 12:   "12 hours",
		time.Hour * 24:   "1 day",
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
