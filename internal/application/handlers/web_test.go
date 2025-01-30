package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/contact"
	"github.com/jonesrussell/goforms/internal/presentation/view"
	"github.com/jonesrussell/goforms/internal/test/utils"
)

// Helper function to create a new request with JSON body
func newRequest(method, path string, body interface{}) (*http.Request, error) {
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Helper function to create a new context
func newContext(req *http.Request) echo.Context {
	rec := httptest.NewRecorder()
	e := echo.New()
	return e.NewContext(req, rec)
}

func TestWebHandler_handleHome(t *testing.T) {
	mockLogger := &utils.MockLogger{
		DebugFunc: func(msg string, fields ...interface{}) {
			fmt.Printf("DEBUG: %s %v\n", msg, fields)
		},
		ErrorFunc: func(msg string, fields ...logging.Field) {
			fmt.Printf("ERROR: %s %v\n", msg, fields)
		},
		SyncFunc: func() error {
			return nil // Mock Sync behavior
		},
	}
	renderer := &view.Renderer{}             // Initialize your view renderer here
	contactService := &contact.MockService{} // Use the mock service

	handler := NewWebHandler(mockLogger, WithRenderer(renderer), WithContactService(contactService))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := newContext(req)

	if err := handler.handleHome(c); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestWebHandler_handleDemo(t *testing.T) {
	mockLogger := &utils.MockLogger{
		DebugFunc: func(msg string, fields ...interface{}) {
			fmt.Printf("DEBUG: %s %v\n", msg, fields)
		},
		ErrorFunc: func(msg string, fields ...logging.Field) {
			fmt.Printf("ERROR: %s %v\n", msg, fields)
		},
		SyncFunc: func() error {
			return nil // Mock Sync behavior
		},
	}

	renderer := &view.Renderer{}             // Initialize your view renderer here
	contactService := &contact.MockService{} // Mock contact service if needed

	handler := NewWebHandler(mockLogger, WithRenderer(renderer), WithContactService(contactService))

	req := httptest.NewRequest(http.MethodGet, "/demo", nil)
	rec := httptest.NewRecorder()
	c := newContext(req)

	if err := handler.handleDemo(c); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestWebHandler_handleSignup(t *testing.T) {
	mockLogger := &utils.MockLogger{
		DebugFunc: func(msg string, fields ...interface{}) {
			fmt.Printf("DEBUG: %s %v\n", msg, fields)
		},
		ErrorFunc: func(msg string, fields ...logging.Field) {
			fmt.Printf("ERROR: %s %v\n", msg, fields)
		},
		SyncFunc: func() error {
			return nil // Mock Sync behavior
		},
	}

	renderer := &view.Renderer{}             // Initialize your view renderer here
	contactService := &contact.MockService{} // Mock contact service if needed

	handler := NewWebHandler(mockLogger, WithRenderer(renderer), WithContactService(contactService))

	req := httptest.NewRequest(http.MethodGet, "/signup", nil)
	rec := httptest.NewRecorder()
	c := newContext(req)

	if err := handler.handleSignup(c); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestWebHandler_handleLogin(t *testing.T) {
	mockLogger := &utils.MockLogger{
		DebugFunc: func(msg string, fields ...interface{}) {
			fmt.Printf("DEBUG: %s %v\n", msg, fields)
		},
		ErrorFunc: func(msg string, fields ...logging.Field) {
			fmt.Printf("ERROR: %s %v\n", msg, fields)
		},
		SyncFunc: func() error {
			return nil // Mock Sync behavior
		},
	}

	renderer := &view.Renderer{}             // Initialize your view renderer here
	contactService := &contact.MockService{} // Mock contact service if needed

	handler := NewWebHandler(mockLogger, WithRenderer(renderer), WithContactService(contactService))

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()
	c := newContext(req)

	if err := handler.handleLogin(c); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
