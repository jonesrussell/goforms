package user

import (
	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/application/repositories/database"

	"github.com/jonesrussell/goforms/internal/application/logging"
)

// Module provides the user dependencies
var Module = fx.Module("user",
	fx.Provide(
		NewStore,   // Ensure this is using the correct type
		NewService, // Provide the user service
	),
)

// NewStore creates a new user repository
func NewStore(db *database.DB, logger logging.Logger) UserRepository {
	return &store{ // Ensure you are returning an instance of store
		db:     db,
		logger: logger,
	}
}

// NewService creates a new user service
func NewService(repo UserRepository) Service {
	return &service{repo: repo}
}
