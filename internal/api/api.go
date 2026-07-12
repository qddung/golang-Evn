package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	_ "github.com/homework/lab/docs"
	"github.com/homework/lab/internal/config"
	"github.com/homework/lab/internal/handler"
	"github.com/homework/lab/internal/repository"
	"github.com/homework/lab/internal/service"
	"github.com/homework/lab/pkg/helpers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Engine interface for app engine
type Engine interface {
	Run() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// engine struct for app engine
type engine struct {
	app   *gin.Engine
	cfg   *config.Config
	redis *redis.Client
}

// NewEngine creates a new engine instance
func NewEngine(cfg *config.Config, redisClient ...*redis.Client) Engine {
	var rClient *redis.Client
	if len(redisClient) > 0 && redisClient[0] != nil {
		rClient = redisClient[0]
	} else {
		rClient = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
	}

	api := &engine{
		app:   gin.New(),
		cfg:   cfg,
		redis: rClient,
	}

	api.initRoutes(cfg)
	return api
}

// config Run starts the app engine
func (e *engine) Run() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// config ServeHTTP serves the app engine
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

type handlers struct {
	healthCheck handler.HealthCheck
	shorten     handler.ShorternUrl
}

func (e *engine) InitHandlers(cfg *config.Config) handlers {
	serviceName := cfg.ServiceName
	instanceID := cfg.InstanceID
	// create handler
	healthCheckRepository := repository.NewPing(e.redis)
	healthCheckService := service.NewHealthCheck(serviceName, instanceID, healthCheckRepository)
	healthCheckHandler := handler.NewHealthCheck(healthCheckService)

	// create shrotten url handler
	urlStorage := repository.NewURLStorage(e.redis)
	shortenService := service.NewShorternUrl(urlStorage, helpers.NewKeyGenerator())
	shortenURLHandler := handler.NewShortenURL(shortenService)
	return handlers{healthCheckHandler, shortenURLHandler}
}

// initRoutes initializes the routes for the app engine
func (e *engine) initRoutes(cfg *config.Config) {
	allHandlers := e.InitHandlers(cfg)

	e.app.GET("/ping", allHandlers.healthCheck.Ping)
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1Routes := e.app.Group("/api/v1")
	{
		v1Routes.POST("/shorten", allHandlers.shorten.ShortenUrl)
	}
}
