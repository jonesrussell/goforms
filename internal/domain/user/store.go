package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories/database"
	"github.com/jonesrussell/goforms/internal/domain/common"
)

// Repository defines the methods for user data access
type Repository interface {
	SignUp(ctx context.Context, signup *Signup) (*common.User, error)
	DeleteUser(ctx context.Context, id uint) error
	GetByEmail(email string) (*common.User, error)
	GetUserByID(ctx context.Context, id uint) (*common.User, error)
	ListUsers(ctx context.Context) ([]common.User, error)
	Login(ctx context.Context, login *Login) (*TokenPair, error)
	Logout(ctx context.Context, token string) error
	UpdateSubmissionStatus(ctx context.Context, id int64, status string) error
	UpdateUser(ctx context.Context, u *common.User) error
	IsTokenBlacklisted(token string) bool
	Create(u *common.User) error
	Get(id uint) (*common.User, error)
	Update(u *common.User) error
	Delete(id uint) error
	List() ([]common.User, error)
}

// StoreImpl implements both Repository and TokenRepository interfaces
type StoreImpl struct {
	db     *database.DB
	logger logging.Logger
}

// NewRepository creates a new user repository
func NewRepository(db *database.DB, logger logging.Logger) Repository {
	return &StoreImpl{
		db:     db,
		logger: logger,
	}
}

// GenerateTokens generates access and refresh tokens for a user
func (s *StoreImpl) GenerateTokens(user *common.User) (*TokenPair, error) {
	// TODO: Implement actual token generation logic (e.g., using JWT)
	accessToken := "generated_access_token"
	refreshToken := "generated_refresh_token"

	s.logger.Debug("Tokens generated successfully", logging.String("email", user.Email))
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// InvalidateToken invalidates a given token
func (s *StoreImpl) InvalidateToken(token string) error {
	// TODO: Implement actual token invalidation logic (e.g., add to blacklist)
	s.logger.Debug("Token invalidated successfully", logging.String("token", token))
	return nil
}

// SignUp implements the Repository interface
func (s *StoreImpl) SignUp(ctx context.Context, signup *Signup) (*common.User, error) {
	// Hash the password before saving
	hashedPassword, err := hashPassword(signup.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", logging.Error(err))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := &common.User{
		Email:          signup.Email,
		HashedPassword: hashedPassword,
		Active:         true, // Default active status
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = s.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to sign up user: %w", err)
	}

	return newUser, nil
}

// DeleteUser implements the Repository interface
func (s *StoreImpl) DeleteUser(ctx context.Context, id uint) error {
	query := `DELETE FROM users WHERE id = ?`

	s.logger.Debug("Deleting user", logging.Uint("id", id))

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		s.logger.Error("Failed to delete user", logging.Error(err))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("Failed to retrieve rows affected", logging.Error(err))
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if rowsAffected == 0 {
		s.logger.Warn("No user found to delete", logging.Uint("id", id))
		return fmt.Errorf("user not found: %d", id)
	}

	s.logger.Info("User deleted successfully", logging.Uint("id", id))
	return nil
}

// GetByEmail implements the Repository interface
func (s *StoreImpl) GetByEmail(email string) (*common.User, error) {
	query := `
		SELECT id, email, hashed_password, first_name, last_name, role, active, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	s.logger.Debug("Fetching user by email", logging.String("email", email))

	var user common.User
	err := s.db.GetContext(context.Background(), &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("User not found", logging.String("email", email))
			return nil, nil // User not found
		}
		s.logger.Error("Failed to fetch user by email", logging.Error(err))
		return nil, fmt.Errorf("failed to fetch user by email: %w", err)
	}

	s.logger.Debug("User fetched successfully", logging.Uint("id", user.ID), logging.String("email", user.Email))
	return &user, nil
}

// GetUserByID implements the Repository interface
func (s *StoreImpl) GetUserByID(ctx context.Context, id uint) (*common.User, error) {
	query := `
		SELECT id, email, hashed_password, first_name, last_name, role, active, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	s.logger.Debug("Fetching user by ID", logging.Uint("id", id))

	var user common.User
	err := s.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("User not found", logging.Uint("id", id))
			return nil, nil // User not found
		}
		s.logger.Error("Failed to fetch user by ID", logging.Error(err))
		return nil, fmt.Errorf("failed to fetch user by ID: %w", err)
	}

	s.logger.Debug("User fetched successfully", logging.Uint("id", user.ID), logging.String("email", user.Email))
	return &user, nil
}

