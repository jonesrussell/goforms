package usermock

import (
	"errors"
	"sync"

	"github.com/jonesrussell/goforms/internal/domain/user"
	"github.com/stretchr/testify/mock"
)

var (
	ErrNotFound = errors.New("user not found")
	ErrNoValue  = errors.New("no value returned")
)

// UserStore is a mock implementation of user.Store interface
type UserStore struct {
	mock.Mock
	users     map[uint]*user.User
	emailMap  map[string]uint
	mu        sync.RWMutex
	nextID    uint
	createErr error
	getErr    error
	updateErr error
	deleteErr error
	listErr   error
}

var _ user.Store = (*UserStore)(nil) // Ensure UserStore implements user.Store

// NewUserStore creates a new mock user store
func NewUserStore() *UserStore {
	return &UserStore{
		users:    make(map[uint]*user.User),
		emailMap: make(map[string]uint),
		nextID:   1,
	}
}

// SetError sets the error to be returned by the specified operation
func (m *UserStore) SetError(op string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch op {
	case "create":
		m.createErr = err
	case "get":
		m.getErr = err
	case "update":
		m.updateErr = err
	case "delete":
		m.deleteErr = err
	case "list":
		m.listErr = err
	}
}

// Create implements Store.Create
func (m *UserStore) Create(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

// GetByID implements Store.GetByID
func (m *UserStore) GetByID(id uint) (*user.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, user.ErrUserNotFound
	}
	u, ok := args.Get(0).(*user.User)
	if !ok {
		return nil, errors.New("invalid type assertion for user")
	}
	return u, args.Error(1)
}

// GetByEmail implements Store.GetByEmail
func (m *UserStore) GetByEmail(email string) (*user.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, ErrNoValue
	}
	u, ok := args.Get(0).(*user.User)
	if !ok {
		return nil, errors.New("invalid type assertion for user")
	}
	return u, args.Error(1)
}

// Update implements Store.Update
func (m *UserStore) Update(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

// Delete implements Store.Delete
func (m *UserStore) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// List implements Store.List
func (m *UserStore) List() ([]user.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	users, ok := args.Get(0).([]user.User)
	if !ok {
		return nil, errors.New("invalid type assertion for users")
	}
	return users, args.Error(1)
}
