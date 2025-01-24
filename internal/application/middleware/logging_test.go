package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/jonesrussell/goforms/internal/application/middleware"
	"github.com/jonesrussell/goforms/internal/infrastructure/logging"
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
		mockLogger.ExpectInfo("http request",
			logging.String("method", "GET"),
			logging.String("path", "/test"),
			logging.String("ip", "192.0.2.1"),
		)

		// Create middleware
		mw := middleware.LoggingMiddleware(mockLogger)

		// Create test handler
		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "test response")
		}

		// Execute middleware
		err := mw(handler)(c)

		// Assert
		assert.NoError(t, err)
		mockLogger.AssertExpectations(t)
	})

	t.Run("logs error when handler fails", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create mock logger
		mockLogger := mocklogging.NewMockLogger()
		mockLogger.ExpectInfo("http request",
			logging.String("method", "GET"),
			logging.String("path", "/test"),
			logging.String("ip", "192.0.2.1"),
		)
		mockLogger.ExpectError("request failed",
			logging.String("error", "test error"),
		)

		// Create middleware
		mw := middleware.LoggingMiddleware(mockLogger)

		// Create test handler that returns error
		handler := func(c echo.Context) error {
			return echo.NewHTTPError(http.StatusInternalServerError, "test error")
		}

		// Execute middleware
		err := mw(handler)(c)

		// Assert
		assert.Error(t, err)
		mockLogger.AssertExpectations(t)
	})
}

func TestLoggingMiddleware_RealIP(t *testing.T) {
	// Create a mock logger for testing
	mockLogger := mocklogging.NewMockLogger()
	mockLogger.ExpectInfo("http request",
		logging.String("method", "GET"),
		logging.String("path", "/test"),
		logging.String("ip", "192.168.1.1"),
	)

	// Create Echo instance
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Real-IP", "192.168.1.1")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware
	middleware := middleware.LoggingMiddleware(mockLogger)
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute request
	_ = handler(c)

	// Verify logs
	mockLogger.AssertExpectations(t)
}
