package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/user"
)

// Config holds middleware configuration
type Config struct {
	Logger      logging.Logger
	JWTSecret   string
	UserService user.Service
	EnableCSRF  bool
}

// Setup configures all middleware for an Echo instance
func Setup(e *echo.Echo, cfg *Config) {
	// Security
	m := New(cfg.Logger)
	m.Setup(e)

	// Logging
	e.Use(LoggingMiddleware(cfg.Logger))

	// Auth if secret provided
	if cfg.JWTSecret != "" && cfg.UserService != nil {
		userService := cfg.UserService
		e.Use(NewJWTMiddleware(userService, cfg.JWTSecret))
	}

	// CSRF if enabled
	if cfg.EnableCSRF {
		e.Use(CSRFMiddleware())
	}
}
