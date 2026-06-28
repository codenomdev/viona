package config

import (
	"context"
	"fmt"
	"os"

	appConfig "github.com/codenomdev/viona/config"
	"github.com/codenomdev/viona/pkg/log"
	"github.com/qiangxue/go-env"
	"gopkg.in/yaml.v3"
)

// Config represents an application configuration.
type Config struct {
	// Main host config
	HOST HostConfig `yaml:"HOST"`
	// Postgres sql
	POSTGRES PostgresConfig `yaml:"POSTGRES"`
	// Sony flake
	SONYFLAKE SonyflakeConfig `yaml:"SONYFLAKE"`
	// UI
	UI UIConfig `yaml:"UI"`
}

type depsKey struct{}

// Load returns an application configuration which is populated from the given configuration file and environment variables.
func LoadConfig(filename string, logger log.Logger) (*Config, error) {
	c := Config{
		HOST:      HostConfigInit(),
		POSTGRES:  DatabaseConfigInit(),
		SONYFLAKE: NewSonyFlakeInit(),
		UI:        UIConfigInit(),
	}

	// get config validate
	cfgLoad := ValidateConfEnv(filename)
	// load from YAML config file
	bytes, err := appConfig.DefaultConfig.ReadFile(cfgLoad)

	if err != nil {
		logger.Error(fmt.Sprintf("failed to load application configuration: %s, path: %s", err, cfgLoad))
		return nil, err
	}

	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	var logFn env.LogFunc

	if os.Getenv("APP_ENV") == "dev" {
		logFn = func(format string, args ...interface{}) {
			logger.Debug(fmt.Sprintf(format, args...))
		}
	} else {
		logFn = func(format string, args ...interface{}) {
			logger.Info(fmt.Sprintf(format, args...))
		}
	}

	// load from environment variables
	// if not set environment variables, will use config with yaml & default config if existing config.
	//
	// load from environment variables HOST config prefixed with "HOST_"
	// e.g: HOST_PORT, HOST_DOMAIN, CONTEXT_TIMEOUT, SERVER_DEBUG, SERVER_WRITE_TIMEOUT, SERVER_READ_TIMEOUT.
	if err = env.New("HOST_", logFn).Load(&c.HOST); err != nil {
		return nil, err
	}

	// load from environment variables HOST SSL prefixed with "SSL_"
	// e.g: SSL_ENABLE, SSL_CERT_KEY, SSL_PRIV_KEY
	if err := env.New("SSL_", logFn).Load(&c.HOST.SSL); err != nil {
		return nil, err
	}
	// load from environment variables DATABASE POSTGRE SQL prefixed with "DATABASE_"
	// e.g: DATABASE_HOST, DATABASE_PORT, DATABASE_USERNAME, DATABASE_PASSWORD,
	// 		DATABASE_NAME, DATABASE_SSL_MODE, DATABASE_TIMEZONE
	if err := env.New("DATABASE_", logFn).Load(&c.POSTGRES); err != nil {
		return nil, err
	}

	// load from environment
	if err := env.New("UI_", logFn).Load(&c.UI); err != nil {
		return nil, err
	}

	return &c, nil
}

// Validate load config file environment
// auto detect env for load config yaml.
func ValidateConfEnv(configPath string) string {
	flag := "config"
	switch os.Getenv("APP_ENV") {
	case "dev", "development": // get identity environment, can "dev" or "development"
		return fmt.Sprintf("%s.dev.yaml", flag)
	case "prod", "production": // get identity environment, can "prod" or "production"
		return fmt.Sprintf("%s.prod.yaml", flag)
	default:
		return configPath
	}
}

func ToContext(ctx context.Context, deps *Config) context.Context {
	return context.WithValue(ctx, depsKey{}, deps)
}

func FromContext(ctx context.Context) *Config {
	if v, ok := ctx.Value(depsKey{}).(*Config); ok {
		return v
	}
	return nil
}
