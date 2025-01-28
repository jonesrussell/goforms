package contact

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/jonesrussell/goforms/internal/application/logging"
)

// ContactStore defines the methods for contact data access
type Store interface {
	Create(ctx context.Context, sub *Submission) error
	List(ctx context.Context) ([]Submission, error)
	Get(ctx context.Context, id int64) (*Submission, error)
	UpdateStatus(ctx context.Context, id int64, status Status) error
}

// store implements the ContactStore interface
type store struct {
	db     *sqlx.DB
	logger logging.Logger
}

// NewStore creates a new contact store
func NewStore(db *sqlx.DB, logger logging.Logger) Store {
	return &store{
		db:     db,
		logger: logger,
	}
}

// Create stores a new contact form submission
func (s *store) Create(ctx context.Context, sub *Submission) error {
	query := `
		INSERT INTO contact_submissions (name, email, message, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	s.logger.Debug("creating contact submission",
		logging.String("email", sub.Email),
	)

	_, err := s.db.ExecContext(ctx, query, sub.Name, sub.Email, sub.Message, sub.Status)
	if err != nil {
		s.logger.Error("failed to create contact submission", logging.Error(err))
		return err
	}

	return nil
}

// Get retrieves a contact form submission by ID
func (s *store) Get(ctx context.Context, id int64) (*Submission, error) {
	query := `
		SELECT id, name, email, message, status, created_at, updated_at
		FROM contact_submissions
		WHERE id = ?
	`

	s.logger.Debug("getting contact submission by ID", logging.Int64("id", id))

	var submission Submission
	if err := s.db.GetContext(ctx, &submission, query, id); err != nil {
		s.logger.Error("failed to get contact submission", logging.Error(err))
		return nil, fmt.Errorf("failed to get contact submission: %w", err)
	}

	return &submission, nil
}

// List retrieves all contact form submissions
func (s *store) List(ctx context.Context) ([]Submission, error) {
	query := `
		SELECT id, name, email, message, status, created_at, updated_at
		FROM contact_submissions
		ORDER BY created_at DESC
	`

	s.logger.Debug("listing contact submissions")

	var submissions []Submission
	if err := s.db.SelectContext(ctx, &submissions, query); err != nil {
		s.logger.Error("failed to list contact submissions", logging.Error(err))
		return nil, fmt.Errorf("failed to list contact submissions: %w", err)
	}

	return submissions, nil
}

// UpdateStatus updates the status of a contact form submission
func (s *store) UpdateStatus(ctx context.Context, id int64, status Status) error {
	query := `
		UPDATE contact_submissions
		SET status = ?, updated_at = NOW()
		WHERE id = ?
	`

	s.logger.Debug("updating contact submission status", logging.Int64("id", id), logging.String("status", string(status)))

	result, err := s.db.ExecContext(ctx, query, status, id)
	if err != nil {
		s.logger.Error("failed to update contact submission status", logging.Error(err))
		return fmt.Errorf("failed to update contact submission status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("contact submission not found: %d", id)
	}

	return nil
}
