package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/service"
)

type HealthCheck interface {
	Ping(c *gin.Context)
}

type HealthResponse struct {
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}

type healthCheck struct {
	svc service.HealthCheck
}

func NewHealthCheck(svc service.HealthCheck) HealthCheck {
	return &healthCheck{
		svc: svc,
	}
}
func (h *healthCheck) Ping(c *gin.Context) {
	result, _ := h.svc.Ping()
	handlerResponse := HealthResponse{
		Message:     result.Message,
		ServiceName: result.ServiceName,
		InstanceID:  result.InstanceID,
	}
	c.JSON(http.StatusOK, handlerResponse)
}
