package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/service"
)

type HealthCheck interface {
	Ping(c *gin.Context)
}

type healthCheck struct {
	svc service.HealthCheck
}

func NewHealthCheck() HealthCheck {
	return &healthCheck{}
}
func (h *healthCheck) Ping(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
