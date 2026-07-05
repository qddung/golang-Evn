package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// Config struct for app config
type Config struct {
	AppPort     string
	ServiceName string
	InstanceID  string
}

func LoadConfig() (*Config, error) {
	envPath, err := findEnvFile()
	if err != nil {
		return nil, err
	}
	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			return nil, fmt.Errorf("load .env from %s: %w", envPath, err)
		}
	}

	instanceID := os.Getenv("INSTANCE_ID")
	serviceName := os.Getenv("SERVICE_NAME")
	appPort := os.Getenv("APP_PORT")
	if instanceID == "" {
		instanceID = uuid.New().String()
	}
	cfg := &Config{
		AppPort:     appPort,
		ServiceName: serviceName,
		InstanceID:  instanceID,
	}
	return cfg, nil
}

func findEnvFile() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for dir := wd; ; dir = filepath.Dir(dir) {
		candidate := filepath.Join(dir, ".env")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}

	return "", nil
}
