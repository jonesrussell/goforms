package repositories

import (
	"fmt"

	"github.com/jonesrussell/goforms/internal/domain/common"
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
	// Check if user exists in the mock database
	if u, ok := db.data[args[0].(string)]; ok {
		// Use a pointer to modify the value of dest
		switch d := dest.(type) {
		case *common.User:
			*d = *(u.(*common.User)) // Dereference u to assign the value
		default:
			return fmt.Errorf("unsupported destination type")
		}
		return nil
	}
	return fmt.Errorf("user not found")
}

// Create simulates inserting a record into the mock database.
func (db *MockDB) Create(user *common.User) error {
	// Simulate a simple insert logic
	db.data[user.Email] = user
	return nil
}

// Close simulates closing the database connection.
func (db *MockDB) Close() error {
	return nil
}
