package user

import (
	"fmt"
	"time"

	"github.com/jonesrussell/goforms/internal/domain/user/models"
)

// Repository defines the methods for user data access
type Repository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List() ([]models.User, error)
}

// Signup represents the user signup request
type Signup struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// Login represents the user login request
type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// NewUserFromSignup creates a new User from a Signup request
func NewUserFromSignup(signup *Signup) (*models.User, error) {
	user := &models.User{
		Email:     signup.Email,
		FirstName: signup.FirstName,
		LastName:  signup.LastName,
		Role:      "user",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.SetPassword(signup.Password); err != nil {
		return nil, fmt.Errorf("failed to set password: %w", err)
	}

	return user, nil
}
