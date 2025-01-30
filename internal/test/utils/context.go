package utils

import (
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

// NewContext creates a new echo.Context for testing purposes.
func NewContext(req *http.Request) echo.Context {
	rec := httptest.NewRecorder()
	e := echo.New()
	return e.NewContext(req, rec)
}
