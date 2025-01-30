package user

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories/database"
	"github.com/jonesrussell/goforms/internal/domain/common"
)

// Repository defines the interface for user repository operations.
type Repository interface {
	Create(u *common.User) error
	Get(id uint) (*common.User, error)
	GetByEmail(email string) (*common.User, error)
	List() ([]common.User, error)
	Update(u *common.User) error
	Delete(id uint) error
	SignUp(user *Signup) (*common.User, error)
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
func (r *userRepository) Create(u *common.User) error {
	// Log the user details being saved (excluding the password)
	r.logger.Debug("Saving user to database", logging.Any("user", map[string]interface{}{
		"email":      u.Email,
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"role":       u.Role,
		"active":     u.Active,
	}))

	// Use the logger to log the creation attempt
	r.logger.Debug("Creating user", logging.String("email", u.Email))

	// Implement the logic to create a user in the database
	_, err := r.db.Exec("INSERT INTO users (email, password, first_name, last_name, role, active) VALUES (?, ?, ?, ?, ?, ?)",
		u.Email, u.Password, u.FirstName, u.LastName, u.Role, u.Active)
	if err != nil {
		r.logger.Error("Failed to create user", err)
		return err
	}

	return nil
}

// Get retrieves a user by ID
func (r *userRepository) Get(id uint) (*common.User, error) {
	r.logger.Debug("Getting user by ID", logging.Uint("id", id))
	u := &common.User{}
	err := r.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Role, &u.Active)
	if err != nil {
		r.logger.Error("Failed to get user by ID", err)
		return nil, err
	}
	return u, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*common.User, error) {
	var u common.User
	err := r.db.Get(&u, "SELECT id, email, hashed_password, created_at, updated_at FROM users WHERE email = ?", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &u, nil
}

// List retrieves all users
func (r *userRepository) List() ([]common.User, error) {
	r.logger.Debug("Listing users")
	rows, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		r.logger.Error("Failed to list users", err)
		return nil, err
	}
	defer rows.Close()

	var users []common.User
	for rows.Next() {
		u := common.User{}
		if err := rows.Scan(
			&u.Email,
			&u.Password,
			&u.FirstName,
			&u.LastName,
			&u.Role,
			&u.Active,
		); err != nil {
			r.logger.Error("Failed to scan user", err)
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// Update modifies an existing user
func (r *userRepository) Update(u *common.User) error {
	r.logger.Debug("Updating user", logging.Uint("id", u.ID))
	_, err := r.db.Exec("UPDATE users SET email = ?, password = ?, first_name = ?, last_name = ?, role = ?, active = ? WHERE id = ?",
		u.Email, u.Password, u.FirstName, u.LastName, u.Role, u.Active, u.ID)
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

// SignUp implements the SignUp method of the Repository interface.
func (r *userRepository) SignUp(user *Signup) (*common.User, error) {
	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error("Failed to hash password", logging.Error(err))
		return nil, err
	}

	// Create a new common.User instance
	newUser := &common.User{
		Email:          user.Email,
		HashedPassword: string(hashedPassword),
		Active:         true, // Set default active status
	}

	// Implement the logic to create a user in the database
	_, err = r.db.Exec("INSERT INTO users (email, password, first_name, last_name, role, active) VALUES (?, ?, ?, ?, ?, ?)",
		newUser.Email, newUser.Password, newUser.FirstName, newUser.LastName, newUser.Role, newUser.Active)
	if err != nil {
		r.logger.Error("Failed to create user", err)
		return nil, err
	}

	return newUser, nil
}
