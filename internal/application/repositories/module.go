package repositories

import (
	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories/database"
	"github.com/jonesrussell/goforms/internal/domain/contact"
	"github.com/jonesrussell/goforms/internal/domain/user"
)

// Module provides persistence dependencies
var Module = fx.Module("repositories",
	// Database
	fx.Provide(
		database.NewConfig,
		database.NewDB,
	),

	// Stores
	fx.Provide(
		fx.Annotate(
			contact.NewStore,
			fx.As(new(contact.Store)),
		),
		fx.Annotate(
			user.NewStore,
			fx.As(new(user.Store)),
		),
	),
)

// StoreParams contains dependencies for creating stores
type StoreParams struct {
	fx.In

	DB     *database.DB
	Logger logging.Logger
}

// NewStores creates all database stores
func NewStores(p StoreParams) error {
	p.Logger.Debug("creating database stores",
		logging.Bool("db_available", p.DB != nil),
	)
	return nil
}
