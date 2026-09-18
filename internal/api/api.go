package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/homework/lab/docs"
	_ "github.com/homework/lab/docs"
	"github.com/homework/lab/internal/api/middleware"
	bookmark_cache "github.com/homework/lab/internal/cache/bookmark"
	"github.com/homework/lab/internal/config"
	"github.com/homework/lab/internal/connection"
	bookmark_handler "github.com/homework/lab/internal/handler/bookmark"
	health_check_handler "github.com/homework/lab/internal/handler/health_check"
	"github.com/homework/lab/internal/handler/shorten"
	user_handler "github.com/homework/lab/internal/handler/user"
	bookmark_repository "github.com/homework/lab/internal/repository/bookmark"
	"github.com/homework/lab/internal/repository/cache"
	health_check_repository "github.com/homework/lab/internal/repository/health_check"
	url_repository "github.com/homework/lab/internal/repository/shorten"
	userRepository "github.com/homework/lab/internal/repository/user"
	bookmark_service "github.com/homework/lab/internal/service/bookmark"
	health_check_service "github.com/homework/lab/internal/service/health_check"
	shorten_service "github.com/homework/lab/internal/service/shorten"
	user_service "github.com/homework/lab/internal/service/user"
	base62_helper "github.com/homework/lab/pkg/helpers/base62"
	"github.com/homework/lab/pkg/helpers/hasher"
	jwt_pkg "github.com/homework/lab/pkg/jwt"
	base62_lib "github.com/homework/lab/pkg/lib/base62"
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
	app          *gin.Engine
	cfg          *config.Config
	connector    connection.DBConnector
	jwtGenerator jwt_pkg.JwtGenerator
	jwtValidator jwt_pkg.JwtValidator
}

type EnginOpt struct {
	App          *gin.Engine
	Cfg          *config.Config
	Connector    connection.DBConnector
	JwtGenerator jwt_pkg.JwtGenerator
	JwtValidator jwt_pkg.JwtValidator
}

// NewEngine creates a new engine instance
func NewEngine(opt *EnginOpt) Engine {
	api := &engine{
		app:          opt.App,
		cfg:          opt.Cfg,
		connector:    opt.Connector,
		jwtGenerator: opt.JwtGenerator,
		jwtValidator: opt.JwtValidator,
	}

	api.initRoutes(opt.Cfg)
	return api
}

// config Run starts the app engine
func (e *engine) Run() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// override config ServeHTTP serves the app engine
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

type handlers struct {
	healthCheck health_check_handler.HealthCheck
	shorten     shorten.ShorternUrl
	user        user_handler.UserHandler
	bookmark    bookmark_handler.BookmarkHandler
	config      *config.Config
}

func (e *engine) InitHandlers(cfg *config.Config) handlers {
	serviceName := cfg.ServiceName
	instanceID := cfg.InstanceID
	redisClient := e.connector.GetRedisClient()
	sqlDB := e.connector.GetSqlDB()
	cacheRedis := cache.NewCache(redisClient)

	// lib
	base62_lib := base62_lib.NewStdEncoding()
	// create helper
	hasher := hasher.NewHasher()
	// code_gen := code_gen.NewKeyGenerator()
	base62_helper := base62_helper.New(base62_lib)
	// create repository
	healthCheckRepository := health_check_repository.NewPing(redisClient)
	urlStorage := url_repository.NewURLStorage(redisClient)
	userRepository := userRepository.NewUserRepository(sqlDB)
	bookmarkRepo := bookmark_repository.NewBookmarkRepository(sqlDB, base62_helper)
	// create service
	healthCheckService := health_check_service.NewHealthCheck(serviceName, instanceID, healthCheckRepository)
	shortenService := shorten_service.NewShorternUrl(urlStorage, base62_helper, bookmarkRepo)
	userService := user_service.NewUserService(userRepository, hasher, e.jwtGenerator)
	bookmarkSvc := bookmark_service.NewBookmarkService(bookmarkRepo)

	// create cache
	bookmarkCache := bookmark_cache.NewBookmarkCacheInstance(bookmarkSvc, cacheRedis)
	// create handler
	healthCheckHandler := health_check_handler.NewHealthCheck(healthCheckService)
	shortenURLHandler := shorten.NewShortenURL(shortenService)
	userHandler := user_handler.NewUserHandler(userService)
	bookmarkHandler := bookmark_handler.NewBookmarkHandler(bookmarkCache)

	return handlers{healthCheckHandler, shortenURLHandler, userHandler, bookmarkHandler, cfg}
}

func (e *engine) initRoutes(cfg *config.Config) {
	allHandlers := e.InitHandlers(cfg)

	e.app.GET("/health-check", allHandlers.healthCheck.Ping)

	docs.SwaggerInfo.BasePath = allHandlers.config.BasePath
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	jwtMiddleware := middleware.NewJwtAuthMiddleware(e.jwtValidator)

	v1Routes := e.app.Group("/v1")
	{
		// --- Public API
		v1Routes.POST("/users/login", allHandlers.user.Login)
		v1Routes.POST("/links/shorten", allHandlers.shorten.ShortenUrl)
		v1Routes.GET("/links/redirect/:code", allHandlers.shorten.Redirect)
		v1Routes.POST("/users/register", allHandlers.user.Register)

		// -- Private Api
		v1Routes.Use(jwtMiddleware.JwtAuth()) // middelware
		v1Routes.GET("/self/info", allHandlers.user.GetUserInfo)
		v1Routes.PUT("/self/info", allHandlers.user.UpdateUserInfo)

		v1Routes.POST("/bookmarks", allHandlers.bookmark.CreateBookmark)
		v1Routes.GET("/bookmarks", allHandlers.bookmark.GetBookmarks)
		v1Routes.PUT("/bookmarks/:id", allHandlers.bookmark.UpdateBookmark)
		v1Routes.DELETE("/bookmarks/:id", allHandlers.bookmark.DeleteBookmark)

	}
}
