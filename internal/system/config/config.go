package config

import (
	"fmt"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	Port     int    `envconfig:"HTTP_PORT"`
	LogLevel string `envconfig:"LOG_LEVEL"`
	LRUSize  int    `envconfig:"LRU_SIZE,default=3"`
	Jobs     *Jobs
	Redis    *Redis
}

func Init() (*Config, error) {
	cfg := Config{}
	if err := envconfig.Init(&cfg); err != nil {
		return nil, fmt.Errorf("failed to init config from env: %w", err)
	}
	return &cfg, nil
}
