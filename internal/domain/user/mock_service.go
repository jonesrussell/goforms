package user

import (
	"context"
	"errors"
)

// Define domain-specific errors
var (
	ErrInvalidEmail = errors.New("invalid email address")
)

// MockService is a mock implementation of the user.Service interface for testing purposes.
type MockService struct{}

// Ensure MockService implements the user.Service interface
var _ Service = (*MockService)(nil)

// Implement all methods required by the Service interface
func (m *MockService) GetByEmail(email string) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}
	return &User{Email: email}, nil
}

func (m *MockService) SignUp(signupRequest *Signup) (*User, error) {
	return &User{Email: signupRequest.Email}, nil
}

func (m *MockService) Login(ctx context.Context, login *Login) (*TokenPair, error) {
	if login.Email == "" || login.Password == "" {
		return nil, ErrInvalidCredentials
	}
	return &TokenPair{AccessToken: "mockAccessToken", RefreshToken: "mockRefreshToken"}, nil
}

func (m *MockService) UpdateSubmissionStatus(ctx context.Context, id int64, status string) error {
	return nil
}

func (m *MockService) GetUserByID(ctx context.Context, id uint) (*User, error) {
	return &User{ID: id}, nil
}

func (m *MockService) ListUsers(ctx context.Context) ([]User, error) {
	return []User{
		{ID: 1, Email: "user1@example.com"},
		{ID: 2, Email: "user2@example.com"},
	}, nil
}

func (m *MockService) DeleteUser(ctx context.Context, id uint) error {
	return nil
}

func (m *MockService) Signup( /* parameters */ ) (*User, error) {
	return &User{Email: "example@example.com"}, nil
}

func (m *MockService) Logout(ctx context.Context, token string) error {
	return nil
}

func (m *MockService) UpdateUser(ctx context.Context, user *User) error {
	return nil
}

func (m *MockService) IsTokenBlacklisted(token string) bool {
	// Implement mock behavior, e.g., always return false for testing
	return false
}
