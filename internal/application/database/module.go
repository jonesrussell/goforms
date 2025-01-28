package database

import (
	"github.com/jonesrussell/goforms/internal/application/config"
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories/database"
	"go.uber.org/fx"
)

// Module provides the database connection
var Module = fx.Module("database",
	fx.Provide(
		func(lc fx.Lifecycle, cfg *config.Config, logger logging.Logger) (*database.DB, error) {
			return database.NewDB(lc, cfg, logger) // Use the NewDB from connection.go
		},
	),
)
