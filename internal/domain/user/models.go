package user

import (
	"github.com/jonesrussell/goforms/internal/domain/common"
)

// Request represents a user request.
type Request struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// Signup represents the user signup request
type Signup struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Login represents the user login credentials
type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// TokenPair represents a pair of tokens for authentication
type TokenPair struct {
	AccessToken  string `json:"access_token" db:"access_token"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
}

// ConvertSignupToUser converts a Signup struct to a User struct
func ConvertSignupToUser(signup *Signup) *common.User {
	return &common.User{
		Email:          signup.Email,
		HashedPassword: "",     // Set this later after hashing
		Role:           "user", // Set a default role or modify as needed
		Active:         true,   // Set default active status
	}
}
