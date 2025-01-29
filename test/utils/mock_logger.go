package utils

import "github.com/jonesrussell/goforms/internal/application/logging"

// MockLogger is a mock implementation of the logging.Logger interface
type MockLogger struct {
	ErrorFunc func(msg string, fields ...logging.Field)
	DebugFunc func(msg string, fields ...interface{})
}

func (m *MockLogger) Info(msg string, fields ...logging.Field) {
	// Implement mock logic
}

func (m *MockLogger) Error(msg string, fields ...logging.Field) {
	m.ErrorFunc(msg, fields)
}

func (m *MockLogger) Debug(msg string, fields ...interface{}) {
	m.DebugFunc(msg, fields)
}

func (m *MockLogger) Warn(msg string, fields ...logging.Field) {
	// Implement mock logic
}

func (m *MockLogger) WithPrefix(prefix string) logging.Logger {
	// Implement mock logic
	return m
}

func (m *MockLogger) LogWithPrefix(level string, prefix, msg string, fields ...logging.Field) {
	// Implement mock logic
}
