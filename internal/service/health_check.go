package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/homework/lab/internal/repository"
)

var (
	errRedisNotAvailable = errors.New("redis not available")
)

//go:generate mockery --name=HealthCheck --filename=health_check.go --outpkg=mocks
type HealthStatusResult struct {
	Message     string
	ServiceName string
	InstanceID  string
}

type HealthCheck interface {
	Ping(ctx context.Context) (HealthStatusResult, error)
}

type healthCheck struct {
	serviceName    string
	instanceID     string
	pingRepository repository.Ping
}

// Nhận cấu hình từ ngoài truyền vào (Dependency Injection)
func NewHealthCheck(serviceName string, instanceID string, pingRepository repository.Ping) HealthCheck {
	return &healthCheck{
		serviceName:    serviceName,
		instanceID:     instanceID,
		pingRepository: pingRepository,
	}
}

func (h *healthCheck) Ping(ctx context.Context) (HealthStatusResult, error) {
	// Hàm này giờ đây cực kỳ "sạch", không phụ thuộc vào file .env nào cả
	pingErr := h.pingRepository.Ping(ctx)
	if pingErr != nil {
		return HealthStatusResult{
			Message:     "Error",
			ServiceName: h.serviceName,
			InstanceID:  h.instanceID,
		}, errors.Wrap(errRedisNotAvailable, pingErr.Error())
	}
	return HealthStatusResult{
		Message:     "OK",
		ServiceName: h.serviceName,
		InstanceID:  h.instanceID,
	}, nil
}
