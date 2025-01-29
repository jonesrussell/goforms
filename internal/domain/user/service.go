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
type Service struct {
	repo      Repository      // User repository for user-related operations
	tokenRepo TokenRepository // Token repository for token-related operations
	logger    logging.Logger
}

// TokenService defines the methods for token-related operations
type TokenService interface {
	IsTokenBlacklisted(token string) bool
	ValidateToken(token string) (string, error)
	Login(ctx context.Context, login *Login) (*TokenPair, error)
	Logout(ctx context.Context, token string) error
}

// Helper function for error logging and wrapping
func logAndWrapError(logger logging.Logger, msg string, err error) error {
	logger.Error(msg, logging.Error(err))
	return fmt.Errorf("%s: %w", msg, err)
}

func (s *Service) SignUp(email, password, firstName string) (*User, error) {
	user := &User{
		Email:          email,
		FirstName:      firstName,
		HashedPassword: password, // Ensure you are storing the password correctly
	}

	err := s.repo.Create(user)
	if err != nil {
		return nil, err // Return the error if creation fails
	}
	return user, nil // Return the created user and nil error
}

func (s *Service) GetUserByID(ctx context.Context, id uint) (*User, error) {
	user, err := s.repo.Get(id)
	if err != nil {
		return nil, logAndWrapError(s.logger, "failed to get user", err)
	}
	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, user *User) error {
	if err := s.repo.Update(user); err != nil {
		return logAndWrapError(s.logger, "failed to update user", err)
	}
	return nil
}

func (s *Service) DeleteUser(ctx context.Context, id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return logAndWrapError(s.logger, "failed to delete user", err)
	}
	return nil
}

func (s *Service) ListUsers(ctx context.Context) ([]User, error) {
	users, err := s.repo.List()
	if err != nil {
		return nil, logAndWrapError(s.logger, "failed to list users", err)
	}
	return users, nil
}

func (s *Service) GetByEmail(email string) (*User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, logAndWrapError(s.logger, "failed to get user by email", err)
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, login *Login) (*TokenPair, error) {
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
func (s *Service) Logout(ctx context.Context, token string) error {
	err := s.tokenRepo.BlacklistToken(token)
	if err != nil {
		return fmt.Errorf("failed to logout user: %w", err)
	}
	return nil
}

func (s *Service) IsTokenBlacklisted(token string) bool {
	// Implement your logic to check if the token is blacklisted
	return s.tokenRepo.IsTokenBlacklisted(token) // Assuming this method exists in your TokenRepository
}

// NewService creates a new user service
func NewService(repo Repository, tokenRepo TokenRepository, logger logging.Logger) *Service {
	return &Service{
		repo:      repo,
		tokenRepo: tokenRepo,
		logger:    logger,
	}
}
