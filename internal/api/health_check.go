package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/homework/lab/internal/config"
	"github.com/homework/lab/internal/handler"
	"github.com/homework/lab/internal/service"
)

// Engine interface for app engine
type Engine interface {
	Run() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// engine struct for app engine
type engine struct {
	app *gin.Engine
	cfg *config.Config
}

// NewEngine creates a new engine instance
func NewEngine(cfg *config.Config) Engine {
	api := &engine{
		app: gin.New(),
		cfg: cfg,
	}

	api.initRoutes(cfg.ServiceName, cfg.InstanceID)
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

// initRoutes initializes the routes for the app engine
func (e *engine) initRoutes(serviceName string, instanceID string) { // (serviceName, instanceID) {
	// create service
	healthCheckService := service.NewHealthCheck(serviceName, instanceID)
	// create handler
	healthCheckHandler := handler.NewHealthCheck(healthCheckService)
	// health check
	e.app.GET("/ping", healthCheckHandler.Ping)
}
