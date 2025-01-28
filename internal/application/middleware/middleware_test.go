package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	mocklogging "github.com/jonesrussell/goforms/test/mocks/logging"
)

func TestMiddlewareSetup(t *testing.T) {
	// Create mock logger
	mockLogger := mocklogging.NewMockLogger()
	mockLogger.ExpectDebug("creating new middleware manager")
	mockLogger.ExpectDebug("setting up middleware")
	mockLogger.ExpectDebug("adding security headers middleware")
	mockLogger.ExpectDebug("adding request ID middleware")
	mockLogger.ExpectDebug("middleware setup complete")

	// Create middleware manager
	mw := New(mockLogger)

	// Create Echo instance
	e := echo.New()

	// Setup middleware
	mw.Setup(e)

	// Verify logger calls
	if err := mockLogger.Verify(); err != nil {
		t.Errorf("logger expectations not met: %v", err)
	}
}

func TestRequestIDMiddleware(t *testing.T) {
	mockLogger := mocklogging.NewMockLogger()
	mockLogger.ExpectDebug("creating new middleware manager")
	mockLogger.ExpectDebug("processing request ID middleware")
	mockLogger.ExpectDebug("request ID middleware complete")

	e := echo.New()
	m := New(mockLogger)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Add middleware
	e.Use(m.requestID())

	// Create test handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	// Process the request
	e.ServeHTTP(rec, req)

	// Assert response
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	if err := mockLogger.Verify(); err != nil {
		t.Errorf("logger expectations not met: %v", err)
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	mockLogger := mocklogging.NewMockLogger()
	mockLogger.ExpectDebug("processing security headers")
	mockLogger.ExpectDebug("security headers processing complete")

	e := echo.New()
	m := New(mockLogger)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Add middleware
	e.Use(m.securityHeaders())

	// Create test handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	// Process the request
	e.ServeHTTP(rec, req)

	// Check security headers
	expectedHeaders := map[string]string{
		"Content-Security-Policy":      "default-src 'self';",
		"X-Content-Type-Options":       "nosniff",
		"X-Frame-Options":              "DENY",
		"X-XSS-Protection":             "1; mode=block",
		"Referrer-Policy":              "no-referrer",
		"Permissions-Policy":           "geolocation=(), microphone=(), camera=()",
		"Cross-Origin-Opener-Policy":   "same-origin",
		"Cross-Origin-Embedder-Policy": "require-corp",
		"Cross-Origin-Resource-Policy": "same-origin",
	}

	for header, expected := range expectedHeaders {
		got := rec.Header().Get(header)
		if got != expected {
			t.Errorf("expected %s header to be %q, got %q", header, expected, got)
		}
	}

	if err := mockLogger.Verify(); err != nil {
		t.Errorf("logger expectations not met: %v", err)
	}
}

func TestMiddleware(t *testing.T) {
	mockLogger := mocklogging.NewMockLogger()

	// Set expectation
	mockLogger.ExpectDebug("set security header")

	// Call the function that uses the logger
	// Example: middleware.SetSecurityHeader(mockLogger)

	// Verify the expectations
	if err := mockLogger.Verify(); err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
}
