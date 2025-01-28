package logging

import (
	"go.uber.org/fx"
)

// Module provides the logging dependencies
var Module = fx.Module("logging",
	fx.Provide(
		NewFactory, // Provide the logger factory
		fx.Annotate(func(f *Factory) Logger {
			return f.CreateFromConfig() // Provide the logger using the factory
		}, fx.As(new(Logger))), // Provide the logger
	),
)

// Factory creates loggers based on configuration
type Factory struct{}

// NewFactory creates a new logger factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreateFromConfig creates a logger from configuration
func (f *Factory) CreateFromConfig() Logger {
	// Here you can implement logic to create a logger based on config
	return NewTestLogger() // Placeholder for actual config-based logger creation
}

// CreateTestLogger creates a logger for testing
func (f *Factory) CreateTestLogger() Logger {
	return NewTestLogger()
}
