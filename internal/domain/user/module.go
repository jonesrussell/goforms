package user

import (
	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories/database"
)

// Module provides the user dependencies
var Module = fx.Module("user",
	fx.Provide(
		NewStore,   // Provide the user repository
		NewService, // Provide the user service
	),
	fx.Invoke(func(logger logging.Logger) {
		logger.Debug("User module initialized")
	}),
)

// NewStore creates a new user repository
func NewStore(db *database.DB, logger logging.Logger) Repository {
	return &store{
		db:     db,
		logger: logger,
	}
}

// NewService creates a new user service
func NewService(repo Repository, tokenRepo TokenRepository, logger logging.Logger) *service {
	return &service{
		repo:      repo,
		tokenRepo: tokenRepo,
		logger:    logger,
	}
}
