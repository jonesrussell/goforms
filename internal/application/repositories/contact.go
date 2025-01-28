package repositories

import (
	"context"
	"fmt"

	"github.com/jonesrussell/goforms/internal/application/database"
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/contact"
)

// ContactRepository implements contact.Store interface
type ContactRepository struct {
	db     *database.Database
	logger logging.Logger
}

// NewContactRepository creates a new contact repository
func NewContactRepository(db *database.Database, logger logging.Logger) contact.Store {
	logger.Debug("creating contact repository",
		logging.Bool("db_available", db != nil),
	)
	return &ContactRepository{
		db:     db,
		logger: logger,
	}
}

// Create stores a new contact form submission
func (r *ContactRepository) Create(ctx context.Context, submission *contact.Submission) error {
	query := `
		INSERT INTO contact_submissions (name, email, message, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	r.logger.Debug("creating contact submission",
		logging.String("email", submission.Email),
		logging.String("status", string(submission.Status)),
	)

	_, err := r.db.ExecContext(ctx, query, submission.Name, submission.Email, submission.Message, submission.Status)
	if err != nil {
		r.logger.Error("failed to create contact submission", logging.Error(err))
		return fmt.Errorf("failed to create contact submission: %w", err)
	}

	return nil
}

// List returns all contact form submissions
func (r *ContactRepository) List(ctx context.Context) ([]contact.Submission, error) {
	query := `
		SELECT id, name, email, message, status, created_at, updated_at
		FROM contact_submissions
		ORDER BY created_at DESC
	`

	r.logger.Debug("listing contact submissions")

	var submissions []contact.Submission
	if err := r.db.SelectContext(ctx, &submissions, query); err != nil {
		r.logger.Error("failed to list contact submissions", logging.Error(err))
		return nil, fmt.Errorf("failed to list contact submissions: %w", err)
	}

	return submissions, nil
}

// Get returns a specific contact form submission
func (r *ContactRepository) Get(ctx context.Context, id int64) (*contact.Submission, error) {
	query := `
		SELECT id, name, email, message, status, created_at, updated_at
		FROM contact_submissions
		WHERE id = ?
	`

	r.logger.Debug("getting contact submission", logging.Int64("id", id))

	var submission contact.Submission
	if err := r.db.GetContext(ctx, &submission, query, id); err != nil {
		r.logger.Error("failed to get contact submission", logging.Error(err))
		return nil, fmt.Errorf("failed to get contact submission: %w", err)
	}

	return &submission, nil
}

// UpdateStatus updates the status of a contact form submission
func (r *ContactRepository) UpdateStatus(ctx context.Context, id int64, status contact.Status) error {
	query := `
		UPDATE contact_submissions
		SET status = ?, updated_at = NOW()
		WHERE id = ?
	`

	r.logger.Debug("updating contact submission status", logging.Int64("id", id), logging.String("status", string(status)))

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		r.logger.Error("failed to update contact submission status", logging.Error(err))
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

// Other required methods for contact.Store interface...
