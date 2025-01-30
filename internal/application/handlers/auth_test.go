package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/common"
	"github.com/jonesrussell/goforms/internal/domain/user"
	"github.com/jonesrussell/goforms/internal/test/utils"
)

// MockService implementation
type MockService struct {
	users map[string]*common.User // Simulate a user store
}

// GetUserByID implements user.Service.
func (m *MockService) GetUserByID(ctx context.Context, id uint) (*common.User, error) {
	panic("unimplemented")
}

// IsTokenBlacklisted implements user.Service.
func (m *MockService) IsTokenBlacklisted(token string) bool {
	panic("unimplemented")
}

// ListUsers implements user.Service.
func (m *MockService) ListUsers(ctx context.Context) ([]common.User, error) {
	panic("unimplemented")
}

// Login implements user.Service.
func (m *MockService) Login(ctx context.Context, login *user.Login) (*user.TokenPair, error) {
	panic("unimplemented")
}

// Logout implements user.Service.
func (m *MockService) Logout(ctx context.Context, token string) error {
	panic("unimplemented")
}

// UpdateSubmissionStatus implements user.Service.
func (m *MockService) UpdateSubmissionStatus(ctx context.Context, id int64, status string) error {
	panic("unimplemented")
}

// UpdateUser implements user.Service.
func (m *MockService) UpdateUser(ctx context.Context, u *common.User) error {
	panic("unimplemented")
}

func (m *MockService) GetByEmail(email string) (*common.User, error) {
	if u, exists := m.users[email]; exists {
		return u, nil // User exists
	}
	return nil, nil // User does not exist
}

func (m *MockService) SignUp(signup *user.Signup) (*common.User, error) {
	if _, exists := m.users[signup.Email]; exists {
		return nil, fmt.Errorf("user already exists") // Simulate user already exists
	}
	// Create a new user and add to the mock store
	newUser := &common.User{Email: signup.Email} // Only set the Email field
	m.users[signup.Email] = newUser
	return newUser, nil
}

func (m *MockService) DeleteUser(ctx context.Context, userID uint) error {
	// Assuming you have a way to map userID to email or user
	// For simplicity, let's say we are using email as the identifier
	for email, user := range m.users {
		if user.ID == userID { // Assuming user has an ID field
			delete(m.users, email) // Remove the user from the mock store
			return nil
		}
	}
	return fmt.Errorf("user not found") // Simulate user not found
}

// newRequest creates a new HTTP request with the specified method, URL, and body.
func newRequest(method, url string, body interface{}) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func TestAuthHandler_handleSignup(t *testing.T) {
	mockLogger := &utils.MockLogger{
		DebugFunc: func(msg string, fields ...interface{}) {
			fmt.Printf("DEBUG: %s %v\n", msg, fields)
		},
		ErrorFunc: func(msg string, fields ...logging.Field) {
			fmt.Printf("ERROR: %s %v\n", msg, fields)
		},
		SyncFunc: func() error {
			return nil // Mock Sync behavior
		},
	}

	mockUserService := &MockService{users: make(map[string]*common.User)}
	handler := NewAuthHandler(mockLogger, mockUserService)

	validInput := user.Signup{
		Email:    "uniqueuser@example.com",
		Password: "securepassword",
	}

	// Act
	req, err := newRequest(http.MethodPost, "/api/v1/auth/signup", validInput)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	c := utils.NewContext(req)

	// Call the handler
	if err := handler.handleSignup(c); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Assert
	if c.Response().Status != http.StatusCreated {
		t.Errorf("expected status 201, got %d", c.Response().Status)
	}
}

func TestAuthHandler_handleLogin(t *testing.T) {
	mockLogger := &utils.MockLogger{}
	mockUserService := &user.MockService{}
	handler := NewAuthHandler(mockLogger, mockUserService)

	loginInput := map[string]string{
		"email":    "uniqueuser@example.com",
		"password": "securepassword",
	}

	req, err := newRequest(http.MethodPost, "/api/v1/auth/login", loginInput)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	c := utils.NewContext(req)

	if err := handler.handleLogin(c); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if c.Response().Status != http.StatusOK {
		t.Errorf("expected status 200, got %d", c.Response().Status)
	}
}

func TestAuthHandler_handleLogout(t *testing.T) {
	mockLogger := &utils.MockLogger{}
	mockUserService := &user.MockService{}
	handler := NewAuthHandler(mockLogger, mockUserService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer mocktoken")
	c := utils.NewContext(req)

	if err := handler.handleLogout(c); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if c.Response().Status != http.StatusOK {
		t.Errorf("expected status 200, got %d", c.Response().Status)
	}
}
