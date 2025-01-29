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
	DeleteUser(ctx context.Context, id uint) error
	GetByEmail(email string) (*User, error)
	GetUserByID(ctx context.Context, id uint) (*User, error)
	ListUsers(ctx context.Context) ([]User, error)
	Login(ctx context.Context, login *Login) (*TokenPair, error)
	Logout(ctx context.Context, token string) error
	SignUp(signup *Signup) (*User, error)
	UpdateSubmissionStatus(ctx context.Context, id int64, status string) error
	UpdateUser(ctx context.Context, user *User) error
	IsTokenBlacklisted(token string) bool
}

// Service defines the methods for user-related operations
type ServiceImpl struct {
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

func (s *ServiceImpl) SignUp(signup *Signup) (*User, error) {
	user := ConvertSignupToUser(signup) // Convert Signup to User

	// Log the user details before saving (excluding the password)
	s.logger.Debug("Preparing to create user", logging.Any("user", user))

	// Hash the password before saving
	if err := user.SetPassword(signup.Password); err != nil {
		return nil, err // Return error if password hashing fails
	}

	err := s.repo.Create(user) // Save the user to the database
	if err != nil {
		s.logger.Error("Failed to create user", logging.Error(err))
		return nil, err // Return the error if creation fails
	}

	s.logger.Debug("User created successfully", logging.Any("user", user))
	return user, nil // Return the created user and nil error
}

func (s *ServiceImpl) GetUserByID(ctx context.Context, id uint) (*User, error) {
	user, err := s.repo.Get(id)
	if err != nil {
		return nil, logAndWrapError(s.logger, "failed to get user", err)
	}
	return user, nil
}

func (s *ServiceImpl) UpdateUser(ctx context.Context, user *User) error {
	if err := s.repo.Update(user); err != nil {
		return logAndWrapError(s.logger, "failed to update user", err)
	}
	return nil
}

func (s *ServiceImpl) DeleteUser(ctx context.Context, id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return logAndWrapError(s.logger, "failed to delete user", err)
	}
	return nil
}

func (s *ServiceImpl) ListUsers(ctx context.Context) ([]User, error) {
	users, err := s.repo.List()
	if err != nil {
		return nil, logAndWrapError(s.logger, "failed to list users", err)
	}
	return users, nil
}

func (s *ServiceImpl) GetByEmail(email string) (*User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, logAndWrapError(s.logger, "failed to get user by email", err)
	}
	return user, nil
}

func (s *ServiceImpl) Login(ctx context.Context, login *Login) (*TokenPair, error) {
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
func (s *ServiceImpl) Logout(ctx context.Context, token string) error {
	err := s.tokenRepo.BlacklistToken(token)
	if err != nil {
		return fmt.Errorf("failed to logout user: %w", err)
	}
	return nil
}

func (s *ServiceImpl) IsTokenBlacklisted(token string) bool {
	// Implement your logic to check if the token is blacklisted
	return s.tokenRepo.IsTokenBlacklisted(token) // Assuming this method exists in your TokenRepository
}

func (s *ServiceImpl) UpdateSubmissionStatus(ctx context.Context, id int64, status string) error {
	// Implement your logic to update the submission status
	// For example, you might want to update a submission in the database
	return nil // Return nil for now, or implement actual logic
}

// NewService creates a new user service
func NewService(repo Repository, tokenRepo TokenRepository, logger logging.Logger) *ServiceImpl {
	return &ServiceImpl{
		repo:      repo,
		tokenRepo: tokenRepo,
		logger:    logger,
	}
}
