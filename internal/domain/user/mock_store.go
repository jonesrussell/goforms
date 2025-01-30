package user

import (
	"fmt"

	"github.com/jonesrussell/goforms/internal/domain/common"
)

// MockStore is a mock implementation of the Store interface for testing purposes
type MockStore struct {
	users map[uint]*common.User
}

// NewMockStore creates a new instance of MockStore
func NewMockStore() *MockStore {
	return &MockStore{
		users: make(map[uint]*common.User),
	}
}

// Create simulates storing a new user
func (m *MockStore) Create(u *common.User) error {
	m.users[u.ID] = u
	return nil
}

// Get retrieves a user by ID
func (m *MockStore) Get(id uint) (*common.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetByEmail retrieves a user by email
func (m *MockStore) GetByEmail(email string) (*common.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil // User not found
}

// Update modifies an existing user
func (m *MockStore) Update(u *common.User) error {
	m.users[u.ID] = u
	return nil
}

// Delete removes a user by ID
func (m *MockStore) Delete(id uint) error {
	if _, exists := m.users[id]; !exists {
		return fmt.Errorf("user not found")
	}
	delete(m.users, id)
	return nil
}

// List retrieves all users
func (m *MockStore) List() ([]common.User, error) {
	userList := make([]common.User, 0, len(m.users)) // Pre-allocate with the length of users
	for _, user := range m.users {
		userList = append(userList, *user)
	}
	return userList, nil
}
