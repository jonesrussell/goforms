package mocklogging

import (
	"fmt"
	"strings"
	"sync"
	"time"

	forbidden_zap "go.uber.org/zap"

	"github.com/jonesrussell/goforms/internal/application/logging"
)

// AnyValue is a placeholder for any field value
type AnyValue struct{}

// logCall represents a single logging call
type logCall struct {
	level   string
	message string
	fields  map[string]interface{}
}

// MockLogger is a mock implementation of logging.Logger
type MockLogger struct {
	mu       sync.Mutex
	calls    []logCall
	expected []logCall
}

// NewMockLogger creates a new mock logger
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

func (m *MockLogger) recordCall(level, message string, fields ...logging.Field) {
	m.mu.Lock()
	defer m.mu.Unlock()
	fieldMap := make(map[string]interface{})
	for _, field := range fields {
		fieldMap[field.Key] = field.Interface
	}
	m.calls = append(m.calls, logCall{level: level, message: message, fields: fieldMap})
}

func (m *MockLogger) Info(message string, fields ...logging.Field) {
	m.recordCall("info", message, fields...)
}

func (m *MockLogger) Error(message string, fields ...logging.Field) {
	m.recordCall("error", message, fields...)
}

func (m *MockLogger) Debug(message string, fields ...logging.Field) {
	m.recordCall("debug", message, fields...)
}

func (m *MockLogger) Warn(message string, fields ...logging.Field) {
	m.recordCall("warn", message, fields...)
}

// Field creation methods
func (m *MockLogger) Int64(key string, value int64) logging.Field {
	return forbidden_zap.Int64(key, value)
}

func (m *MockLogger) Int(key string, value int) logging.Field {
	return forbidden_zap.Int(key, value)
}

func (m *MockLogger) Int32(key string, value int32) logging.Field {
	return forbidden_zap.Int32(key, value)
}

func (m *MockLogger) Uint64(key string, value uint64) logging.Field {
	return forbidden_zap.Uint64(key, value)
}

func (m *MockLogger) Uint(key string, value uint) logging.Field {
	return forbidden_zap.Uint(key, value)
}

func (m *MockLogger) Uint32(key string, value uint32) logging.Field {
	return forbidden_zap.Uint32(key, value)
}

func (m *MockLogger) String(key string, value string) logging.Field {
	return forbidden_zap.String(key, value)
}

func (m *MockLogger) Bool(key string, value bool) logging.Field {
	return forbidden_zap.Bool(key, value)
}

func (m *MockLogger) ErrorField(err error) logging.Field {
	return forbidden_zap.Error(err)
}

func (m *MockLogger) Duration(key string, value time.Duration) logging.Field {
	return forbidden_zap.Duration(key, value)
}

// Expectation methods
func (m *MockLogger) ExpectInfo(message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	call := logCall{level: "info", message: message, fields: make(map[string]interface{})}
	m.expected = append(m.expected, call)
}

func (m *MockLogger) ExpectError(message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	call := logCall{level: "error", message: message, fields: make(map[string]interface{})}
	m.expected = append(m.expected, call)
}

func (m *MockLogger) ExpectDebug(message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	call := logCall{level: "debug", message: message, fields: make(map[string]interface{})}
	m.expected = append(m.expected, call)
}

func (m *MockLogger) ExpectWarn(message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	call := logCall{level: "warn", message: message, fields: make(map[string]interface{})}
	m.expected = append(m.expected, call)
}

// Verify checks if all expected calls were made
func (m *MockLogger) Verify() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.expected) != len(m.calls) {
		return fmt.Errorf("expected %d calls but got %d", len(m.expected), len(m.calls))
	}

	for i, exp := range m.expected {
		got := m.calls[i]
		if exp.level != got.level {
			return fmt.Errorf("call %d: expected level %q but got %q", i, exp.level, got.level)
		}
		if !strings.Contains(got.message, exp.message) {
			return fmt.Errorf("call %d: expected message %q but got %q", i, exp.message, got.message)
		}
	}

	return nil
}

// Reset clears all calls and expectations
func (m *MockLogger) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = nil
	m.expected = nil
}

// Ensure MockLogger implements logging.Logger
var _ logging.Logger = (*MockLogger)(nil)

// WithFields adds field expectations to a log call
func (m *MockLogger) WithFields(fields map[string]interface{}) logging.Field {
	// You can implement this to return a logging.Field that your logger can use
	return logging.Field{
		Key:       "fields",
		Interface: fields,
	}
}
