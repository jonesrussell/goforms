package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/contact"
	"github.com/jonesrussell/goforms/internal/presentation/view"
	"github.com/jonesrussell/goforms/test/utils"
)

func TestWebHandler_handleHome(t *testing.T) {
	e := echo.New()
	mockLogger := &utils.MockLogger{
		DebugFunc: func(msg string, fields ...interface{}) {},
		ErrorFunc: func(msg string, fields ...logging.Field) {},
		InfoFunc:  func(msg string, fields ...logging.Field) {},
		WarnFunc:  func(msg string, fields ...logging.Field) {},
	}

	renderer := &view.Renderer{}             // Initialize your view renderer here
	contactService := &contact.MockService{} // Use the mock service

	handler := NewWebHandler(mockLogger, WithRenderer(renderer), WithContactService(contactService))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.handleHome(c); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestWebHandler_handleDemo(t *testing.T) {
	e := echo.New()
	mockLogger := &utils.MockLogger{
		DebugFunc: func(msg string, fields ...interface{}) {},
		ErrorFunc: func(msg string, fields ...logging.Field) {},
		InfoFunc:  func(msg string, fields ...logging.Field) {},
		WarnFunc:  func(msg string, fields ...logging.Field) {},
	}

	renderer := &view.Renderer{}             // Initialize your view renderer here
	contactService := &contact.MockService{} // Mock contact service if needed

	handler := NewWebHandler(mockLogger, WithRenderer(renderer), WithContactService(contactService))

	req := httptest.NewRequest(http.MethodGet, "/demo", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.handleDemo(c); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestWebHandler_handleSignup(t *testing.T) {
	e := echo.New()
	mockLogger := &utils.MockLogger{
		DebugFunc: func(msg string, fields ...interface{}) {},
		ErrorFunc: func(msg string, fields ...logging.Field) {},
		InfoFunc:  func(msg string, fields ...logging.Field) {},
		WarnFunc:  func(msg string, fields ...logging.Field) {},
	}

	renderer := &view.Renderer{}             // Initialize your view renderer here
	contactService := &contact.MockService{} // Mock contact service if needed

	handler := NewWebHandler(mockLogger, WithRenderer(renderer), WithContactService(contactService))

	req := httptest.NewRequest(http.MethodGet, "/signup", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.handleSignup(c); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestWebHandler_handleLogin(t *testing.T) {
	e := echo.New()
	mockLogger := &utils.MockLogger{
		DebugFunc: func(msg string, fields ...interface{}) {},
		ErrorFunc: func(msg string, fields ...logging.Field) {},
		InfoFunc:  func(msg string, fields ...logging.Field) {},
		WarnFunc:  func(msg string, fields ...logging.Field) {},
	}

	renderer := &view.Renderer{}             // Initialize your view renderer here
	contactService := &contact.MockService{} // Mock contact service if needed

	handler := NewWebHandler(mockLogger, WithRenderer(renderer), WithContactService(contactService))

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.handleLogin(c); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
