package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestSignupHandler(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/signup", nil) // Assuming signup is a POST request
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := SignupHandler(c); err != nil {
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
