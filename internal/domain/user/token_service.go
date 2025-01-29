package user

import (
	"context"
	"fmt"

	"github.com/jonesrussell/goforms/internal/application/logging"
)

type tokenService struct {
	userRepo  Repository      // User repository for user-related operations
	tokenRepo TokenRepository // Token repository for token-related operations
	logger    logging.Logger
}

// Login authenticates the user and returns a pair of tokens
func (t *tokenService) Login(ctx context.Context, login *Login) (*TokenPair, error) {
	user, err := t.userRepo.GetByEmail(login.Email) // Use userRepo to get user by email
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

func NewTokenService(userRepo Repository, tokenRepo TokenRepository, logger logging.Logger) *tokenService {
	return &tokenService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		logger:    logger,
	}
}
