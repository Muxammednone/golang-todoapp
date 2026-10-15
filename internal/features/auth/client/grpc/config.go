package auth_client_grpc

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr    string        `envconfig:"ADDR" default:"127.0.0.1:50051"`
	Timeout time.Duration `envconfig:"TIMEOUT" default:"5s"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("SSO", &cfg); err != nil {
		return nil, fmt.Errorf("process envconfig: %w", err)
	}
	if cfg.Addr == "" {
		return nil, fmt.Errorf("missing SSO_ADDR environment variable")
	}
	if cfg.Timeout <= 0 {
		return nil, fmt.Errorf("invalid SSO_TIMEOUT environment variable: must be greater than 0")
	}

	return &cfg, nil
}

func NewConfigMust() *Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("failed to create config: %w", err))
	}
	return cfg
}
