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

// @Summary Health check endpoint
// @Description Check the health of the API
// @Tags health_check
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /ping [get]
func (h *healthCheck) Ping(c *gin.Context) {
	result, _ := h.svc.Ping(c)
	handlerResponse := HealthResponse{
		Message:     result.Message,
		ServiceName: result.ServiceName,
		InstanceID:  result.InstanceID,
	}
	c.JSON(http.StatusOK, handlerResponse)
}
