package contact

import (
	"context"
	"errors"

	"github.com/jonesrussell/goforms/internal/domain/contact"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

// List returns a list of submissions based on expectations
func (m *MockStore) List(ctx context.Context) ([]*contact.Submission, error) {
	args := m.Called(ctx)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	submissions, ok := result.([]*contact.Submission)
	if !ok {
		return nil, errors.New("unexpected type for submissions")
	}
	return submissions, args.Error(1)
}

// Get returns a single submission based on expectations
func (m *MockStore) Get(ctx context.Context, id string) (*contact.Submission, error) {
	args := m.Called(ctx, id)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	submission, ok := result.(*contact.Submission)
	if !ok {
		return nil, errors.New("unexpected type for submission")
	}
	return submission, args.Error(1)
}

// Create creates a new submission based on expectations
func (m *MockStore) Create(ctx context.Context, submission *contact.Submission) error {
	args := m.Called(ctx, submission)
	return args.Error(0)
}

// UpdateStatus updates the status of a contact submission
func (m *MockStore) UpdateStatus(ctx context.Context, id, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

// Expect sets up an expectation for a method call
func (m *MockStore) Expect(method string, args ...any) *MockStore {
	m.Mock.On(method, args...)
	return m
}

// Return sets up return values for a method call
func (m *MockStore) Return(method string, returns ...any) *MockStore {
	m.Mock.On(method, returns...)
	return m
}

func (m *MockStore) ExpectationsWereMet(t mock.TestingT) error {
	m.Mock.AssertExpectations(t)
	return nil
}

func (m *MockStore) On(method string, args ...any) *mock.Call {
	return m.Mock.On(method, args...)
}

func (m *MockStore) AssertExpectations(t mock.TestingT) bool {
	return m.Mock.AssertExpectations(t)
}
