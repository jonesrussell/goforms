package contact

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories/database"
	"github.com/jonesrussell/goforms/internal/domain/common"
)

// Repository implements the Store interface
type Repository struct {
	db     *database.DB
	logger logging.Logger
}

// NewContactStore creates a new contact store
func NewContactStore(db *database.DB, logger logging.Logger) Store {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new contact submission
func (s *Repository) Create(ctx context.Context, sub *common.Submission) error {
	query := `
		INSERT INTO contact_submissions (name, email, message, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	result, err := s.db.ExecContext(ctx, query,
		sub.Name,
		sub.Email,
		sub.Message,
		sub.Status,
	)
	if err != nil {
		s.logger.Error("failed to create contact submission", logging.Error(err))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	sub.ID = id
	return nil
}

// List returns all contact submissions
func (s *Repository) List(ctx context.Context) ([]common.Submission, error) {
	var subs []common.Submission
	query := `
		SELECT id, name, email, message, status, created_at, updated_at
		FROM contact_submissions
		ORDER BY created_at DESC
	`

	if err := s.db.SelectContext(ctx, &subs, query); err != nil {
		s.logger.Error("failed to list contact submissions", logging.Error(err))
		return nil, err
	}

	return subs, nil
}

// Get returns a contact submission by ID
func (s *Repository) Get(ctx context.Context, id int64) (*common.Submission, error) {
	var sub common.Submission
	query := `
		SELECT id, name, email, message, status, created_at, updated_at
		FROM contact_submissions
		WHERE id = ?
	`

	err := s.db.GetContext(ctx, &sub, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		s.logger.Error("failed to get contact submission", logging.Error(err))
		return nil, err
	}

	return &sub, nil
}

// UpdateStatus updates the status of a contact submission
func (s *Repository) UpdateStatus(ctx context.Context, id int64, status common.Status) error {
	query := `
		UPDATE contact_submissions
		SET status = ?, updated_at = NOW()
		WHERE id = ?
	`

	result, err := s.db.ExecContext(ctx, query, string(status), id)
	if err != nil {
		s.logger.Error("failed to update contact submission status", logging.Error(err))
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("contact submission not found")
	}

	return nil
}
