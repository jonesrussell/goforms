package user

import (
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/user"
)

type Service struct {
	repo      user.Repository
	tokenRepo user.TokenRepository
	logger    logging.Logger
}

// NewService creates a new user service
func NewService(repo user.Repository, tokenRepo user.TokenRepository, logger logging.Logger) *Service {
	return &Service{
		repo:      repo,
		tokenRepo: tokenRepo,
		logger:    logger,
	}
}
