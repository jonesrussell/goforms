// Package database provides database connection and management
package database

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/jonesrussell/goforms/internal/application/logging"
)

// Database wraps the SQL connection pool
type Database struct {
	DB     *sqlx.DB
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

// Connect establishes a connection to the database
func Connect(dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return nil, err
	}
	return db, nil
}

type DB struct {
	// Add fields for your database connection, e.g., *sql.DB
}

// NewDB initializes a new DB instance.
func NewDB() *DB {
	// Initialize and return a new DB instance...
	return &DB{}
}
