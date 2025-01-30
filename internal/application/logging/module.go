package logging

import (
	"fmt"

	"go.uber.org/fx"
	forbidden_zap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jonesrussell/goforms/internal/application/loggingconfig"
)

// Declare logInstance at the package level
var logInstance Logger // Declare logInstance

// Module provides logging-related dependencies.
var Module = fx.Options(
	fx.Provide(
		loggingconfig.NewConfig, // Ensure this is only provided here
		NewLogger,
		NewFactory, // Added Factory to provide logger creation capabilities
		NewFxLogger,
	),
)

// NewLogger creates a new logger instance based on the environment configuration
func NewLogger(cfg loggingconfig.LoggerConfigInterface) Logger {
	var zapLog *forbidden_zap.Logger
	var zapSugaredLog *forbidden_zap.SugaredLogger
	var err error

	// Check if the environment is development
	if cfg.GetEnv() == "development" {
		zapConfig := forbidden_zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Enable colorized output
		zapLog, err = zapConfig.Build()                                        // Build the development logger
		if err != nil {
			logInstance.Error("Failed to create development logger", forbidden_zap.Error(err)) // Use logInstance
			panic(fmt.Errorf("failed to create development logger: %w", err))                  // Log error before panicking
		}
		zapSugaredLog = zapLog.Sugar() // Convert to SugaredLogger
	} else {
		prodLog, err := forbidden_zap.NewProduction() // Create the production logger
		if err != nil {
			logInstance.Error("Failed to create production logger", forbidden_zap.Error(err)) // Use logInstance
			panic(fmt.Errorf("failed to create production logger: %w", err))                  // Log error before panicking
		}
		zapSugaredLog = prodLog.Sugar() // Convert to SugaredLogger
	}

	return &logger{log: zapSugaredLog} // Return the logger
}

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
