# User Mocks

This package provides mock implementations for the user domain interfaces.

## Mocks

### Service (`service.go`)
Implements `user.Service` interface from `internal/domain/user/service.go`

Methods:
- SignUp(ctx context.Context, signup *user.Signup) (*user.User, error)
- DeleteUser(ctx context.Context, id uint) error
- GetUserByEmail(ctx context.Context, email string) (*user.User, error)
- GetUserByID(ctx context.Context, id uint) (*user.User, error)
- IsTokenBlacklisted(token string) bool
- ListUsers(ctx context.Context) ([]user.User, error)

## Usage

```go
// Example test setup
mockSvc := mockuser.NewMockService()

// Set expectations
mockSvc.ExpectSignUp(ctx, signup, &user.User{}, nil) // expect successful signup
mockSvc.ExpectGetUserByEmail(ctx, "test@example.com", &user.User{}, nil) // expect user found

// Verify expectations were met
if err := mockSvc.Verify(); err != nil {
    t.Error(err)
}
``` 