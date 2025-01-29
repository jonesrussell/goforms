// database_test.go
package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/test/utils"
)

func TestDatabase_Close_Success(t *testing.T) {
	// Arrange
	mockDB, _, err := sqlmock.New() // Create a mock database
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	mockLogger := &utils.MockLogger{
		DebugFunc: func(msg string, fields ...interface{}) {},
		ErrorFunc: func(msg string, fields ...logging.Field) {},
		InfoFunc:  func(msg string, fields ...logging.Field) {},
		WarnFunc:  func(msg string, fields ...logging.Field) {},
	}
	db := &Database{
		DB:     sqlx.NewDb(mockDB, "mysql"), // Use the mock DB
		logger: mockLogger,
	}

	// Act
	err = db.Close()

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDatabase_Begin_Success(t *testing.T) {
	// Arrange
	mockDB, _, err := sqlmock.New() // Create a mock database
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	mockLogger := &utils.MockLogger{
		DebugFunc: func(msg string, fields ...interface{}) {
			// Log debug messages with structured fields
		},
		ErrorFunc: func(msg string, fields ...logging.Field) {
			// Log error messages with structured fields
		},
		InfoFunc: func(msg string, fields ...logging.Field) {
			// Log info messages with structured fields
		},
		WarnFunc: func(msg string, fields ...logging.Field) {
			// Log warning messages with structured fields
		},
	}

	db := &Database{
		DB:     sqlx.NewDb(mockDB, "mysql"), // Use the mock DB
		logger: mockLogger,
	}

	// Act
	tx, err := db.Begin()

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if tx == nil {
		t.Error("expected a transaction, got nil")
	}
}

func TestConnect_Success(t *testing.T) {
	// Arrange
	dataSourceName := "user:password@tcp(localhost:3306)/dbname" // Use a valid DSN for testing

	// Act
	db, err := Connect(dataSourceName)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if db == nil {
		t.Error("expected a database connection, got nil")
	}
}
