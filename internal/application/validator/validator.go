package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator wraps the validator.Validate instance
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

// Validate validates the provided struct using the validator instance
func (cv *CustomValidator) Validate(i interface{}) error {
	if cv.validator == nil {
		return fmt.Errorf("validator not initialized")
	}
	return cv.validator.Struct(i)
}

// Ensure CustomValidator implements echo.Validator at compile time
var _ echo.Validator = &CustomValidator{}
