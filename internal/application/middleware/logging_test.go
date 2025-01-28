package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/jonesrussell/goforms/internal/application/middleware"
	mocklogging "github.com/jonesrussell/goforms/test/mocks/logging"
)

func TestLoggingMiddleware(t *testing.T) {
	t.Run("logs request and response details", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create mock logger
		mockLogger := mocklogging.NewMockLogger()
		mockLogger.ExpectInfo("http request")

		// Create middleware
		mw := middleware.LoggingMiddleware(mockLogger)

		// Create test handler
		handler := func(c echo.Context) error {
			c.Response().WriteHeader(http.StatusOK)
			return nil
		}

		// Execute middleware
		err := mw(handler)(c)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if err := mockLogger.Verify(); err != nil {
			t.Errorf("logger expectations not met: %v", err)
		}
	})

	t.Run("logs error when handler fails", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create mock logger
		mockLogger := mocklogging.NewMockLogger()
		mockLogger.ExpectInfo("http request")

		// Create middleware
		mw := middleware.LoggingMiddleware(mockLogger)

		// Create test handler that returns error
		handler := func(c echo.Context) error {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return echo.NewHTTPError(http.StatusInternalServerError, "test error")
		}

		// Execute middleware
		err := mw(handler)(c)
		if err == nil {
			t.Error("expected error, got nil")
		}

		if err := mockLogger.Verify(); err != nil {
			t.Errorf("logger expectations not met: %v", err)
		}
	})
}

func TestLoggingMiddleware_RealIP(t *testing.T) {
	// Create a mock logger for testing
	mockLogger := mocklogging.NewMockLogger()
	mockLogger.ExpectInfo("http request")

	// Create Echo instance
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Real-IP", "192.168.1.1")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware
	mw := middleware.LoggingMiddleware(mockLogger)
	handler := mw(func(c echo.Context) error {
		c.Response().WriteHeader(http.StatusOK)
		return c.String(http.StatusOK, "success")
	})

	// Execute request
	err := handler(c)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify logs
	if err := mockLogger.Verify(); err != nil {
		t.Errorf("logger expectations not met: %v", err)
	}
}

func TestLogging(t *testing.T) {
	mockLogger := mocklogging.NewMockLogger()

	// Set expectation without WithFields
	mockLogger.ExpectInfo("http request")

	// Call the function that uses the logger
	// Example: middleware.HandleRequest(mockLogger)

	// Verify the expectations
	if err := mockLogger.Verify(); err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
}
