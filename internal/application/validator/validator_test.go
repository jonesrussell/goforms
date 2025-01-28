package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

type TestStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=0,lte=130"`
}

type TestSignup struct {
	Email     string `validate:"required,email"`
	Password  string `validate:"required,min=8"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
}

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	if v == nil {
		t.Error("NewValidator() returned nil")
	}
	if _, ok := interface{}(v).(*CustomValidator); !ok {
		t.Error("NewValidator() did not return *CustomValidator")
	}
}

func TestValidate(t *testing.T) {
	v := NewValidator()

	t.Run("valid struct", func(t *testing.T) {
		test := TestStruct{
			Name:  "Test User",
			Email: "test@example.com",
			Age:   25,
		}
		err := v.Validate(test)
		if err != nil {
			t.Errorf("Validate() error = %v, want nil", err)
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		test := TestStruct{
			Name:  "Test User",
			Email: "invalid-email",
			Age:   25,
		}
		err := v.Validate(test)
		if err == nil {
			t.Error("Validate() error = nil, want error")
		}
		validationErr, ok := err.(validator.ValidationErrors)
		if !ok {
			t.Error("Validate() error is not validator.ValidationErrors")
		}
		if !containsTag(validationErr[0].Tag(), "email") {
			t.Error("Validate() error does not contain 'email' tag")
		}
	})

	t.Run("missing required field", func(t *testing.T) {
		test := TestStruct{
			Email: "test@example.com",
			Age:   25,
		}
		err := v.Validate(test)
		if err == nil {
			t.Error("Validate() error = nil, want error")
		}
		validationErr, ok := err.(validator.ValidationErrors)
		if !ok {
			t.Error("Validate() error is not validator.ValidationErrors")
		}
		if !containsTag(validationErr[0].Tag(), "required") {
			t.Error("Validate() error does not contain 'required' tag")
		}
	})

	t.Run("age out of range", func(t *testing.T) {
		test := TestStruct{
			Name:  "Test User",
			Email: "test@example.com",
			Age:   150,
		}
		err := v.Validate(test)
		if err == nil {
			t.Error("Validate() error = nil, want error")
		}
		validationErr, ok := err.(validator.ValidationErrors)
		if !ok {
			t.Error("Validate() error is not validator.ValidationErrors")
		}
		if !containsTag(validationErr[0].Tag(), "lte") {
			t.Error("Validate() error does not contain 'lte' tag")
		}
	})

	t.Run("non-struct value", func(t *testing.T) {
		err := v.Validate("not a struct")
		if err == nil {
			t.Error("Validate() error = nil, want error")
		}
	})
}

func TestValidator(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name    string
		input   TestSignup
		wantErr bool
	}{
		{
			name: "valid signup",
			input: TestSignup{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			input: TestSignup{
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			input: TestSignup{
				Email:     "notanemail",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			wantErr: true,
		},
		{
			name: "missing first name",
			input: TestSignup{
				Email:    "test@example.com",
				Password: "password123",
				LastName: "User",
			},
			wantErr: true,
		},
		{
			name: "missing last name",
			input: TestSignup{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
			},
			wantErr: true,
		},
		{
			name: "password too short",
			input: TestSignup{
				Email:     "test@example.com",
				Password:  "short",
				FirstName: "Test",
				LastName:  "User",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func containsTag(tag, substr string) bool {
	return tag == substr
}
