// Package logging provides a unified logging interface using zap
package logging

import (
	"time"

	forbidden_zap "go.uber.org/zap"
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
	Debug(msg string, fields ...interface{})
	// Warn logs a message at warn level with optional fields
	Warn(msg string, fields ...Field)
}

// Field represents a logging field.
type Field interface{}

// logger implements the shared.Logger interface using zap
type logger struct {
	log *forbidden_zap.SugaredLogger
}

// NewLogger creates a new logger instance
func NewLogger() Logger {
	zapLog, err := forbidden_zap.NewProduction()
	if err != nil {
		// Handle error appropriately, e.g., log it or panic
		panic(err) // For simplicity, panic on error
	}
	return &logger{log: zapLog.Sugar()}
}

// NewTestLogger creates a logger suitable for testing
func NewTestLogger() Logger {
	config := forbidden_zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	zapLog, _ := config.Build()
	return &logger{log: zapLog.Sugar()}
}

// Info logs an info message
func (l *logger) Info(msg string, fields ...Field) {
	l.log.Infow(msg, convertFields(fields)...)
}

// Error logs an error message
func (l *logger) Error(msg string, fields ...Field) {
	l.log.Errorw(msg, convertFields(fields)...)
}

// Debug logs a debug message
func (l *logger) Debug(msg string, fields ...interface{}) {
	l.log.Debugw(msg, fields...)
}

// Warn logs a warning message
func (l *logger) Warn(msg string, fields ...Field) {
	l.log.Warnw(msg, convertFields(fields)...)
}

// convertFields converts custom fields to zap fields.
func convertFields(fields []Field) []interface{} {
	converted := make([]interface{}, len(fields))
	for i, field := range fields {
		converted[i] = field // Assuming Field is an interface{}
	}
	return converted
}

// String creates a string field
func String(key string, value string) forbidden_zap.Field {
	return forbidden_zap.String(key, value)
}

// Bool creates a boolean field
func Bool(key string, value bool) forbidden_zap.Field {
	return forbidden_zap.Bool(key, value)
}

// Error creates an error field
func Error(value error) forbidden_zap.Field {
	return forbidden_zap.Error(value)
}

// Uint creates an unsigned integer field
func Uint(key string, value uint) forbidden_zap.Field {
	return forbidden_zap.Uint(key, value)
}

// Int creates an integer field
func Int(key string, value int) forbidden_zap.Field {
	return forbidden_zap.Int(key, value)
}

// Int64 creates an int64 field
func Int64(key string, value int64) forbidden_zap.Field {
	return forbidden_zap.Int64(key, value)
}

// Any creates an interface{} field
func Any(key string, value interface{}) forbidden_zap.Field {
	return forbidden_zap.Any(key, value)
}

// Duration creates a duration field
func Duration(key string, value time.Duration) forbidden_zap.Field {
	return forbidden_zap.Duration(key, value)
}

// FxEventLogger is a logger that integrates with fx
type FxEventLogger struct {
	Logger Logger
}

// Implement any methods you need for FxEventLogger
