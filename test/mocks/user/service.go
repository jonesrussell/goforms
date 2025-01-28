package mockuser

import (
	"context"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jonesrussell/goforms/internal/domain/user"
)

// MockService is a mock implementation of user.Service
type MockService struct {
	mu sync.RWMutex

	// Function fields to customize mock behavior
	SignUpFunc             func(ctx context.Context, signup *user.Signup) (*user.User, error)
	LoginFunc              func(ctx context.Context, login *user.Login) (*user.TokenPair, error)
	LogoutFunc             func(ctx context.Context, token string) error
	ValidateTokenFunc      func(tokenString string) (*jwt.Token, error)
	IsTokenBlacklistedFunc func(token string) bool
	DeleteUserFunc         func(ctx context.Context, id uint) error
	GetUserByEmailFunc     func(ctx context.Context, email string) (*user.User, error)
	GetUserByIDFunc        func(ctx context.Context, id uint) (*user.User, error)
	ListUsersFunc          func(ctx context.Context) ([]user.User, error)
	RefreshTokenFunc       func(ctx context.Context, refreshToken string) (*user.TokenPair, error)
	UpdateUserFunc         func(ctx context.Context, user *user.User) error

	// Call tracking
	calls struct {
		SignUp             []struct{ Signup *user.Signup }
		Login              []struct{ Login *user.Login }
		Logout             []struct{ Token string }
		ValidateToken      []struct{ Token string }
		IsTokenBlacklisted []struct{ Token string }
		DeleteUser         []struct{ ID uint }
		GetUserByEmail     []struct{ Email string }
		GetUserByID        []struct{ ID uint }
		ListUsers          []struct{}
		RefreshToken       []struct{ Token string }
		UpdateUser         []struct{ User *user.User }
	}
}

// NewMockService creates a new mock user service
func NewMockService() *MockService {
	return &MockService{}
}

// ExpectSignUp sets up expectations for SignUp
func (m *MockService) ExpectSignUp(ctx context.Context, signup *user.Signup, returnUser *user.User, returnErr error) {
	m.SignUpFunc = func(ctx context.Context, s *user.Signup) (*user.User, error) {
		return returnUser, returnErr
	}
}

// SignUp implements the Service interface
func (m *MockService) SignUp(ctx context.Context, signup *user.Signup) (*user.User, error) {
	m.mu.Lock()
	m.calls.SignUp = append(m.calls.SignUp, struct{ Signup *user.Signup }{Signup: signup})
	m.mu.Unlock()

	return m.SignUpFunc(ctx, signup)
}

// ExpectLogin sets up expectations for Login
func (m *MockService) ExpectLogin(ctx context.Context, login *user.Login, returnTokens *user.TokenPair, returnErr error) {
	m.LoginFunc = func(ctx context.Context, l *user.Login) (*user.TokenPair, error) {
		return returnTokens, returnErr
	}
}

// Login implements the Service interface
func (m *MockService) Login(ctx context.Context, login *user.Login) (*user.TokenPair, error) {
	m.mu.Lock()
	m.calls.Login = append(m.calls.Login, struct{ Login *user.Login }{Login: login})
	m.mu.Unlock()

	return m.LoginFunc(ctx, login)
}

// ExpectLogout sets up expectations for Logout
func (m *MockService) ExpectLogout(ctx context.Context, token string, returnErr error) {
	m.LogoutFunc = func(ctx context.Context, t string) error {
		return returnErr
	}
}

// Logout implements the Service interface
func (m *MockService) Logout(ctx context.Context, token string) error {
	m.mu.Lock()
	m.calls.Logout = append(m.calls.Logout, struct{ Token string }{Token: token})
	m.mu.Unlock()

	return m.LogoutFunc(ctx, token)
}

// ExpectValidateToken sets up expectations for ValidateToken
func (m *MockService) ExpectValidateToken(token string, returnToken *jwt.Token, returnErr error) {
	m.ValidateTokenFunc = func(t string) (*jwt.Token, error) {
		return returnToken, returnErr
	}
}

// ValidateToken implements the Service interface
func (m *MockService) ValidateToken(token string) (*jwt.Token, error) {
	m.mu.Lock()
	m.calls.ValidateToken = append(m.calls.ValidateToken, struct{ Token string }{Token: token})
	m.mu.Unlock()

	return m.ValidateTokenFunc(token)
}

