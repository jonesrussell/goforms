package application

import (
	"github.com/labstack/echo/v4"

	"github.com/jonesrussell/goforms/internal/application/middleware"
	"github.com/jonesrussell/goforms/internal/application/validator"
	"github.com/jonesrussell/goforms/internal/infrastructure/logging"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(e *echo.Echo, handlers ...interface{ Register(e *echo.Echo) }) {
	// Register API handlers
	for _, handler := range handlers {
		handler.Register(e)
	}
}

// NewEcho creates a new Echo instance with common middleware and routes
func NewEcho(log logging.Logger) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Register validator
	log.Debug("creating validator instance")
	v := validator.NewValidator()
	if v == nil {
		log.Error("validator instance is nil")
		panic("validator instance is nil")
	}
	log.Debug("validator instance created")

	// Test validate a simple struct
	type test struct {
		Field string `validate:"required"`
	}
	err := v.Validate(test{})
	log.Debug("validator test", logging.Error(err))

	e.Validator = v
	log.Debug("validator registered with echo")

	// Setup middleware
	mw := middleware.New(log)
	mw.Setup(e)

	return e
}
