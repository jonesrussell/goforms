package mocksubscription

import (
	"context"
	"sync"

	"github.com/jonesrussell/goforms/internal/domain/subscription"
)

// MockService is a mock implementation of subscription.Service
type MockService struct {
	mu sync.RWMutex

	// Function fields to customize mock behavior
	CreateSubscriptionFunc func(ctx context.Context, sub *subscription.Subscription) error
	GetSubscriptionFunc    func(ctx context.Context, id int64) (*subscription.Subscription, error)
	ListSubscriptionsFunc  func(ctx context.Context) ([]subscription.Subscription, error)
	UpdateStatusFunc       func(ctx context.Context, id int64, status subscription.Status) error
	DeleteSubscriptionFunc func(ctx context.Context, id int64) error

	// Call tracking
	calls struct {
		CreateSubscription []struct{ Sub *subscription.Subscription }
		GetSubscription    []struct{ ID int64 }
		ListSubscriptions  []struct{}
		UpdateStatus       []struct {
			ID     int64
			Status subscription.Status
		}
		DeleteSubscription []struct{ ID int64 }
	}
}

// NewMockService creates a new mock subscription service
func NewMockService() *MockService {
	return &MockService{}
}

// ExpectCreateSubscription sets up expectations for CreateSubscription
func (m *MockService) ExpectCreateSubscription(ctx context.Context, sub *subscription.Subscription, returnErr error) {
	m.CreateSubscriptionFunc = func(ctx context.Context, s *subscription.Subscription) error {
		return returnErr
	}
}

// CreateSubscription implements the Service interface
func (m *MockService) CreateSubscription(ctx context.Context, sub *subscription.Subscription) error {
	m.mu.Lock()
	m.calls.CreateSubscription = append(m.calls.CreateSubscription, struct{ Sub *subscription.Subscription }{Sub: sub})
	m.mu.Unlock()

	return m.CreateSubscriptionFunc(ctx, sub)
}

// ExpectGetSubscription sets up expectations for GetSubscription
func (m *MockService) ExpectGetSubscription(ctx context.Context, id int64, returnSub *subscription.Subscription, returnErr error) {
	m.GetSubscriptionFunc = func(ctx context.Context, i int64) (*subscription.Subscription, error) {
		return returnSub, returnErr
	}
}

// GetSubscription implements the Service interface
func (m *MockService) GetSubscription(ctx context.Context, id int64) (*subscription.Subscription, error) {
	m.mu.Lock()
	m.calls.GetSubscription = append(m.calls.GetSubscription, struct{ ID int64 }{ID: id})
	m.mu.Unlock()

	return m.GetSubscriptionFunc(ctx, id)
}

// ExpectListSubscriptions sets up expectations for ListSubscriptions
func (m *MockService) ExpectListSubscriptions(ctx context.Context, returnSubs []subscription.Subscription, returnErr error) {
	m.ListSubscriptionsFunc = func(ctx context.Context) ([]subscription.Subscription, error) {
		return returnSubs, returnErr
	}
}

// ListSubscriptions implements the Service interface
func (m *MockService) ListSubscriptions(ctx context.Context) ([]subscription.Subscription, error) {
	m.mu.Lock()
	m.calls.ListSubscriptions = append(m.calls.ListSubscriptions, struct{}{})
	m.mu.Unlock()

	return m.ListSubscriptionsFunc(ctx)
}

// ExpectUpdateStatus sets up expectations for UpdateStatus
func (m *MockService) ExpectUpdateStatus(ctx context.Context, id int64, status subscription.Status, returnErr error) {
	m.UpdateStatusFunc = func(ctx context.Context, i int64, s subscription.Status) error {
		return returnErr
	}
}

// UpdateStatus implements the Service interface
func (m *MockService) UpdateStatus(ctx context.Context, id int64, status subscription.Status) error {
	m.mu.Lock()
	m.calls.UpdateStatus = append(m.calls.UpdateStatus, struct {
		ID     int64
		Status subscription.Status
	}{ID: id, Status: status})
	m.mu.Unlock()

	return m.UpdateStatusFunc(ctx, id, status)
}

// ExpectDeleteSubscription sets up expectations for DeleteSubscription
func (m *MockService) ExpectDeleteSubscription(ctx context.Context, id int64, returnErr error) {
	m.DeleteSubscriptionFunc = func(ctx context.Context, i int64) error {
		return returnErr
	}
}

// DeleteSubscription implements the Service interface
func (m *MockService) DeleteSubscription(ctx context.Context, id int64) error {
	m.mu.Lock()
	m.calls.DeleteSubscription = append(m.calls.DeleteSubscription, struct{ ID int64 }{ID: id})
	m.mu.Unlock()

	return m.DeleteSubscriptionFunc(ctx, id)
}

// Verify verifies that all expected calls were made
func (m *MockService) Verify() error {
	return nil
}
