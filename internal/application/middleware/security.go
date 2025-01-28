package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/jonesrussell/goforms/internal/application/logging"
)

// SecurityHeadersMiddleware adds security headers to all responses
func SecurityHeadersMiddleware(logger logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Debug("processing security headers")

			// Add security headers
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self';")
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response().Header().Set("Referrer-Policy", "no-referrer")
			c.Response().Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
			c.Response().Header().Set("Cross-Origin-Opener-Policy", "same-origin")
			c.Response().Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
			c.Response().Header().Set("Cross-Origin-Resource-Policy", "same-origin")

			logger.Debug("security headers processing complete")
			return next(c)
		}
	}
}
