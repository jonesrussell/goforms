package user

import (
	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories/database"
)

// Module provides the user-related dependencies
var Module = fx.Options(
	fx.Provide(NewStore),
	fx.Invoke(func(logger logging.Logger) {
		logger.Debug("User module initialized")
	}),
)

// NewStore creates a new user repository
func NewStore(db *database.DB, logger logging.Logger) Store {
	return &StoreImpl{
		db:     db,
		logger: logger,
	}
}
