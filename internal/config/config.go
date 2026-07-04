package config

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

// Config struct for app config
type Config struct {
	AppPort     string `envconfig:"APP_PORT" default:"8080"`
	ServiceName string `envconfig:"SERVICE_NAME" default:"bookmark_service"`
	InstanceID  string `envconfig:"INSTANCE_ID"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}
	if cfg.InstanceID == "" {
		uuid := uuid.New()
		cfg.InstanceID = uuid.String()
	}

	return cfg, nil
}
