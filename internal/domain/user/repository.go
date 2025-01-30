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
	// Log the user details being saved (excluding the password)
	r.logger.Debug("Saving user to database", logging.Any("user", map[string]interface{}{
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.Role,
		"active":     user.Active,
	}))

	// Use the logger to log the creation attempt
	r.logger.Debug("Creating user", logging.String("email", user.Email))

	// Implement the logic to create a user in the database
	_, err := r.db.Exec("INSERT INTO users (email, hashed_password, first_name, last_name, role, active) VALUES (?, ?, ?, ?, ?, ?)",
		user.Email, user.HashedPassword, user.FirstName, user.LastName, user.Role, user.Active)
	if err != nil {
		r.logger.Error("Failed to create user", err)
		return err
	}

	return nil
}

// Get retrieves a user by ID
func (r *userRepository) Get(id uint) (*User, error) {
	r.logger.Debug("Getting user by ID", logging.Uint("id", id))
	user := &User{}
	err := r.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user.Email, &user.HashedPassword, &user.FirstName, &user.LastName, &user.Role, &user.Active)
	if err != nil {
		r.logger.Error("Failed to get user by ID", err)
		return nil, err
	}
	return user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*User, error) {
	r.logger.Debug("Getting user by email", logging.String("email", email))
	user := &User{}
	err := r.db.QueryRow("SELECT * FROM users WHERE email = ?", email).Scan(&user.Email, &user.HashedPassword, &user.FirstName, &user.LastName, &user.Role, &user.Active)
	if err != nil {
		r.logger.Error("Failed to get user by email", err)
		return nil, err
	}
	r.logger.Debug("User retrieved", logging.Any("user", user))
	return user, nil
}

// List retrieves all users
func (r *userRepository) List() ([]User, error) {
	r.logger.Debug("Listing users")
	rows, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		r.logger.Error("Failed to list users", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.Email, &user.HashedPassword, &user.FirstName, &user.LastName, &user.Role, &user.Active); err != nil {
			r.logger.Error("Failed to scan user", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Update modifies an existing user
func (r *userRepository) Update(user *User) error {
	r.logger.Debug("Updating user", logging.Uint("id", user.ID))
	_, err := r.db.Exec("UPDATE users SET email = ?, hashed_password = ?, first_name = ?, last_name = ?, role = ?, active = ? WHERE id = ?",
		user.Email, user.HashedPassword, user.FirstName, user.LastName, user.Role, user.Active, user.ID)
	if err != nil {
		r.logger.Error("Failed to update user", err)
		return err
	}
	return nil
}

// Delete removes a user by ID
func (r *userRepository) Delete(id uint) error {
	r.logger.Debug("Deleting user", logging.Uint("id", id))
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		r.logger.Error("Failed to delete user", err)
		return err
	}
	return nil
}
