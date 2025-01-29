package utils

import (
	"log" // Import log package for default logging

	"github.com/jonesrussell/goforms/internal/application/logging"
)

// MockLogger is a mock implementation of the logging.Logger interface
type MockLogger struct {
	ErrorFunc func(msg string, fields ...logging.Field)
	DebugFunc func(msg string, fields ...interface{})
	InfoFunc  func(msg string, fields ...logging.Field)
	WarnFunc  func(msg string, fields ...logging.Field)
}

// Log implements logging.Logger.
func (m *MockLogger) Log(message string) {
	panic("unimplemented")
}

func (m *MockLogger) Info(msg string, fields ...logging.Field) {
	if m.InfoFunc != nil {
		m.InfoFunc(msg, fields...)
	}
}

func (m *MockLogger) Error(msg string, fields ...logging.Field) {
	if m.ErrorFunc != nil {
		m.ErrorFunc(msg, fields...)
	}
}

func (m *MockLogger) Debug(msg string, fields ...interface{}) {
	if m.DebugFunc == nil {
		// Handle the case where DebugFunc is not set
		log.Printf("Debug called but DebugFunc is not set: %s", msg) // Log a warning
		return
	}
	m.DebugFunc(msg, fields...)
}

func (m *MockLogger) Warn(msg string, fields ...logging.Field) {
	if m.WarnFunc != nil {
		m.WarnFunc(msg, fields...)
	}
}

func (m *MockLogger) WithPrefix(prefix string) logging.Logger {
	// Implement mock logic
	return m
}

func (m *MockLogger) LogWithPrefix(level string, prefix, msg string, fields ...logging.Field) {
	// Implement mock logic
}
