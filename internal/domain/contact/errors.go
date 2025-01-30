package contact

import "errors"

var (
	// ErrSubmissionNotFound is returned when a contact submission is not found
	ErrSubmissionNotFound = errors.New("contact submission not found")
	ErrNameRequired       = errors.New("name is required")
	ErrEmailRequired      = errors.New("email is required")
	ErrMessageRequired    = errors.New("message is required")
	ErrInvalidStatus      = errors.New("invalid status provided")
)
