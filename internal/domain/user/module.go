package user

import (
	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/application/repositories/database"

	"github.com/jonesrussell/goforms/internal/application/logging"
)

// Module provides the user dependencies
var Module = fx.Module("user",
	fx.Provide(
		NewStore,   // Provide the user repository
		NewService, // Provide the user service
	),
)

// NewStore creates a new user repository
func NewStore(db *database.DB, logger logging.Logger) UserRepository {
	return &store{
		db:     db,
		logger: logger,
	}
}

// NewService creates a new user service
func NewService(repo UserRepository, logger logging.Logger) Service {
	return &service{repo: repo, logger: logger}
}
