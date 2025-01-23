package contact

import (
	"context"
	"fmt"
	"time"

	"github.com/jonesrussell/goforms/internal/logger"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

// Service defines the interface for contact form operations
type Service interface {
	CreateSubmission(ctx context.Context, submission *Submission) error
	ListSubmissions(ctx context.Context) ([]Submission, error)
	GetSubmission(ctx context.Context, id int64) (*Submission, error)
	UpdateSubmissionStatus(ctx context.Context, id int64, status Status) error
}

type Store interface {
	Create(ctx context.Context, submission *Submission) error
	List(ctx context.Context) ([]Submission, error)
	GetByID(ctx context.Context, id int64) (*Submission, error)
	UpdateStatus(ctx context.Context, id int64, status Status) error
}

// ServiceImpl implements the Service interface
type ServiceImpl struct {
	store  Store
	logger logger.Logger
}

// Submission represents a contact form submission
type Submission struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Message   string    `json:"message" db:"message"`
	Status    Status    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewService creates a new ServiceImpl instance
func NewService(store Store, logger logger.Logger) Service {
	return &ServiceImpl{
		store:  store,
		logger: logger,
	}
}

func (s *ServiceImpl) wrapError(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

func (s *ServiceImpl) CreateSubmission(ctx context.Context, submission *Submission) error {
	if submission.Status == "" {
		submission.Status = StatusPending
	}

	if err := s.store.Create(ctx, submission); err != nil {
		s.logger.Error("failed to create contact submission",
			logger.Error(err),
			logger.String("email", submission.Email),
		)
		return s.wrapError(err, "failed to create contact submission")
	}
	return nil
}

func (s *ServiceImpl) ListSubmissions(ctx context.Context) ([]Submission, error) {
	submissions, err := s.store.List(ctx)
	if err != nil {
		s.logger.Error("failed to list contact submissions", logger.Error(err))
		return nil, s.wrapError(err, "failed to list contact submissions")
	}
	return submissions, nil
}

func (s *ServiceImpl) GetSubmission(ctx context.Context, id int64) (*Submission, error) {
	submission, err := s.store.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get contact submission",
			logger.Error(err),
			logger.Int("id", int(id)),
		)
		return nil, s.wrapError(err, "failed to get contact submission")
	}
	return submission, nil
}

func (s *ServiceImpl) UpdateSubmissionStatus(ctx context.Context, id int64, status Status) error {
	if err := s.store.UpdateStatus(ctx, id, status); err != nil {
		s.logger.Error("failed to update contact submission status",
			logger.Error(err),
			logger.Int("id", int(id)),
			logger.String("status", string(status)),
		)
		return s.wrapError(err, "failed to update contact submission status")
	}
	return nil
}
