package sqldb

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// config represents the configuration of redis
type config struct {
	Host     string `default:"localhost:6379" envconfig:"PG_SQL_HOST"`
	Password string `default:"" envconfig:"PG_SQL_PWD"`
	DbName   string `default:"" envconfig:"PG_SQL_DB"`
	User     string `default:"" envconfig:"PG_SQL_USER"`
	Port     string `default:"5432" envconfig:"PG_SQL_PORT"`
	SslMode  string `default:"disable" envconfig:"PG_SQL_SSLMODE"`
	TimeZone string `default:"Asia/Shanghai" envconfig:"PG_SQL_TIMEZONE"`
}

// newConfig creates a new config struct based on the environment variables
func getConfig() (*config, error) {
	cfg := &config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *config) getDSN() string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DbName,
		cfg.Port,
		cfg.SslMode,
		cfg.TimeZone,
	)
	return dsn
}
