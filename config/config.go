// Configuration loading (env, files)
package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL environment variable is required")
	}
	if cfg.Port == "" {
		cfg.Port = "8080" // default port
	}
	return cfg, nil
}
