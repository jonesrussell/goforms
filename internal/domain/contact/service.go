package contact

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/domain/common"
)

type Service interface {
	Submit(ctx context.Context, sub *common.Submission) error
	ListSubmissions(ctx context.Context) ([]common.Submission, error)
	GetSubmission(ctx context.Context, id int64) (*common.Submission, error)
	UpdateSubmissionStatus(ctx context.Context, id int64, status string) error
}

// Ensure Service implements ContactService interface
var _ Service = (*ServiceImpl)(nil)

// ServiceImpl handles contact form business logic
type ServiceImpl struct {
	store  Store
	logger logging.Logger
}

// Submit handles a new contact form submission
func (s *ServiceImpl) Submit(ctx context.Context, sub *common.Submission) error {
	if sub.Name == "" {
		return ErrNameRequired
	}
	if sub.Email == "" {
		return ErrEmailRequired
	}
	if sub.Message == "" {
		return ErrMessageRequired
	}

	sub.Status = common.StatusPending // Use the constant from common
	sub.CreatedAt = time.Now()
	sub.UpdatedAt = sub.CreatedAt

	if err := s.store.Create(ctx, sub); err != nil {
		s.logger.Error("failed to create submission",
			logging.Error(err),
			logging.String("email", sub.Email),
		)
		return err
	}

	s.logger.Info("submission created",
		logging.String("email", sub.Email),
		logging.String("status", string(sub.Status)), // Convert to string if needed
	)

	return nil
}

// ListSubmissions returns all contact form submissions
func (s *ServiceImpl) ListSubmissions(ctx context.Context) ([]common.Submission, error) {
	submissions, err := s.store.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list submissions: %w", err)
	}
	return submissions, nil
}

// GetSubmission returns a specific contact form submission
func (s *ServiceImpl) GetSubmission(ctx context.Context, id int64) (*common.Submission, error) {
	submission, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get submission: %w", err)
	}
	return submission, nil
}

// UpdateSubmissionStatus updates the status of a contact form submission
func (s *ServiceImpl) UpdateSubmissionStatus(ctx context.Context, id int64, status string) error {
	var statusEnum common.Status

	// Convert string to common.Status
	switch status {
	case string(common.StatusPending):
		statusEnum = common.StatusPending
	case string(common.StatusApproved):
		statusEnum = common.StatusApproved
	case string(common.StatusRejected):
		statusEnum = common.StatusRejected
	default:
		return ErrInvalidStatus
	}

	// Convert statusEnum to string when calling UpdateStatus
	if err := s.store.UpdateStatus(ctx, id, statusEnum); err != nil {
		s.logger.Error("failed to update submission status",
			logging.Error(err),
			logging.String("id", strconv.FormatInt(id, 10)),
			logging.String("status", string(statusEnum)),
		)
		return fmt.Errorf("failed to update submission status: %w", err)
	}

	s.logger.Info("submission status updated",
		logging.String("id", strconv.FormatInt(id, 10)),
		logging.String("status", string(statusEnum)),
	)

	return nil
}

// NewService creates a new service instance.
func NewService(store Store) Service {
	return &ServiceImpl{store: store}
}
