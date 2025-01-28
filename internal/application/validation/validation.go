package validation

import "github.com/go-playground/validator/v10"

// New creates a new validator instance
func New() *validator.Validate {
	return validator.New()
}
