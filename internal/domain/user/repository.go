package user

import (
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories/database"
)

// Repository defines the interface for user repository operations.
type Repository interface {
	Create(user *User) error
	Get(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	List() ([]User, error)
	Update(user *User) error
	Delete(id uint) error
}

// userRepository implements the Repository interface.
type userRepository struct {
	logger logging.Logger
	db     *database.DB
}

// NewUserRepository creates a new user repository.
func NewUserRepository(logger logging.Logger, db *database.DB) Repository {
	return &userRepository{
		logger: logger,
		db:     db,
	}
}

// Create stores a new user
func (r *userRepository) Create(user *User) error {
	// Use the logger to log the creation attempt
	r.logger.Debug("Creating user", logging.String("email", user.Email))

	// Implement the logic to create a user in the database
	_, err := r.db.Exec("INSERT INTO users (email, hashed_password) VALUES (?, ?)", user.Email, user.HashedPassword)
	if err != nil {
		r.logger.Error("Failed to create user", logging.Error(err))
		return err
	}

	return nil
}

// Get retrieves a user by ID
func (r *userRepository) Get(id uint) (*User, error) {
	r.logger.Debug("Getting user by ID", logging.Uint("id", id))
	// Implement the logic to get a user by ID from the database
	return nil, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*User, error) {
	r.logger.Debug("Getting user by email", logging.String("email", email))
	// Implement the logic to get a user by email from the database
	return nil, nil
}

// List retrieves all users
func (r *userRepository) List() ([]User, error) {
	r.logger.Debug("Listing users")
	// Implement the logic to list all users from the database
	return nil, nil
}

// Update modifies an existing user
func (r *userRepository) Update(user *User) error {
	r.logger.Debug("Updating user", logging.Uint("id", user.ID))
	// Implement the logic to update a user in the database
	return nil
}

// Delete removes a user by ID
func (r *userRepository) Delete(id uint) error {
	r.logger.Debug("Deleting user", logging.Uint("id", id))
	// Implement the logic to delete a user from the database
	return nil
}
