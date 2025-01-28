package mockcontact

import (
	"context"
	"fmt"
	"sync"

	"github.com/jonesrussell/goforms/internal/domain/contact"
)

// Service is a mock implementation of the contact service
type Service struct {
	mu       sync.Mutex
	calls    []mockCall
	expected []mockCall
}

// NewService creates a new mock service
func NewService() *Service {
	return &Service{}
}

// recordCall records a method call
func (m *Service) recordCall(method string, args []interface{}) []interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := mockCall{method: method, args: args}
	m.calls = append(m.calls, call)

	// Find matching expectation
	for _, exp := range m.expected {
		if exp.method == method && matchArgs(exp.args, args) {
			return exp.ret
		}
	}
	return nil
}

// ExpectSubmit sets up an expectation for Submit method
func (m *Service) ExpectSubmit(ctx context.Context, sub *contact.Submission, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.expected = append(m.expected, mockCall{
		method: "Submit",
		args:   []interface{}{ctx, sub},
		ret:    []interface{}{err},
	})
}

// ExpectListSubmissions sets up an expectation for ListSubmissions method
func (m *Service) ExpectListSubmissions(ctx context.Context, ret []contact.Submission, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.expected = append(m.expected, mockCall{
		method: "ListSubmissions",
		args:   []interface{}{ctx},
		ret:    []interface{}{ret, err},
	})
}

// ExpectGetSubmission sets up an expectation for GetSubmission method
func (m *Service) ExpectGetSubmission(ctx context.Context, id int64, ret *contact.Submission, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.expected = append(m.expected, mockCall{
		method: "GetSubmission",
		args:   []interface{}{ctx, id},
		ret:    []interface{}{ret, err},
	})
}

// ExpectUpdateSubmissionStatus sets up an expectation for UpdateSubmissionStatus method
func (m *Service) ExpectUpdateSubmissionStatus(ctx context.Context, id int64, status contact.Status, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.expected = append(m.expected, mockCall{
		method: "UpdateSubmissionStatus",
		args:   []interface{}{ctx, id, status},
		ret:    []interface{}{err},
	})
}

// Submit mocks the Submit method
func (m *Service) Submit(ctx context.Context, sub *contact.Submission) error {
	ret := m.recordCall("Submit", []interface{}{ctx, sub})
	if ret == nil || ret[0] == nil {
		return nil
	}
	return ret[0].(error)
}

// ListSubmissions mocks the ListSubmissions method
func (m *Service) ListSubmissions(ctx context.Context) ([]contact.Submission, error) {
	ret := m.recordCall("ListSubmissions", []interface{}{ctx})
	if ret == nil {
		return nil, nil
	}
	var subs []contact.Submission
	if ret[0] != nil {
		subs = ret[0].([]contact.Submission)
	}
	var err error
	if ret[1] != nil {
		err = ret[1].(error)
	}
	return subs, err
}

// GetSubmission mocks the GetSubmission method
func (m *Service) GetSubmission(ctx context.Context, id int64) (*contact.Submission, error) {
	ret := m.recordCall("GetSubmission", []interface{}{ctx, id})
	if ret == nil {
		return nil, nil
	}
	var sub *contact.Submission
	if ret[0] != nil {
		sub = ret[0].(*contact.Submission)
	}
	var err error
	if ret[1] != nil {
		err = ret[1].(error)
	}
	return sub, err
}

// UpdateSubmissionStatus mocks the UpdateSubmissionStatus method
func (m *Service) UpdateSubmissionStatus(ctx context.Context, id int64, status contact.Status) error {
	ret := m.recordCall("UpdateSubmissionStatus", []interface{}{ctx, id, status})
	if ret == nil || ret[0] == nil {
		return nil
	}
	return ret[0].(error)
}

// Verify checks if all expectations were met
func (m *Service) Verify() error {
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
func (m *Service) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = m.calls[:0]
	m.expected = m.expected[:0]
}
