package config

import (
	"log"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

// Config struct for app config
type Config struct {
	AppPort     string `envconfig:"APP_PORT" default:"8080"`
	ServiceName string `envconfig:"SERVICE_NAME" default:"book"`
	InstanceID  string `envconfig:"INSTANCE_ID"`
	BasePath    string `envconfig:"BASE_PATH" default:"/"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.NewString()
	}

	return &cfg, nil
}
