package mockuser

import (
	"sync"

	"github.com/jonesrussell/goforms/internal/domain/user"
)

// Store is a mock implementation of user.Repository interface
type Store struct {
	mu sync.RWMutex

	// Function fields to customize mock behavior
	CreateFunc     func(u *user.User) error
	GetByIDFunc    func(id uint) (*user.User, error)
	GetByEmailFunc func(email string) (*user.User, error)
	UpdateFunc     func(user *user.User) error
	DeleteFunc     func(id uint) error
	ListFunc       func() ([]user.User, error)

	// Call tracking
	calls struct {
		Create     []struct{ User *user.User }
		GetByID    []struct{ ID uint }
		GetByEmail []struct{ Email string }
		Update     []struct{ User *user.User }
		Delete     []struct{ ID uint }
		List       []struct{}
	}
}

// NewStore creates a new mock store
func NewStore() *Store {
	return &Store{}
}

// Create implements the Store interface
func (m *Store) Create(u *user.User) error {
	m.mu.Lock()
	m.calls.Create = append(m.calls.Create, struct{ User *user.User }{User: u})
	m.mu.Unlock()

	if m.CreateFunc != nil {
		return m.CreateFunc(u)
	}
	return nil
}

// GetByID implements the Store interface
func (m *Store) GetByID(id uint) (*user.User, error) {
	m.mu.Lock()
	m.calls.GetByID = append(m.calls.GetByID, struct{ ID uint }{ID: id})
	m.mu.Unlock()

	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

// GetByEmail implements the Store interface
func (m *Store) GetByEmail(email string) (*user.User, error) {
	m.mu.Lock()
	m.calls.GetByEmail = append(m.calls.GetByEmail, struct{ Email string }{Email: email})
	m.mu.Unlock()

	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(email)
	}
	return nil, nil
}

// Update implements the Store interface
func (m *Store) Update(user *user.User) error {
	m.mu.Lock()
	m.calls.Update = append(m.calls.Update, struct{ User *user.User }{User: user})
	m.mu.Unlock()

	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return nil
}

// Delete implements the Store interface
func (m *Store) Delete(id uint) error {
	m.mu.Lock()
	m.calls.Delete = append(m.calls.Delete, struct{ ID uint }{ID: id})
	m.mu.Unlock()

	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// List implements the Store interface
func (m *Store) List() ([]user.User, error) {
	m.mu.Lock()
	m.calls.List = append(m.calls.List, struct{}{})
	m.mu.Unlock()

	if m.ListFunc != nil {
		return m.ListFunc()
	}
	return nil, nil
}

// CreateCalls returns info about Create calls
func (m *Store) CreateCalls() []struct{ User *user.User } {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.calls.Create
}

// GetByIDCalls returns info about GetByID calls
func (m *Store) GetByIDCalls() []struct{ ID uint } {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.calls.GetByID
}

// GetByEmailCalls returns info about GetByEmail calls
func (m *Store) GetByEmailCalls() []struct{ Email string } {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.calls.GetByEmail
}

// Add SetError method
func (m *Store) SetError(method string, err error) {
	switch method {
	case "create":
		m.CreateFunc = func(u *user.User) error { return err }
	case "getByID":
		m.GetByIDFunc = func(id uint) (*user.User, error) { return nil, err }
	case "getByEmail":
		m.GetByEmailFunc = func(email string) (*user.User, error) { return nil, err }
	case "update":
		m.UpdateFunc = func(user *user.User) error { return err }
	case "delete":
		m.DeleteFunc = func(id uint) error { return err }
	case "list":
		m.ListFunc = func() ([]user.User, error) { return nil, err }
	}
}
