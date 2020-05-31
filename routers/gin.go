package routers

import (
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/stekostas/fwd-dog/cache/adapters"
	"github.com/stekostas/fwd-dog/handlers"
	"github.com/stekostas/fwd-dog/services"
	"time"
)

func SetupGinRouter(root string, redisHost string, ttlOptions map[time.Duration]string) *gin.Engine {
	router := gin.New()

	// Set up app dependencies
	redisOptions := &redis.Options{Addr: redisHost}
	redisClient := redis.NewClient(redisOptions)

	redisAdapter := adapters.NewRedisAdapter(redisClient)
	linkGenerator := services.NewLinkGenerator(redisAdapter)
	homepageHandler := handlers.NewHomepageHandler(ttlOptions, linkGenerator)
	redirectHandler := handlers.NewRedirectHandler(redisAdapter)

	// Define the required middleware
	router.Use(gin.Logger())
	router.Use(nice.Recovery(handlers.RecoveryHandler))

	router.Static("/assets/css", root+"/public/css")
	router.Static("/assets/js", root+"/public/js")
	router.LoadHTMLGlob(root + "/templates/*")

	// Define the HTTP 404 or redirect handler
	router.NoRoute(redirectHandler.Redirect)

	// Register the app routes
	router.GET("/", homepageHandler.Get)
	router.POST("/", homepageHandler.Post)
	router.GET("/_about", handlers.AboutPageHandler)
	router.GET("/_credits", handlers.CreditsPageHandler)

	return router
}
