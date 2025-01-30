package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go" // Ensure this is present

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/common"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenBlacklisted   = errors.New("token is blacklisted")
)

// Service defines the methods for user services
type Service interface {
	DeleteUser(ctx context.Context, id uint) error
	GetByEmail(email string) (*common.User, error)
	GetUserByID(ctx context.Context, id uint) (*common.User, error)
	ListUsers(ctx context.Context) ([]common.User, error)
	Login(ctx context.Context, login *Login) (*TokenPair, error)
	Logout(ctx context.Context, token string) error
	SignUp(signup *Signup) (*common.User, error)
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

func (s *ServiceImpl) SignUp(signup *Signup) (*common.User, error) {
	u := ConvertSignupToUser(signup) // Convert Signup to User

	// Log the user details before saving (excluding the password)
	s.logger.Debug("Preparing to create user", logging.Any("user", u))

	// Hash the password before saving
	if err := u.SetPassword(signup.Password); err != nil {
		return nil, err // Return error if password hashing fails
	}

	err := s.repo.Create(u) // Save the user to the database
	if err != nil {
		s.logger.Error("Failed to create user", logging.Error(err))
		return nil, err // Return the error if creation fails
	}

	s.logger.Debug("User created successfully", logging.Any("user", u))
	return u, nil // Return the created user
}

func (s *ServiceImpl) GetUserByID(ctx context.Context, id uint) (*common.User, error) {
	u, err := s.repo.Get(id)
	if err != nil {
		return nil, logAndWrapError(s.logger, "failed to get user", err)
	}
	return u, nil
}

func (s *ServiceImpl) UpdateUser(ctx context.Context, u *common.User) error {
	if err := s.repo.Update(u); err != nil {
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

func (s *ServiceImpl) ListUsers(ctx context.Context) ([]common.User, error) {
	users, err := s.repo.List()
	if err != nil {
		return nil, logAndWrapError(s.logger, "failed to list users", err)
	}
	return users, nil
}

func (s *ServiceImpl) GetByEmail(email string) (*common.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, logAndWrapError(s.logger, "failed to get user by email", err)
	}
	return user, nil
}

func (s *ServiceImpl) Login(ctx context.Context, login *Login) (*TokenPair, error) {
	u, err := s.repo.GetByEmail(login.Email)
	if err != nil {
		return nil, err // Handle error
	}
	if u == nil || !u.CheckPassword(login.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Proceed with token generation
	return s.generateTokens(u) // Call the generateTokens function
}

// generateTokens generates access and refresh tokens for the user
func (s *ServiceImpl) generateTokens(u *common.User) (*TokenPair, error) {
	// Example secret key for signing tokens (use a secure method to manage secrets)
	secretKey := []byte("your_secret_key")

	// Generate Access Token
	accessTokenClaims := jwt.MapClaims{
		"email": u.Email,
		"role":  u.Role,
		"exp":   time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	signedAccessToken, err := accessToken.SignedString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate Refresh Token
	refreshTokenClaims := jwt.MapClaims{
		"email": u.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(), // Token expires in 7 days
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedRefreshToken, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
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

// Implement the methods defined in the Service interface
func (s *ServiceImpl) SomeMethod() error {
	// Implementation here
	return nil
}

// NewService creates a new user service
func NewService(repo Repository, tokenRepo TokenRepository, logger logging.Logger) Service {
	return &ServiceImpl{
		repo:      repo,
		tokenRepo: tokenRepo,
		logger:    logger,
	}
}
