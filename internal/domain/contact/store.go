package contact

import (
	"context"
	"fmt"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories/database"
	"github.com/jonesrussell/goforms/internal/domain/common"
)

// Repository defines the methods for contact data access.
type Repository interface {
	CreateSubmission(ctx context.Context, sub *common.Submission) error
	GetSubmission(ctx context.Context, id int64) (*common.Submission, error)
	ListSubmissions(ctx context.Context) ([]common.Submission, error)
	UpdateSubmissionStatus(ctx context.Context, id int64, status string) error
	// Add other necessary methods
}

// StoreImpl implements the Repository interface.
type StoreImpl struct {
	db     *database.DB
	logger logging.Logger
}

// CreateSubmission stores a new contact form submission
func (s *StoreImpl) CreateSubmission(ctx context.Context, sub *common.Submission) error {
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

// GetSubmission retrieves a contact form submission by ID
func (s *StoreImpl) GetSubmission(ctx context.Context, id int64) (*common.Submission, error) {
	query := `
		SELECT id, name, email, message, status, created_at, updated_at
		FROM contact_submissions
		WHERE id = ?
	`

	s.logger.Debug("getting contact submission by ID", logging.Int64("id", id))

	var submission common.Submission
	if err := s.db.GetContext(ctx, &submission, query, id); err != nil {
		s.logger.Error("failed to get contact submission", logging.Error(err))
		return nil, fmt.Errorf("failed to get contact submission: %w", err)
	}

	return &submission, nil
}

// ListSubmissions retrieves all contact form submissions
func (s *StoreImpl) ListSubmissions(ctx context.Context) ([]common.Submission, error) {
	query := `
		SELECT id, name, email, message, status, created_at, updated_at
		FROM contact_submissions
		ORDER BY created_at DESC
	`

	s.logger.Debug("listing contact submissions")

	var submissions []common.Submission
	if err := s.db.SelectContext(ctx, &submissions, query); err != nil {
		s.logger.Error("failed to list contact submissions", logging.Error(err))
		return nil, fmt.Errorf("failed to list contact submissions: %w", err)
	}

	return submissions, nil
}

// UpdateSubmissionStatus updates the status of a contact form submission
func (s *StoreImpl) UpdateSubmissionStatus(ctx context.Context, id int64, status string) error {
	query := `
		UPDATE contact_submissions
		SET status = ?, updated_at = NOW()
		WHERE id = ?
	`

	s.logger.Debug("updating contact submission status", logging.Int64("id", id), logging.String("status", status))

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

// NewRepository creates a new contact repository
func NewRepository(db *database.DB, logger logging.Logger) Repository {
	return &StoreImpl{
		db:     db,
		logger: logger,
	}
}
