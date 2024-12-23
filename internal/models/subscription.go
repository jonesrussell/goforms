package models

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// Subscription represents a newsletter subscription
type Subscription struct {
	ID        int64     `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Validate checks if the subscription data is valid
func (s *Subscription) Validate() error {
	if s.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email is required")
	}

	// Check for spaces in email
	if strings.Contains(s.Email, " ") {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid email format")
	}

	// Basic email format validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(s.Email) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid email format")
	}

	return nil
}

// SubscriptionStore defines the interface for subscription storage operations
type SubscriptionStore interface {
	CreateSubscription(ctx context.Context, sub *Subscription) error
}

// subscriptionStore implements SubscriptionStore
type subscriptionStore struct {
	db DB
}

// NewSubscriptionStore creates a new subscription store
func NewSubscriptionStore(db DB) SubscriptionStore {
	return &subscriptionStore{db: db}
}

// CreateSubscription implements the subscription creation
func (s *subscriptionStore) CreateSubscription(ctx context.Context, sub *Subscription) error {
	// Create a timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check for existing subscription
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM subscriptions WHERE email = ?)"
	err := s.db.QueryRowxContext(ctx, checkQuery, sub.Email).Scan(&exists)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check subscription")
	}
	if exists {
		return echo.NewHTTPError(http.StatusConflict, "Email already subscribed")
	}

	query := `
		INSERT INTO subscriptions (email, created_at)
		VALUES (?, ?)
		RETURNING id`

	sub.CreatedAt = time.Now()
	row := s.db.QueryRowxContext(ctx, query,
		sub.Email,
		sub.CreatedAt,
	)

	err = row.Scan(&sub.ID)
	if err != nil {
		if err == context.DeadlineExceeded {
			return echo.NewHTTPError(http.StatusGatewayTimeout, "Database operation timed out")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create subscription")
	}

	return nil
}
