package mockcontact

import (
	"context"

	"github.com/jonesrussell/goforms/internal/domain/contact"
)

// MockStore is a mock implementation of the contact.Store interface
type MockStore struct {
	CreateFunc       func(ctx context.Context, sub *contact.Submission) error
	GetFunc          func(ctx context.Context, id int64) (*contact.Submission, error)
	ListFunc         func(ctx context.Context) ([]contact.Submission, error)
	UpdateStatusFunc func(ctx context.Context, id int64, status contact.Status) error
}

// NewMockStore creates a new instance of MockStore
func NewMockStore() *MockStore {
	return &MockStore{}
}

// Create calls the mocked CreateFunc
func (m *MockStore) Create(ctx context.Context, sub *contact.Submission) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, sub)
	}
	return nil
}

// Get calls the mocked GetFunc
func (m *MockStore) Get(ctx context.Context, id int64) (*contact.Submission, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}

// List calls the mocked ListFunc
func (m *MockStore) List(ctx context.Context) ([]contact.Submission, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}

// UpdateStatus calls the mocked UpdateStatusFunc
func (m *MockStore) UpdateStatus(ctx context.Context, id int64, status contact.Status) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, id, status)
	}
	return nil
}
