package utils

import "github.com/jonesrussell/goforms/internal/application/logging"

// MockLogger is a mock implementation of the logging.Logger interface
type MockLogger struct {
	DebugFunc         func(msg string, fields ...interface{})
	ErrorFunc         func(msg string, fields ...logging.Field)
	InfoFunc          func(msg string, fields ...logging.Field)
	WarnFunc          func(msg string, fields ...logging.Field)
	FatalFunc         func(msg string, fields ...logging.Field)
	LogFunc           func(msg string)
	LogWithPrefixFunc func(prefix string, msg string, level string, fields ...logging.Field)
	SyncFunc          func() error
	WithPrefixFunc    func(prefix string) logging.Logger
}

// Implement the logging.Logger interface methods
func (m *MockLogger) Debug(msg string, fields ...interface{}) {
	m.DebugFunc(msg, fields...)
}

func (m *MockLogger) Error(msg string, fields ...logging.Field) {
	m.ErrorFunc(msg, fields...)
}

func (m *MockLogger) Info(msg string, fields ...logging.Field) {
	m.InfoFunc(msg, fields...)
}

func (m *MockLogger) Warn(msg string, fields ...logging.Field) {
	m.WarnFunc(msg, fields...)
}

// Implement the Fatal method
func (m *MockLogger) Fatal(msg string, fields ...logging.Field) {
	m.FatalFunc(msg, fields...)
}

// Implement the Log method
func (m *MockLogger) Log(msg string) {
	if m.LogFunc != nil {
		m.LogFunc(msg)
	}
}

// Implement the LogWithPrefix method
func (m *MockLogger) LogWithPrefix(prefix string, msg string, level string, fields ...logging.Field) {
	if m.LogWithPrefixFunc != nil {
		m.LogWithPrefixFunc(prefix, msg, level, fields...)
	}
}

// Implement the Sync method
func (m *MockLogger) Sync() error {
	if m.SyncFunc != nil {
		return m.SyncFunc()
	}
	return nil
}

// Implement the WithPrefix method
func (m *MockLogger) WithPrefix(prefix string) logging.Logger {
	if m.WithPrefixFunc != nil {
		return m.WithPrefixFunc(prefix)
	}
	return m
}

// Implement the methods as needed
