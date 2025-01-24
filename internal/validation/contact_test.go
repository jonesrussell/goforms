package validation

import (
	"testing"

	"github.com/jonesrussell/goforms/internal/core/contact"
	"github.com/stretchr/testify/assert"
)

func TestValidateContact(t *testing.T) {
	tests := []struct {
		name    string
		sub     *contact.Submission
		wantErr bool
	}{
		{
			name: "valid submission",
			sub: &contact.Submission{
				Email:   "test@example.com",
				Name:    "Test User",
				Message: "Test message",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			sub: &contact.Submission{
				Name:    "Test User",
				Message: "Test message",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			sub: &contact.Submission{
				Email:   "invalid-email",
				Name:    "Test User",
				Message: "Test message",
			},
			wantErr: true,
		},
		{
			name: "missing name",
			sub: &contact.Submission{
				Email:   "test@example.com",
				Message: "Test message",
			},
			wantErr: true,
		},
		{
			name: "missing message",
			sub: &contact.Submission{
				Email: "test@example.com",
				Name:  "Test User",
			},
			wantErr: true,
		},
		{
			name:    "empty submission",
			sub:     &contact.Submission{},
			wantErr: true,
		},
		{
			name:    "nil submission",
			sub:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateContact(tt.sub)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
