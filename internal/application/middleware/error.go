package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorHandler is a custom error handler middleware
func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			// Log the error
			c.Logger().Error(err)

			// Return a generic error response
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
		}
		return nil
	}
}
