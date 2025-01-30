package repositories

import (
	"errors"

	"github.com/jonesrussell/goforms/internal/domain/user"
)

// MockDB is a mock implementation of the database interface for testing purposes.
type MockDB struct {
	data map[string]interface{}
}

// NewMockDB creates a new instance of MockDB.
func NewMockDB() *MockDB {
	return &MockDB{
		data: make(map[string]interface{}),
	}
}

// Get simulates retrieving a record from the mock database.
func (db *MockDB) Get(dest interface{}, query string, args ...interface{}) error {
	// Simulate a simple retrieval logic
	if user, ok := db.data[args[0].(string)]; ok {
		dest = user
		return nil
	}
	return errors.New("not found")
}

// Create simulates inserting a record into the mock database.
func (db *MockDB) Create(user *user.User) error {
	// Simulate a simple insert logic
	db.data[user.Email] = user
	return nil
}

// Close simulates closing the database connection.
func (db *MockDB) Close() error {
	return nil
}
