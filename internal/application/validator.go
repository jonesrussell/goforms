package application

import (
	"github.com/go-playground/validator/v10"

	"github.com/jonesrussell/goforms/internal/application/validation"
)

// CustomValidator for request validation
type CustomValidator struct {
	validator *validator.Validate
}

// Validate implements echo.Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// NewValidator creates a new validator instance
func NewValidator() *CustomValidator {
	return &CustomValidator{
		validator: validation.New(),
	}
}
