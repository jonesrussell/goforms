package subscription

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Status represents the status of a subscription
type Status string

const (
	// StatusPending indicates a pending subscription
	StatusPending Status = "pending"
	// StatusActive indicates an active subscription
	StatusActive Status = "active"
	// StatusCancelled indicates a cancelled subscription
	StatusCancelled Status = "cancelled"
)

// Subscription represents a newsletter subscription
type Subscription struct {
	ID        int64     `json:"id" db:"id"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Name      string    `json:"name" db:"name" validate:"required"`
	Status    Status    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Validate validates the subscription
func (s *Subscription) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
