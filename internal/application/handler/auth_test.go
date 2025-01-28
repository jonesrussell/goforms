package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/jonesrussell/goforms/internal/application/validator"
	"github.com/jonesrussell/goforms/internal/domain/user"
	mocklogging "github.com/jonesrussell/goforms/test/mocks/logging"
	mockuser "github.com/jonesrussell/goforms/test/mocks/user"
)

type signupTestCase struct {
	name       string
	input      map[string]interface{}
	wantStatus int
	wantErr    bool
}

func TestSignupHandler(t *testing.T) {
	tests := []signupTestCase{
		{
			name: "valid signup",
			input: map[string]interface{}{
				"email":     "test@example.com",
				"password":  "password123",
				"firstName": "Test",
				"lastName":  "User",
			},
			wantStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "missing required fields",
			input: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid email",
			input: map[string]interface{}{
				"email":     "notanemail",
				"password":  "password123",
				"firstName": "Test",
				"lastName":  "User",
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup
			e := echo.New()
			e.Validator = validator.NewValidator()

			mockLogger := mocklogging.NewMockLogger()
			mockUserSvc := mockuser.NewMockService()

			if !tt.wantErr {
				signup := &user.Signup{
					Email:     tt.input["email"].(string),
					Password:  tt.input["password"].(string),
					FirstName: tt.input["firstName"].(string),
					LastName:  tt.input["lastName"].(string),
				}
				mockUserSvc.ExpectSignUp(context.Background(), signup, &user.User{}, nil)
			}

			h := NewAuthHandler(mockLogger, WithUserService(mockUserSvc))

			// Create request
			jsonBytes, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("failed to marshal input: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(jsonBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Test
			err = h.handleSignup(c)

			// Assert
			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Errorf("expected status %d but got %d", tt.wantStatus, rec.Code)
			}

			if err := mockUserSvc.Verify(); err != nil {
				t.Errorf("mock expectations not met: %v", err)
			}
		})
	}
}
