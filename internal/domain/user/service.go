package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonesrussell/goforms/internal/application/logging"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenBlacklisted   = errors.New("token is blacklisted")
)

// Service defines the methods for user-related operations
type Service interface {
	SignUp(ctx context.Context, signup *Signup) (*User, error)
	GetUserByID(ctx context.Context, id uint) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context) ([]User, error)
	GetByEmail(email string) (*User, error)
	Login(ctx context.Context, login *Login) (*TokenPair, error)
	Logout(ctx context.Context, token string) error
	IsTokenBlacklisted(token string) bool
}

// TokenService defines the methods for token-related operations
type TokenService interface {
	IsTokenBlacklisted(token string) bool
	ValidateToken(token string) (string, error)
	Login(ctx context.Context, login *Login) (*TokenPair, error)
	Logout(ctx context.Context, token string) error
}

type service struct {
	repo      Repository      // User repository for user-related operations
	tokenRepo TokenRepository // Token repository for token-related operations
	logger    logging.Logger
}

// Helper function for error logging and wrapping
func logAndWrapError(logger logging.Logger, msg string, err error) error {
	logger.Error(msg, logging.Error(err))
	return fmt.Errorf("%s: %w", msg, err)
}

func (s *service) SignUp(ctx context.Context, signup *Signup) (*User, error) {
	existingUser, err := s.repo.GetByEmail(signup.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("failed to create user: %w", ErrEmailAlreadyExists)
	}

	user := &User{
		Email:     signup.Email,
		FirstName: signup.FirstName,
		LastName:  signup.LastName,
		Role:      "user",
		Active:    true,
	}

	if err := user.SetPassword(signup.Password); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *service) GetUserByID(ctx context.Context, id uint) (*User, error) {
	user, err := s.repo.Get(id)
	if err != nil {
		return nil, logAndWrapError(s.logger, "failed to get user", err)
	}
	return user, nil
}

func (s *service) UpdateUser(ctx context.Context, user *User) error {
	if err := s.repo.Update(user); err != nil {
		return logAndWrapError(s.logger, "failed to update user", err)
	}
	return nil
}

func (s *service) DeleteUser(ctx context.Context, id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return logAndWrapError(s.logger, "failed to delete user", err)
	}
	return nil
}

func (s *service) ListUsers(ctx context.Context) ([]User, error) {
	users, err := s.repo.List()
	if err != nil {
		return nil, logAndWrapError(s.logger, "failed to list users", err)
	}
	return users, nil
}

func (s *service) GetByEmail(email string) (*User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, logAndWrapError(s.logger, "failed to get user by email", err)
	}
	return user, nil
}

func (s *service) Login(ctx context.Context, login *Login) (*TokenPair, error) {
	user, err := s.repo.GetByEmail(login.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate user: %w", err)
	}
	if user == nil || !user.CheckPassword(login.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate tokens (this is a placeholder; implement your token generation logic)
	tokens := &TokenPair{
		AccessToken:  "generated-access-token",  // Replace with actual token generation
		RefreshToken: "generated-refresh-token", // Replace with actual token generation
	}

	return tokens, nil
}

// Logout invalidates the user's token
func (s *service) Logout(ctx context.Context, token string) error {
	err := s.tokenRepo.BlacklistToken(token)
	if err != nil {
		return fmt.Errorf("failed to logout user: %w", err)
	}
	return nil
}

func (s *service) IsTokenBlacklisted(token string) bool {
	// Implement your logic to check if the token is blacklisted
	return s.tokenRepo.IsTokenBlacklisted(token) // Assuming this method exists in your TokenRepository
}
