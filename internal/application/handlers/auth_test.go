package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/user"
	"github.com/jonesrussell/goforms/test/utils"
	"github.com/labstack/echo/v4"
)

func TestAuthHandler_handleSignup(t *testing.T) {
	e := echo.New()
	mockLogger := &utils.MockLogger{
		DebugFunc: func(msg string, fields ...interface{}) {},
		ErrorFunc: func(msg string, fields ...logging.Field) {},
		InfoFunc:  func(msg string, fields ...logging.Field) {},
		WarnFunc:  func(msg string, fields ...logging.Field) {},
	}

	mockUserService := &user.MockService{} // Use the mock user service

	handler := NewAuthHandler(mockLogger, mockUserService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", nil) // Assuming signup is a POST request
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.handleSignup(c); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	expected := "Signup successful"
	if rec.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, rec.Body.String())
	}
}

func TestAuthHandler_handleLogin(t *testing.T) {
	e := echo.New()
	mockLogger := &utils.MockLogger{
		DebugFunc: func(msg string, fields ...interface{}) {},
		ErrorFunc: func(msg string, fields ...logging.Field) {},
		InfoFunc:  func(msg string, fields ...logging.Field) {},
		WarnFunc:  func(msg string, fields ...logging.Field) {},
	}

	userService := &user.MockService{} // Use a mock user service

	handler := NewAuthHandler(mockLogger, userService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil) // Assuming login is a POST request
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.handleLogin(c); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestAuthHandler_handleLogout(t *testing.T) {
	e := echo.New()
	mockLogger := &utils.MockLogger{
		DebugFunc: func(msg string, fields ...interface{}) {},
		ErrorFunc: func(msg string, fields ...logging.Field) {},
		InfoFunc:  func(msg string, fields ...logging.Field) {},
		WarnFunc:  func(msg string, fields ...logging.Field) {},
	}

	userService := &user.MockService{} // Use a mock user service

	handler := NewAuthHandler(mockLogger, userService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil) // Assuming logout is a POST request
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set the Authorization header for the logout request
	req.Header.Set("Authorization", "Bearer mocktoken")

	if err := handler.handleLogout(c); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
