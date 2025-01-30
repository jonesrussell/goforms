package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/common"
)

// Define necessary errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenBlacklisted   = errors.New("token is blacklisted")
)

// TokenService defines the methods for token-related operations
type TokenService interface {
	IsTokenBlacklisted(token string) bool
	ValidateToken(token string) (string, error)
	Login(ctx context.Context, login *Login) (*TokenPair, error)
	Logout(ctx context.Context, token string) error
}

// Service defines the methods for user services
type Service interface {
	DeleteUser(ctx context.Context, id uint) error
	GetByEmail(email string) (*common.User, error)
	GetUserByID(ctx context.Context, id uint) (*common.User, error)
	ListUsers(ctx context.Context) ([]common.User, error)
	Login(ctx context.Context, login *Login) (*TokenPair, error)
	Logout(ctx context.Context, token string) error
	SignUp(ctx context.Context, signup *Signup) (*common.User, error)
	UpdateSubmissionStatus(ctx context.Context, id int64, status string) error
	UpdateUser(ctx context.Context, u *common.User) error
	IsTokenBlacklisted(token string) bool
}

// ServiceImpl implements the Service interface
type ServiceImpl struct {
	repo      Repository      // User repository for user-related operations
	tokenRepo TokenRepository // Token repository for token-related operations
	logger    logging.Logger
}

// NewService creates a new user service
func NewService(repo Repository, tokenRepo TokenRepository, logger logging.Logger) Service {
	return &ServiceImpl{
		repo:      repo,
		tokenRepo: tokenRepo,
		logger:    logger,
	}
}

// SignUp handles user registration
func (s *ServiceImpl) SignUp(ctx context.Context, signup *Signup) (*common.User, error) {
	// Check if email already exists
	existingUser, err := s.repo.GetByEmail(signup.Email)
	if err != nil {
		s.logger.Error("Error checking existing user", logging.Error(err))
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Create user
	newUser, err := s.repo.SignUp(ctx, signup)
	if err != nil {
		s.logger.Error("Error signing up user", logging.Error(err))
		return nil, fmt.Errorf("error signing up user: %w", err)
	}

	return newUser, nil
}

// DeleteUser deletes a user by ID
func (s *ServiceImpl) DeleteUser(ctx context.Context, id uint) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		s.logger.Error("Error deleting user", logging.Error(err))
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}

// GetByEmail retrieves a user by email
func (s *ServiceImpl) GetByEmail(email string) (*common.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		s.logger.Error("Error retrieving user by email", logging.Error(err))
		return nil, fmt.Errorf("error retrieving user by email: %w", err)
	}
	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *ServiceImpl) GetUserByID(ctx context.Context, id uint) (*common.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		s.logger.Error("Error retrieving user by ID", logging.Error(err))
		return nil, fmt.Errorf("error retrieving user by ID: %w", err)
	}
	return user, nil
}

// ListUsers lists all users
func (s *ServiceImpl) ListUsers(ctx context.Context) ([]common.User, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		s.logger.Error("Error listing users", logging.Error(err))
		return nil, fmt.Errorf("error listing users: %w", err)
	}
	return users, nil
}

// Login handles user authentication
func (s *ServiceImpl) Login(ctx context.Context, login *Login) (*TokenPair, error) {
	// Authenticate user
	user, err := s.repo.GetByEmail(login.Email)
	if err != nil {
		s.logger.Error("Error fetching user for login", logging.Error(err))
		return nil, fmt.Errorf("error fetching user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := verifyPassword(user.HashedPassword, login.Password); err != nil {
		s.logger.Warn("Invalid password attempt", logging.String("email", login.Email))
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	tokenPair, err := s.tokenRepo.GenerateTokens(user)
	if err != nil {
		s.logger.Error("Error generating tokens", logging.Error(err))
		return nil, fmt.Errorf("error generating tokens: %w", err)
	}

	return tokenPair, nil
}

// Logout handles user logout
func (s *ServiceImpl) Logout(ctx context.Context, token string) error {
	err := s.tokenRepo.InvalidateToken(token)
	if err != nil {
		s.logger.Error("Error invalidating token", logging.Error(err))
		return fmt.Errorf("error invalidating token: %w", err)
	}
	return nil
}

// UpdateSubmissionStatus updates the status of a submission
func (s *ServiceImpl) UpdateSubmissionStatus(ctx context.Context, id int64, status string) error {
	err := s.repo.UpdateSubmissionStatus(ctx, id, status)
	if err != nil {
		s.logger.Error("Error updating submission status", logging.Error(err))
		return fmt.Errorf("error updating submission status: %w", err)
	}
	return nil
}

// UpdateUser updates user information
func (s *ServiceImpl) UpdateUser(ctx context.Context, u *common.User) error {
	err := s.repo.UpdateUser(ctx, u)
	if err != nil {
		s.logger.Error("Error updating user", logging.Error(err))
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

// IsTokenBlacklisted checks if a token is blacklisted
func (s *ServiceImpl) IsTokenBlacklisted(token string) bool {
	return s.repo.IsTokenBlacklisted(token)
}

// Helper functions
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
