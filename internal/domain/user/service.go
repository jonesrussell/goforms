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
	SignUp(ctx context.Context, signup *Signup) (*User, error)
	GetUserByID(ctx context.Context, id uint) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context) ([]User, error)
	GetByEmail(email string) (*User, error)
	IsTokenBlacklisted(token string) bool
	ValidateToken(token string) (string, error)
	Login(ctx context.Context, login *Login) (*TokenPair, error)
	Logout(ctx context.Context, token string) error
}

// service implements the Service interface
type service struct {
	repo   Repository
	logger logging.Logger
}

// SignUp registers a new user and returns the created user
func (s *service) SignUp(ctx context.Context, signup *Signup) (*User, error) {
	// Check if email already exists
	existingUser, err := s.repo.GetByEmail(signup.Email)
	if err != nil {
		s.logger.Error("failed to check existing user", logging.Error(err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("failed to create user: %w", ErrEmailAlreadyExists)
	}

	// Create a new User instance
	user := &User{
		Email:     signup.Email,
		FirstName: signup.FirstName,
		LastName:  signup.LastName,
		Role:      "user", // Set a default role or modify as needed
		Active:    true,   // Set default active status
	}

	// Hash the password
	if err := user.SetPassword(signup.Password); err != nil {
		s.logger.Error("failed to set password", logging.Error(err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Save user
	err = s.repo.Create(user) // Pass the User instance to the repository
	if err != nil {
		s.logger.Error("failed to create user", logging.Error(err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil // Return the created User instance
}

// GetUserByID retrieves a user by ID
func (s *service) GetUserByID(ctx context.Context, id uint) (*User, error) {
	user, err := s.repo.Get(id)
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

// IsTokenBlacklisted checks if the provided token is blacklisted
func (s *service) IsTokenBlacklisted(token string) bool {
	// Implement your logic to check if the token is blacklisted
	// For example, check against a database or in-memory store
	return false // Placeholder return value
}

// ValidateToken checks if the provided token is valid
func (s *service) ValidateToken(token string) (string, error) {
	// Implement your logic to validate the token
	// For example, decode the token and check its validity
	return token, nil // Placeholder return value
}

// Login authenticates a user and returns a token pair
func (s *service) Login(ctx context.Context, login *Login) (*TokenPair, error) {
	// Implement your logic to authenticate the user
	// For example, check the credentials and generate tokens
	return nil, nil // Placeholder return value
}

// Logout invalidates the user's token
func (s *service) Logout(ctx context.Context, token string) error {
	// Implement your logic to invalidate the token
	// For example, mark the token as blacklisted in the database
	return nil // Placeholder return value
}

// Implement the GetByEmail method
func (s *service) GetByEmail(email string) (*User, error) {
	return s.repo.GetByEmail(email) // Assuming repo is your data access layer
}
