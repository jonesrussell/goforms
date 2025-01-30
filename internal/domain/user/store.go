package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories/database"
	"github.com/jonesrussell/goforms/internal/domain/common"
)

// Store defines the methods for user data access
type Store interface {
	Create(u *common.User) error
	Get(id uint) (*common.User, error)
	GetByEmail(email string) (*common.User, error)
	Update(u *common.User) error
	Delete(id uint) error
	List() ([]common.User, error)
}

// StoreImpl implements the UserRepository interface
type StoreImpl struct {
	db     *database.DB
	logger logging.Logger
}

// Create stores a new user
func (s *StoreImpl) Create(u *common.User) error {
	s.logger.Debug("Creating user", logging.String("email", u.Email))

	// Hash the password before saving
	if err := u.SetPassword(u.Password); err != nil {
		s.logger.Error("Failed to set password", logging.Error(err))
		return err
	}

	_, err := s.db.Exec("INSERT INTO users (email, hashed_password, first_name, last_name, role, active) VALUES (?, ?, ?, ?, ?, ?)",
		u.Email, u.HashedPassword, u.FirstName, u.LastName, u.Role, u.Active)
	if err != nil {
		s.logger.Error("Failed to create user", logging.Error(err))
		return err
	}

	s.logger.Info("User created successfully", logging.String("email", u.Email))
	return nil
}

// Get retrieves a user by ID
func (s *StoreImpl) Get(id uint) (*common.User, error) {
	query := `
		SELECT id, email, hashed_password, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	s.logger.Debug("Getting user by ID", logging.Uint("id", id))

	var u common.User
	if err := s.db.Get(&u, query, id); err != nil {
		s.logger.Error("Failed to get user by ID", logging.Error(err), logging.Uint("id", id))
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	s.logger.Debug("User retrieved", logging.Uint("id", u.ID), logging.String("email", u.Email))
	return &u, nil
}

// GetByEmail retrieves a user by email
func (s *StoreImpl) GetByEmail(email string) (*common.User, error) {
	s.logger.Debug("Getting user by email", logging.String("email", email))

	var u common.User
	err := s.db.Get(&u, "SELECT id, email, hashed_password, created_at, updated_at FROM users WHERE email = ?", email)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("User not found by email", logging.String("email", email))
			return nil, nil // User not found, return nil
		}
		s.logger.Error("Failed to get user by email", logging.Error(err))
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	s.logger.Debug("User retrieved", logging.Uint("id", u.ID), logging.String("email", u.Email))
	return &u, nil
}

// Update modifies an existing user
func (s *StoreImpl) Update(u *common.User) error {
	query := `
		UPDATE users
		SET email = ?, hashed_password = ?, updated_at = NOW()
		WHERE id = ?
	`

	s.logger.Debug("Updating user", logging.Uint("id", u.ID), logging.String("email", u.Email))

	err := s.db.WithTx(context.Background(), func(tx *sqlx.Tx) error {
		result, err := tx.Exec(query, u.Email, u.Password, u.ID)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rows == 0 {
			return fmt.Errorf("user not found: %d", u.ID)
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to update user", logging.Error(err), logging.Uint("id", u.ID), logging.String("email", u.Email))
		return fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Info("User updated successfully", logging.Uint("id", u.ID), logging.String("email", u.Email))
	return nil
}

// Delete removes a user by ID
func (s *StoreImpl) Delete(id uint) error {
	query := `
		DELETE FROM users
		WHERE id = ?
	`

	s.logger.Debug("Deleting user", logging.Uint("id", id))

	err := s.db.WithTx(context.Background(), func(tx *sqlx.Tx) error {
		result, err := tx.Exec(query, id)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rows == 0 {
			return fmt.Errorf("user not found: %d", id)
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to delete user", logging.Error(err), logging.Uint("id", id))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Info("User deleted successfully", logging.Uint("id", id))
	return nil
}

// List retrieves all users
func (s *StoreImpl) List() ([]common.User, error) {
	query := `
		SELECT id, email, hashed_password, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	s.logger.Debug("Listing users")

	var users []common.User
	if err := s.db.Select(&users, query); err != nil {
		s.logger.Error("Failed to list users", logging.Error(err))
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	s.logger.Debug("Users retrieved", logging.Int("count", len(users)))
	return users, nil
}
