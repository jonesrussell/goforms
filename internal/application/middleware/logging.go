package middleware

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/jonesrussell/goforms/internal/application/logging"
)

// LoggingMiddleware logs incoming requests
func LoggingMiddleware(logger logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Call the next handler
			err := next(c)

			// Set status based on error
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					c.Response().Status = he.Code
				} else {
					c.Response().Status = echo.ErrInternalServerError.Code
				}
			}

			// Log the request details
			logger.Info("http request",
				logging.String("method", c.Request().Method),
				logging.String("path", c.Request().URL.Path),
				logging.Int("status", c.Response().Status),
				logging.Duration("latency", time.Since(start)),
				logging.String("ip", c.RealIP()),
			)

			return err
		}
	}
}
