package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/homework/lab/docs"
	_ "github.com/homework/lab/docs"
	"github.com/homework/lab/internal/config"
	"github.com/homework/lab/internal/connection"
	health_check_handler "github.com/homework/lab/internal/handler/health_check"
	"github.com/homework/lab/internal/handler/shorten"
	user_handler "github.com/homework/lab/internal/handler/user"
	health_check_repository "github.com/homework/lab/internal/repository/health_check"
	url_repository "github.com/homework/lab/internal/repository/shorten"
	userRepository "github.com/homework/lab/internal/repository/user"
	health_check_service "github.com/homework/lab/internal/service/health_check"
	shorten_service "github.com/homework/lab/internal/service/shorten"
	user_service "github.com/homework/lab/internal/service/user"
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
	app       *gin.Engine
	cfg       *config.Config
	connector connection.DBConnector
}

// NewEngine creates a new engine instance
func NewEngine(cfg *config.Config, conn connection.DBConnector) Engine {
	api := &engine{
		app:       gin.New(),
		cfg:       cfg,
		connector: conn,
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
	healthCheck health_check_handler.HealthCheck
	shorten     shorten.ShorternUrl
	user        user_handler.UserHandler
	config      *config.Config
}

func (e *engine) InitHandlers(cfg *config.Config) handlers {
	serviceName := cfg.ServiceName
	instanceID := cfg.InstanceID
	redisClient := e.connector.GetRedisClient()
	sqlDB := e.connector.GetSqlDB()
	// create handler
	healthCheckRepository := health_check_repository.NewPing(redisClient)
	healthCheckService := health_check_service.NewHealthCheck(serviceName, instanceID, healthCheckRepository)
	healthCheckHandler := health_check_handler.NewHealthCheck(healthCheckService)

	// create shorten url handler
	urlStorage := url_repository.NewURLStorage(redisClient)
	shortenService := shorten_service.NewShorternUrl(urlStorage, helpers.NewKeyGenerator())
	shortenURLHandler := shorten.NewShortenURL(shortenService)

	// create user handler
	userRepository := userRepository.NewUserRepository(sqlDB)
	hasher := helpers.NewHasher()
	userService := user_service.NewUserService(userRepository, hasher)
	userHandler := user_handler.NewUserHandler(userService)
	return handlers{healthCheckHandler, shortenURLHandler, userHandler, cfg}
}

// initRoutes initializes the routes for the app engine
func (e *engine) initRoutes(cfg *config.Config) {
	allHandlers := e.InitHandlers(cfg)

	e.app.GET("/health-check", allHandlers.healthCheck.Ping)

	docs.SwaggerInfo.BasePath = allHandlers.config.BasePath
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1Routes := e.app.Group("/v1")
	{
		v1Routes.POST("/links/shorten", allHandlers.shorten.ShortenUrl)
		v1Routes.GET("/links/redirect/:code", allHandlers.shorten.Redirect)
		v1Routes.POST("/users/register", allHandlers.user.Register)
	}
}
