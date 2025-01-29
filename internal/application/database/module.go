package database

import (
	"go.uber.org/fx"

	"github.com/jonesrussell/goforms/internal/application/config"
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories/database"
)

// Module provides the database connection
var Module = fx.Module("database",
	fx.Provide(
		func(lc fx.Lifecycle, cfg *config.Config, logger logging.Logger) (*database.DB, error) {
			return database.NewDB(lc, cfg, logger) // Ensure this is the only connection method used
		},
	),
)
