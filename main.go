package main

import (
	"github.com/stekostas/fwd-dog/routers"
	"log"
	"os"
	"time"
)

var (
	AppHost      = os.Getenv("APP_HOST")
	RedisAddress = os.Getenv("REDIS_ADDRESS")
	TtlOptions   = map[time.Duration]string{
		time.Minute * 5:    "5 minutes",
		time.Minute * 15:   "15 minutes",
		time.Minute * 30:   "30 minutes",
		time.Hour:          "1 hour",
		time.Hour * 3:      "3 hours",
		time.Hour * 6:      "6 hours",
		time.Hour * 12:     "12 hours",
		time.Hour * 24:     "1 day",
		time.Hour * 24 * 3: "3 days",
		time.Hour * 24 * 7: "1 week",
	}
)

func main() {
	router := routers.SetupGinRouter(".", RedisAddress, TtlOptions)
	log.Fatalln(router.Run(AppHost))
}