// ExpectIsTokenBlacklisted sets up expectations for IsTokenBlacklisted
func (m *MockService) ExpectIsTokenBlacklisted(token string, returnBool bool) {
	m.IsTokenBlacklistedFunc = func(t string) bool {
		return returnBool
	}
}

// IsTokenBlacklisted implements the Service interface
func (m *MockService) IsTokenBlacklisted(token string) bool {
	m.mu.Lock()
	m.calls.IsTokenBlacklisted = append(m.calls.IsTokenBlacklisted, struct{ Token string }{Token: token})
	m.mu.Unlock()

	return m.IsTokenBlacklistedFunc(token)
}

// ExpectDeleteUser sets up expectations for DeleteUser
func (m *MockService) ExpectDeleteUser(ctx context.Context, id uint, returnErr error) {
	m.DeleteUserFunc = func(ctx context.Context, i uint) error {
		return returnErr
	}
}

// DeleteUser implements the Service interface
func (m *MockService) DeleteUser(ctx context.Context, id uint) error {
	m.mu.Lock()
	m.calls.DeleteUser = append(m.calls.DeleteUser, struct{ ID uint }{ID: id})
	m.mu.Unlock()

	return m.DeleteUserFunc(ctx, id)
}

// ExpectGetUserByEmail sets up expectations for GetUserByEmail
func (m *MockService) ExpectGetUserByEmail(ctx context.Context, email string, returnUser *user.User, returnErr error) {
	m.GetUserByEmailFunc = func(ctx context.Context, e string) (*user.User, error) {
		return returnUser, returnErr
	}
}

// GetUserByEmail implements the Service interface
func (m *MockService) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	m.mu.Lock()
	m.calls.GetUserByEmail = append(m.calls.GetUserByEmail, struct{ Email string }{Email: email})
	m.mu.Unlock()

	return m.GetUserByEmailFunc(ctx, email)
}

// ExpectGetUserByID sets up expectations for GetUserByID
func (m *MockService) ExpectGetUserByID(ctx context.Context, id uint, returnUser *user.User, returnErr error) {
	m.GetUserByIDFunc = func(ctx context.Context, i uint) (*user.User, error) {
		return returnUser, returnErr
	}
}

// GetUserByID implements the Service interface
func (m *MockService) GetUserByID(ctx context.Context, id uint) (*user.User, error) {
	m.mu.Lock()
	m.calls.GetUserByID = append(m.calls.GetUserByID, struct{ ID uint }{ID: id})
	m.mu.Unlock()

	return m.GetUserByIDFunc(ctx, id)
}

// ExpectListUsers sets up expectations for ListUsers
func (m *MockService) ExpectListUsers(ctx context.Context, returnUsers []user.User, returnErr error) {
	m.ListUsersFunc = func(ctx context.Context) ([]user.User, error) {
		return returnUsers, returnErr
	}
}

// ListUsers implements the Service interface
func (m *MockService) ListUsers(ctx context.Context) ([]user.User, error) {
	m.mu.Lock()
	m.calls.ListUsers = append(m.calls.ListUsers, struct{}{})
	m.mu.Unlock()

	return m.ListUsersFunc(ctx)
}

// ExpectRefreshToken sets up expectations for RefreshToken
func (m *MockService) ExpectRefreshToken(ctx context.Context, refreshToken string, returnTokens *user.TokenPair, returnErr error) {
	m.RefreshTokenFunc = func(ctx context.Context, t string) (*user.TokenPair, error) {
		return returnTokens, returnErr
	}
}

// RefreshToken implements the Service interface
func (m *MockService) RefreshToken(ctx context.Context, refreshToken string) (*user.TokenPair, error) {
	m.mu.Lock()
	m.calls.RefreshToken = append(m.calls.RefreshToken, struct{ Token string }{Token: refreshToken})
	m.mu.Unlock()

	return m.RefreshTokenFunc(ctx, refreshToken)
}

// ExpectUpdateUser sets up expectations for UpdateUser
func (m *MockService) ExpectUpdateUser(ctx context.Context, user *user.User, returnErr error) {
	m.UpdateUserFunc = func(ctx context.Context, u *user.User) error {
		return returnErr
	}
}

// UpdateUser implements the Service interface
func (m *MockService) UpdateUser(ctx context.Context, user *user.User) error {
	m.mu.Lock()
	m.calls.UpdateUser = append(m.calls.UpdateUser, struct{ User *user.User }{User: user})
	m.mu.Unlock()

	return m.UpdateUserFunc(ctx, user)
}

// Verify verifies that all expected calls were made
func (m *MockService) Verify() error {
	return nil
}
