package mockuser

import (
	"context"
	"fmt"
	"sync"

	"github.com/jonesrussell/goforms/internal/domain/user"
)

// MockService is a mock implementation of user.Service
type MockService struct {
	mu       sync.Mutex
	calls    []mockCall
	expected []mockCall
}

type mockCall struct {
	method string
	args   []interface{}
	ret    []interface{}
}

// NewMockService creates a new mock service
func NewMockService() *MockService {
	return &MockService{}
}

// ExpectSignUp sets up an expectation for SignUp method
func (m *MockService) ExpectSignUp(ctx context.Context, signup *user.Signup, ret *user.User, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.expected = append(m.expected, mockCall{
		method: "SignUp",
		args:   []interface{}{ctx, signup},
		ret:    []interface{}{ret, err},
	})
}

// SignUp implements the user.Service interface
func (m *MockService) SignUp(ctx context.Context, signup *user.Signup) (*user.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := mockCall{
		method: "SignUp",
		args:   []interface{}{ctx, signup},
	}
	m.calls = append(m.calls, call)

	// Find matching expectation
	for _, exp := range m.expected {
		if exp.method == "SignUp" {
			if ret := exp.ret[0]; ret != nil {
				return ret.(*user.User), exp.ret[1].(error)
			}
			return nil, exp.ret[1].(error)
		}
	}
	return nil, nil
}

// Verify checks if all expectations were met
func (m *MockService) Verify() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.calls) != len(m.expected) {
		return fmt.Errorf("expected %d calls but got %d", len(m.expected), len(m.calls))
	}

	return nil
}

// Reset clears all calls and expectations
func (m *MockService) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = m.calls[:0]
	m.expected = m.expected[:0]
}

// Add DeleteUser method to implement the full interface
func (m *MockService) DeleteUser(ctx context.Context, id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := mockCall{
		method: "DeleteUser",
		args:   []interface{}{ctx, id},
	}
	m.calls = append(m.calls, call)

	// Find matching expectation
	for _, exp := range m.expected {
		if exp.method == "DeleteUser" {
			return exp.ret[0].(error)
		}
	}
	return nil
}

// Add corresponding expectation method
func (m *MockService) ExpectDeleteUser(ctx context.Context, id uint, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.expected = append(m.expected, mockCall{
		method: "DeleteUser",
		args:   []interface{}{ctx, id},
		ret:    []interface{}{err},
	})
}

// GetUserByEmail implements the user.Service interface
func (m *MockService) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := mockCall{
		method: "GetUserByEmail",
		args:   []interface{}{ctx, email},
	}
	m.calls = append(m.calls, call)

	// Find matching expectation
	for _, exp := range m.expected {
		if exp.method == "GetUserByEmail" {
			if ret := exp.ret[0]; ret != nil {
				return ret.(*user.User), exp.ret[1].(error)
			}
			return nil, exp.ret[1].(error)
		}
	}
	return nil, nil
}

// ExpectGetUserByEmail sets up an expectation for GetUserByEmail method
func (m *MockService) ExpectGetUserByEmail(ctx context.Context, email string, ret *user.User, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.expected = append(m.expected, mockCall{
		method: "GetUserByEmail",
		args:   []interface{}{ctx, email},
		ret:    []interface{}{ret, err},
	})
}

// GetUserByID implements the user.Service interface
func (m *MockService) GetUserByID(ctx context.Context, id uint) (*user.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := mockCall{
		method: "GetUserByID",
		args:   []interface{}{ctx, id},
	}
	m.calls = append(m.calls, call)

	// Find matching expectation
	for _, exp := range m.expected {
		if exp.method == "GetUserByID" {
			if ret := exp.ret[0]; ret != nil {
				return ret.(*user.User), exp.ret[1].(error)
			}
			return nil, exp.ret[1].(error)
		}
	}
	return nil, nil
}

// ExpectGetUserByID sets up an expectation for GetUserByID method
func (m *MockService) ExpectGetUserByID(ctx context.Context, id uint, ret *user.User, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.expected = append(m.expected, mockCall{
		method: "GetUserByID",
		args:   []interface{}{ctx, id},
		ret:    []interface{}{ret, err},
	})
}

// IsTokenBlacklisted implements the user.Service interface
func (m *MockService) IsTokenBlacklisted(token string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := mockCall{
		method: "IsTokenBlacklisted",
		args:   []interface{}{token},
	}
	m.calls = append(m.calls, call)

	// Find matching expectation
	for _, exp := range m.expected {
		if exp.method == "IsTokenBlacklisted" {
			return exp.ret[0].(bool)
		}
	}
	return false
}

// ExpectIsTokenBlacklisted sets up an expectation for IsTokenBlacklisted method
func (m *MockService) ExpectIsTokenBlacklisted(token string, ret bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.expected = append(m.expected, mockCall{
		method: "IsTokenBlacklisted",
		args:   []interface{}{token},
		ret:    []interface{}{ret},
	})
}
