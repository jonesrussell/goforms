package repositories

import (
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/user/models"
)

// UserStore represents the user store
type UserStore struct {
	users map[uint]*models.User
}

// NewUserStore creates a new user store
func NewUserStore(logger logging.Logger) *UserStore {
	logger.Debug("creating user store")
	return &UserStore{
		users: make(map[uint]*models.User),
	}
}

// GetByEmail retrieves a user by their email address
func (s *UserStore) GetByEmail(email string) (*models.User, error) {
	// Example placeholder logic:
	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil // Or return an error if not found
}

// Implement other methods to match the user.Repository interface...
