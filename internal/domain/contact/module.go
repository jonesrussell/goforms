package contact

import (
	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories/database"
)

// Module provides the contact dependencies
var Module = fx.Module("contact",
	fx.Provide(
		NewStore,   // Provide the contact repository
		NewService, // Provide the contact service
	),
)

// NewStore creates a new contact store
func NewStore(db *database.DB, logger logging.Logger) Store {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// NewService creates a new contact service
func NewService(store Store, logger logging.Logger) Service {
	return &ServiceImpl{
		store:  store,
		logger: logger,
	}
}
