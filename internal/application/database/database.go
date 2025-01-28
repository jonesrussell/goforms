// Package database provides database connection and management
package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/jonesrussell/goforms/internal/application/logging"
)

// Database wraps the SQL connection pool
type Database struct {
	*sqlx.DB
	logger logging.Logger
}

// Close closes the database connection
func (db *Database) Close() error {
	db.logger.Debug("closing database connection")
	if err := db.DB.Close(); err != nil {
		db.logger.Error("failed to close database connection", logging.Error(err))
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	db.logger.Debug("database connection closed successfully")
	return nil
}

// Begin starts a new transaction with detailed logging
func (db *Database) Begin() (*sqlx.Tx, error) {
	db.logger.Debug("beginning database transaction")
	tx, err := db.DB.Beginx()
	if err != nil {
		db.logger.Error("failed to begin transaction", logging.Error(err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	db.logger.Debug("transaction started successfully")
	return tx, nil
}
