// Package logging provides a unified logging interface using zap
package logging

import (
	"time"

	forbidden_zap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the interface for logging operations
//
// This interface abstracts the underlying logging implementation,
// allowing for easy mocking in tests and flexibility to change
// the logging backend without affecting application code.
//
// For testing, use test/mocks.Logger instead of implementing this interface directly.
type Logger interface {
	// Info logs a message at info level with optional fields
	Info(msg string, fields ...Field)
	// Error logs a message at error level with optional fields
	Error(msg string, fields ...Field)
	// Debug logs a message at debug level with optional fields
	Debug(msg string, fields ...Field)
	// Warn logs a message at warn level with optional fields
	Warn(msg string, fields ...Field)
	// Int64 adds an int64 field to the log entry
	Int64(key string, value int64) Field
	// Int adds an int field to the log entry
	Int(key string, value int) Field
	// Int32 adds an int32 field to the log entry
	Int32(key string, value int32) Field
	// Uint64 adds a uint64 field to the log entry
	Uint64(key string, value uint64) Field
	// Uint adds a uint field to the log entry
	Uint(key string, value uint) Field
	// Uint32 adds a uint32 field to the log entry
	Uint32(key string, value uint32) Field
}

// Field represents a logging field
type Field = forbidden_zap.Field

// String creates a string field
func String(key string, value string) Field { return forbidden_zap.String(key, value) }

// Int creates an integer field
func Int(key string, value int) Field { return forbidden_zap.Int(key, value) }

// Int64 creates a 64-bit integer field
func Int64(key string, value int64) Field { return forbidden_zap.Int64(key, value) }

// Uint creates an unsigned integer field
func Uint(key string, value uint) Field { return forbidden_zap.Uint(key, value) }

// Bool creates a boolean field
func Bool(key string, value bool) Field { return forbidden_zap.Bool(key, value) }

// Error creates an error field
func Error(err error) Field { return forbidden_zap.Error(err) }

// Duration creates a duration field
func Duration(key string, value time.Duration) Field { return forbidden_zap.Duration(key, value) }

// Any creates a field with any value
func Any(key string, value interface{}) Field { return forbidden_zap.Any(key, value) }

// logger implements the Logger interface using zap
type logger struct {
	log *forbidden_zap.Logger
}

// NewLogger creates a new logger instance
func NewLogger(debug bool, appName string) Logger {
	// Create encoder config
	encoderConfig := forbidden_zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	var zapLog *forbidden_zap.Logger
	if debug {
		// Development mode with colored output
		config := forbidden_zap.NewDevelopmentConfig()
		config.EncoderConfig = encoderConfig
		config.OutputPaths = []string{"stdout"}
		config.Encoding = "console"

		zapLog, _ = config.Build(
			forbidden_zap.AddCaller(),
			forbidden_zap.AddStacktrace(zapcore.ErrorLevel),
			forbidden_zap.Fields(
				forbidden_zap.String("app", appName),
			),
		)
	} else {
		// Production mode with JSON output
		zapLog, _ = forbidden_zap.NewProduction(
			forbidden_zap.Fields(
				forbidden_zap.String("app", appName),
			),
		)
	}

	return &logger{log: zapLog}
}

// NewTestLogger creates a logger suitable for testing
func NewTestLogger() Logger {
	config := forbidden_zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	zapLog, _ := config.Build()
	return &logger{log: zapLog}
}

func (l *logger) Info(msg string, fields ...Field)  { l.log.Info(msg, fields...) }
func (l *logger) Error(msg string, fields ...Field) { l.log.Error(msg, fields...) }
func (l *logger) Debug(msg string, fields ...Field) { l.log.Debug(msg, fields...) }
func (l *logger) Warn(msg string, fields ...Field)  { l.log.Warn(msg, fields...) }

// Int64 adds an int64 field to the log entry
func (l *logger) Int64(key string, value int64) Field {
	return forbidden_zap.Int64(key, value)
}

// Int adds an int field to the log entry
func (l *logger) Int(key string, value int) Field {
	return forbidden_zap.Int(key, value)
}

// Int32 adds an int32 field to the log entry
func (l *logger) Int32(key string, value int32) Field {
	return forbidden_zap.Int32(key, value)
}

// Uint64 adds a uint64 field to the log entry
func (l *logger) Uint64(key string, value uint64) Field {
	return forbidden_zap.Uint64(key, value)
}

// Uint adds a uint field to the log entry
func (l *logger) Uint(key string, value uint) Field {
	return forbidden_zap.Uint(key, value)
}

// Uint32 adds a uint32 field to the log entry
func (l *logger) Uint32(key string, value uint32) Field {
	return forbidden_zap.Uint32(key, value)
}
