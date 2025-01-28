package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonesrussell/goforms/internal/application/logging"
)

var (
	// ErrUserNotFound indicates that a user was not found
	ErrUserNotFound = errors.New("user not found")
	// ErrEmailAlreadyExists indicates that a user with the given email already exists
	ErrEmailAlreadyExists = errors.New("email already exists")
	// ErrInvalidCredentials indicates that the provided credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInvalidToken indicates that the provided token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenBlacklisted indicates that the token has been blacklisted
	ErrTokenBlacklisted = errors.New("token is blacklisted")
)

// Service defines the methods for user-related operations
type Service interface {
	SignUp(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id uint) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context) ([]User, error)
}

// service implements the Service interface
type service struct {
	repo   UserRepository
	logger logging.Logger
}

// SignUp registers a new user
func (s *service) SignUp(ctx context.Context, user *User) error {
	// Check if email already exists
	existingUser, err := s.repo.GetByEmail(user.Email)
	if err != nil {
		s.logger.Error("failed to check existing user", logging.Error(err))
		return fmt.Errorf("failed to create user: %w", err)
	}
	if existingUser != nil {
		return fmt.Errorf("failed to create user: %w", ErrEmailAlreadyExists)
	}

	// Save user
	err = s.repo.Create(user)
	if err != nil {
		s.logger.Error("failed to create user", logging.Error(err))
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (s *service) GetUserByID(ctx context.Context, id uint) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to get user", logging.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpdateUser updates user information
func (s *service) UpdateUser(ctx context.Context, user *User) error {
	if err := s.repo.Update(user); err != nil {
		s.logger.Error("failed to update user", logging.Error(err))
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser removes a user
func (s *service) DeleteUser(ctx context.Context, id uint) error {
	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("failed to delete user", logging.Error(err))
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ListUsers returns all users
func (s *service) ListUsers(ctx context.Context) ([]User, error) {
	users, err := s.repo.List()
	if err != nil {
		s.logger.Error("failed to list users", logging.Error(err))
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}
