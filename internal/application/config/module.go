package config

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"

	"github.com/jonesrussell/goforms/internal/application/loggingconfig"
)

// Module provides the configuration dependencies
var Module = fx.Module("config",
	fx.Provide(
		New, // Provide the New function to create a Config instance
		fx.Annotate(func(cfg *Config) loggingconfig.LoggerConfig {
			return cfg // Pass the LoggerConfig interface to the logger
		}, fx.As(new(loggingconfig.LoggerConfig))), // Provide the LoggerConfig interface
	),
)

// New creates a new configuration instance
func New() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process config: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}
