package contact

import (
	"time"
)

// Status represents the status of a contact form submission
type Status string

const (
	// StatusPending indicates a pending submission
	StatusPending Status = "pending"
	// StatusApproved indicates an approved submission
	StatusApproved Status = "approved"
	// StatusRejected indicates a rejected submission
	StatusRejected Status = "rejected"
)

// Submission represents a contact form submission
type Submission struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Message   string    `json:"message" db:"message"`
	Status    Status    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
