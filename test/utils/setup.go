package utils

import (
	"github.com/labstack/echo/v4"
)

// TestSetup contains common test setup utilities
type TestSetup struct {
	Echo *echo.Echo
}

// NewTestSetup creates a new test setup with common configurations
func NewTestSetup() *TestSetup {
	e := echo.New()

	return &TestSetup{
		Echo: e,
	}
}

// Close performs any necessary cleanup
func (ts *TestSetup) Close() error {
	// Add any necessary cleanup logic here
	return nil
}
