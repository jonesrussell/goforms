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
	mockDB, mock, err := sqlmock.New() // Create a mock database
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// Set expectation for Close
	mock.ExpectClose()

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
	mockDB, mock, err := sqlmock.New() // Create a mock database
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// Set expectation for Begin
	mock.ExpectBegin()

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
	// This test should also use a mock connection instead of a real one.
	// You can implement a similar mock setup for this test if needed.
}
