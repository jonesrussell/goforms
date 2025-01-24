package utils

import (
	"github.com/labstack/echo/v4"

	"github.com/jonesrussell/goforms/internal/infrastructure/logging"
)

// TestSetup contains common test setup utilities
type TestSetup struct {
	Echo   *echo.Echo
	Logger logging.Logger
}

// NewTestSetup creates a new test setup with common configurations
func NewTestSetup() *TestSetup {
	e := echo.New()
	mockLogger := logging.NewMockLogger()

	return &TestSetup{
		Echo:   e,
		Logger: mockLogger,
	}
}

// Close performs any necessary cleanup
func (ts *TestSetup) Close() error {
	return nil
}
