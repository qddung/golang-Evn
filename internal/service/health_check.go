package service

//go:generate mockery --name=HealthCheck --filename=health_check.go --outpkg=mocks_health_check
type HealthStatusResult struct {
	Message     string
	ServiceName string
	InstanceID  string
}

type HealthCheck interface {
	Ping() (HealthStatusResult, error)
}

type healthCheck struct {
	serviceName string
	instanceID  string
}

// Nhận cấu hình từ ngoài truyền vào (Dependency Injection)
func NewHealthCheck(serviceName string, instanceID string) HealthCheck {
	return &healthCheck{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

func (h *healthCheck) Ping() (HealthStatusResult, error) {
	// Hàm này giờ đây cực kỳ "sạch", không phụ thuộc vào file .env nào cả
	return HealthStatusResult{
		Message:     "OK",
		ServiceName: h.serviceName,
		InstanceID:  h.instanceID,
	}, nil
}
