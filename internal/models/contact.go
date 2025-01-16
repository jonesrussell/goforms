package models

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type ContactSubmission struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Message   string    `db:"message" json:"message"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (c *ContactSubmission) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}

	if strings.TrimSpace(c.Email) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email is required")
	}

	if strings.TrimSpace(c.Message) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "message is required")
	}

	// Email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(c.Email) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid email format")
	}

	return nil
}

type ContactStore interface {
	CreateContact(ctx context.Context, contact *ContactSubmission) error
	GetContacts(ctx context.Context) ([]ContactSubmission, error)
}

type contactStore struct {
	db DB
}

func NewContactStore(db DB) ContactStore {
	return &contactStore{db: db}
}

func (s *contactStore) CreateContact(ctx context.Context, contact *ContactSubmission) error {
	query := `
		INSERT INTO contact_submissions (name, email, message, created_at)
		VALUES (?, ?, ?, NOW())
	`
	result, err := s.db.ExecContext(ctx, query, contact.Name, contact.Email, contact.Message)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	contact.ID = id
	contact.CreatedAt = time.Now()
	return nil
}

func (s *contactStore) GetContacts(ctx context.Context) ([]ContactSubmission, error) {
	query := `
		SELECT id, name, email, message, created_at
		FROM contact_submissions
		ORDER BY created_at DESC
		LIMIT 100
	`
	var submissions []ContactSubmission
	if err := s.db.SelectContext(ctx, &submissions, query); err != nil {
		return nil, err
	}
	return submissions, nil
}