// ListUsers implements the Repository interface
func (s *StoreImpl) ListUsers(ctx context.Context) ([]common.User, error) {
	query := `
		SELECT id, email, hashed_password, first_name, last_name, role, active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	s.logger.Debug("Listing all users")

	var users []common.User
	err := s.db.SelectContext(ctx, &users, query)
	if err != nil {
		s.logger.Error("Failed to list users", logging.Error(err))
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	s.logger.Debug("Users listed successfully", logging.Int("count", len(users)))
	return users, nil
}

// Login implements the Repository interface
func (s *StoreImpl) Login(ctx context.Context, login *Login) (*TokenPair, error) {
	// Implementation of user login (e.g., validate credentials and generate tokens)
	// Placeholder implementation
	return &TokenPair{}, nil
}

// Logout implements the Repository interface
func (s *StoreImpl) Logout(ctx context.Context, token string) error {
	// Implementation of user logout (e.g., invalidate token)
	// Placeholder implementation
	return nil
}

// UpdateSubmissionStatus implements the Repository interface
func (s *StoreImpl) UpdateSubmissionStatus(ctx context.Context, id int64, status string) error {
	// Implementation of updating submission status
	// Placeholder implementation
	return nil
}

// UpdateUser implements the Repository interface
func (s *StoreImpl) UpdateUser(ctx context.Context, u *common.User) error {
	query := `
		UPDATE users
		SET email = ?, hashed_password = ?, first_name = ?, last_name = ?, role = ?, active = ?, updated_at = NOW()
		WHERE id = ?
	`

	s.logger.Debug("Updating user", logging.Uint("id", u.ID))

	result, err := s.db.ExecContext(ctx, query, u.Email, u.HashedPassword, u.FirstName, u.LastName, u.Role, u.Active, u.ID)
	if err != nil {
		s.logger.Error("Failed to update user", logging.Error(err))
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("Failed to retrieve rows affected", logging.Error(err))
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if rowsAffected == 0 {
		s.logger.Warn("No user found to update", logging.Uint("id", u.ID))
		return fmt.Errorf("user not found: %d", u.ID)
	}

	s.logger.Info("User updated successfully", logging.Uint("id", u.ID), logging.String("email", u.Email))
	return nil
}

// IsTokenBlacklisted implements the Repository interface
func (s *StoreImpl) IsTokenBlacklisted(token string) bool {
	// Implementation of token blacklist check
	// Placeholder implementation
	return false
}

// Create implements the Repository interface
func (s *StoreImpl) Create(u *common.User) error {
	s.logger.Debug("Creating user", logging.String("email", u.Email))

	// Hash the password before saving
	hashedPassword, err := hashPassword(u.HashedPassword)
	if err != nil {
		s.logger.Error("Failed to hash password", logging.Error(err))
		return fmt.Errorf("failed to hash password: %w", err)
	}
	u.HashedPassword = hashedPassword

	query := `
		INSERT INTO users (email, hashed_password, first_name, last_name, role, active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err = s.db.ExecContext(context.Background(), query, u.Email, u.HashedPassword, u.FirstName, u.LastName, u.Role, u.Active)
	if err != nil {
		s.logger.Error("Failed to create user", logging.Error(err))
		return fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("User created successfully", logging.String("email", u.Email))
	return nil
}

// Get implements the Repository interface
func (s *StoreImpl) Get(id uint) (*common.User, error) {
	query := `
		SELECT id, email, hashed_password, first_name, last_name, role, active, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	s.logger.Debug("Fetching user by ID", logging.Uint("id", id))

	var user common.User
	err := s.db.GetContext(context.Background(), &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("User not found", logging.Uint("id", id))
			return nil, nil // User not found
		}
		s.logger.Error("Failed to fetch user by ID", logging.Error(err))
		return nil, fmt.Errorf("failed to fetch user by ID: %w", err)
	}

	s.logger.Debug("User fetched successfully", logging.Uint("id", user.ID), logging.String("email", user.Email))
	return &user, nil
}

// Update implements the Repository interface
func (s *StoreImpl) Update(u *common.User) error {
	query := `
		UPDATE users
		SET email = ?, hashed_password = ?, first_name = ?, last_name = ?, role = ?, active = ?, updated_at = NOW()
		WHERE id = ?
	`

	s.logger.Debug("Updating user", logging.Uint("id", u.ID))

	result, err := s.db.ExecContext(context.Background(), query, u.Email, u.HashedPassword, u.FirstName, u.LastName, u.Role, u.Active, u.ID)
	if err != nil {
		s.logger.Error("Failed to update user", logging.Error(err))
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("Failed to retrieve rows affected", logging.Error(err))
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if rowsAffected == 0 {
		s.logger.Warn("No user found to update", logging.Uint("id", u.ID))
		return fmt.Errorf("user not found: %d", u.ID)
	}

	s.logger.Info("User updated successfully", logging.Uint("id", u.ID), logging.String("email", u.Email))
	return nil
}

// Delete implements the Repository interface
func (s *StoreImpl) Delete(id uint) error {
	query := `
		DELETE FROM users
		WHERE id = ?
	`

	s.logger.Debug("Deleting user", logging.Uint("id", id))

	result, err := s.db.ExecContext(context.Background(), query, id)
	if err != nil {
		s.logger.Error("Failed to delete user", logging.Error(err))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("Failed to retrieve rows affected", logging.Error(err))
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if rowsAffected == 0 {
		s.logger.Warn("No user found to delete", logging.Uint("id", id))
		return fmt.Errorf("user not found: %d", id)
	}

	s.logger.Info("User deleted successfully", logging.Uint("id", id))
	return nil
}

// List implements the Repository interface
func (s *StoreImpl) List() ([]common.User, error) {
	query := `
		SELECT id, email, hashed_password, first_name, last_name, role, active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	s.logger.Debug("Listing all users")

	var users []common.User
	err := s.db.SelectContext(context.Background(), &users, query)
	if err != nil {
		s.logger.Error("Failed to list users", logging.Error(err))
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	s.logger.Debug("Users listed successfully", logging.Int("count", len(users)))
	return users, nil
}
