package interfaces

import "github.com/go-playground/validator/v10"

// Validator defines the interface for validation operations
type Validator interface {
	// Struct validates a struct based on validation tags
	Struct(any) error
	// Var validates a single variable using a tag
	Var(any, string) error
	// RegisterValidation adds a custom validation with the given tag
	RegisterValidation(string, func(fl validator.FieldLevel) bool) error
}
