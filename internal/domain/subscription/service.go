package subscription

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jonesrussell/goforms/internal/infrastructure/logging"
)

// Service defines the interface for subscription operations
type Service interface {
	CreateSubscription(ctx context.Context, subscription *Subscription) error
	ListSubscriptions(ctx context.Context) ([]Subscription, error)
	GetSubscription(ctx context.Context, id int64) (*Subscription, error)
	GetSubscriptionByEmail(ctx context.Context, email string) (*Subscription, error)
	UpdateSubscriptionStatus(ctx context.Context, id int64, status Status) error
	DeleteSubscription(ctx context.Context, id int64) error
}

// ServiceImpl handles subscription business logic
type ServiceImpl struct {
	store  Store
	logger logging.Logger
}

// NewService creates a new subscription service
func NewService(store Store, logger logging.Logger) Service {
	return &ServiceImpl{
		store:  store,
		logger: logger,
	}
}

func (s *ServiceImpl) wrapError(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// CreateSubscription creates a new subscription
func (s *ServiceImpl) CreateSubscription(ctx context.Context, subscription *Subscription) error {
	// Validate subscription
	if err := subscription.Validate(); err != nil {
		return err
	}

	// Validate email format
	if !isValidEmail(subscription.Email) {
		return ErrInvalidEmail
	}

	// Check if email already exists
	existing, err := s.store.GetByEmail(ctx, subscription.Email)
	if err != nil && !errors.Is(err, ErrSubscriptionNotFound) {
		s.logger.Error("failed to check existing subscription", logging.Error(err))
		return s.wrapError(err, "failed to check existing subscription")
	}
	if existing != nil {
		return ErrEmailAlreadyExists
	}

	// Set default values
	subscription.Status = StatusPending
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()

	// Create subscription
	if err := s.store.Create(ctx, subscription); err != nil {
		s.logger.Error("failed to create subscription", logging.Error(err))
		return s.wrapError(err, "failed to create subscription")
	}

	return nil
}

// isValidEmail checks if the email format is valid
func isValidEmail(email string) bool {
	// Simple email validation for now
	// In a real application, you might want to use a more robust validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// ListSubscriptions returns all subscriptions
func (s *ServiceImpl) ListSubscriptions(ctx context.Context) ([]Subscription, error) {
	subscriptions, err := s.store.List(ctx)
	if err != nil {
		s.logger.Error("failed to list subscriptions", logging.Error(err))
		return nil, s.wrapError(err, "failed to list subscriptions")
	}

	return subscriptions, nil
}

// GetSubscription returns a subscription by ID
func (s *ServiceImpl) GetSubscription(ctx context.Context, id int64) (*Subscription, error) {
	subscription, err := s.store.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get subscription", logging.Error(err))
		return nil, s.wrapError(err, "failed to get subscription")
	}

	if subscription == nil {
		return nil, ErrSubscriptionNotFound
	}

	return subscription, nil
}

// GetSubscriptionByEmail returns a subscription by email
func (s *ServiceImpl) GetSubscriptionByEmail(ctx context.Context, email string) (*Subscription, error) {
	if email == "" {
		return nil, errors.New("invalid input: email is required")
	}

	subscription, err := s.store.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to get subscription by email", logging.Error(err))
		return nil, s.wrapError(err, "failed to get subscription by email")
	}

	if subscription == nil {
		return nil, ErrSubscriptionNotFound
	}

	return subscription, nil
}

// UpdateSubscriptionStatus updates the status of a subscription
func (s *ServiceImpl) UpdateSubscriptionStatus(ctx context.Context, id int64, status Status) error {
	// Validate status
	switch status {
	case StatusPending, StatusActive, StatusCancelled:
		// Valid status
	default:
		return ErrInvalidStatus
	}

	// Check if subscription exists
	subscription, err := s.store.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get subscription", logging.Error(err))
		return s.wrapError(err, "failed to get subscription")
	}

	if subscription == nil {
		return ErrSubscriptionNotFound
	}

	// Update status
	if err := s.store.UpdateStatus(ctx, id, status); err != nil {
		s.logger.Error("failed to update subscription status", logging.Error(err))
		return s.wrapError(err, "failed to update subscription status")
	}

	return nil
}

// DeleteSubscription removes a subscription
func (s *ServiceImpl) DeleteSubscription(ctx context.Context, id int64) error {
	// Check if subscription exists
	subscription, err := s.store.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get subscription", logging.Error(err))
		return s.wrapError(err, "failed to get subscription")
	}

	if subscription == nil {
		return ErrSubscriptionNotFound
	}

	// Delete subscription
	if err := s.store.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete subscription", logging.Error(err))
		return s.wrapError(err, "failed to delete subscription")
	}

	return nil
}
