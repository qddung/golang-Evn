package service

type HealthCheck interface {
	Ping() (string, error)
}

type healthCheck struct {
}

func NewHealthCheck() HealthCheck {
	return &healthCheck{}
}

func (h *healthCheck) Ping() (string, error) {
	return "OK", nil
}
