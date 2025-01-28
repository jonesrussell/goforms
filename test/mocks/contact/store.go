package mockcontact

import (
	"context"
	"fmt"
	"sync"

	"github.com/jonesrussell/goforms/internal/domain/contact"
)

// Store is a mock implementation of the contact store
type Store struct {
	mu       sync.Mutex
	calls    []mockCall
	expected []mockCall

	CreateFunc       func(ctx context.Context, sub *contact.Submission) error
	ListFunc         func(ctx context.Context) ([]contact.Submission, error)
	GetFunc          func(ctx context.Context, id int64) (*contact.Submission, error)
	UpdateStatusFunc func(ctx context.Context, id int64, status contact.Status) error
}

// NewStore creates a new mock store
func NewStore() *Store {
	return &Store{}
}

// Create mocks the Create method
func (m *Store) Create(ctx context.Context, sub *contact.Submission) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, sub)
	}
	return nil
}

// List mocks the List method
func (m *Store) List(ctx context.Context) ([]contact.Submission, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}

// Get mocks the Get method
func (m *Store) Get(ctx context.Context, id int64) (*contact.Submission, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}

// UpdateStatus mocks the UpdateStatus method
func (m *Store) UpdateStatus(ctx context.Context, id int64, status contact.Status) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, id, status)
	}
	return nil
}

// Verify checks if all expectations were met
func (m *Store) Verify() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.calls) != len(m.expected) {
		return fmt.Errorf("expected %d calls but got %d", len(m.expected), len(m.calls))
	}

	for i, exp := range m.expected {
		got := m.calls[i]
		if exp.method != got.method {
			return fmt.Errorf("call %d: expected method %s but got %s", i, exp.method, got.method)
		}
		if !matchArgs(exp.args, got.args) {
			return fmt.Errorf("call %d: arguments do not match", i)
		}
	}

	return nil
}

// Reset clears all calls and expectations
func (m *Store) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = m.calls[:0]
	m.expected = m.expected[:0]
}
