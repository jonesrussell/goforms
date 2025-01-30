package user_test

import (
	"fmt"
	"testing"

	"github.com/jonesrussell/goforms/internal/domain/common"
	"github.com/jonesrussell/goforms/internal/domain/user"
)

func TestStore_Create(t *testing.T) {
	store := user.NewMockStore() // Use the mock store

	user := &common.User{
		ID:        1, // Assign an ID for testing
		Email:     "test@example.com",
		Password:  "securepassword",
		FirstName: nil,
		LastName:  nil,
		Role:      "user",
		Active:    true,
	}

	err := store.Create(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify that the user was created
	retrievedUser, err := store.GetByEmail(user.Email)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrievedUser == nil {
		t.Fatal("expected user to be created, got nil")
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, retrievedUser.Email)
	}
}

func TestStore_Get(t *testing.T) {
	store := user.NewMockStore()

	user := &common.User{
		Email:     "test@example.com",
		Password:  "securepassword",
		FirstName: nil,
		LastName:  nil,
		Role:      "user",
		Active:    true,
	}

	// Create the user first
	err := store.Create(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Retrieve the user by ID
	retrievedUser, err := store.Get(user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrievedUser == nil {
		t.Fatal("expected user to be retrieved, got nil")
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, retrievedUser.Email)
	}
}

func TestStore_GetByEmail_NotFound(t *testing.T) {
	store := user.NewMockStore()

	// Attempt to get a user by email that doesn't exist
	retrievedUser, err := store.GetByEmail("nonexistent@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrievedUser != nil {
		t.Fatal("expected user to be nil, got non-nil")
	}
}

func TestStore_Update(t *testing.T) {
	store := user.NewMockStore()

	user := &common.User{
		Email:     "test@example.com",
		Password:  "securepassword",
		FirstName: nil,
		LastName:  nil,
		Role:      "user",
		Active:    true,
	}

	// Create the user first
	err := store.Create(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Update the user's email
	user.Email = "updated@example.com"
	err = store.Update(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify the update
	updatedUser, err := store.GetByEmail(user.Email)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updatedUser.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, updatedUser.Email)
	}
}

func TestStore_Delete(t *testing.T) {
	store := user.NewMockStore()

	user := &common.User{
		Email:     "test@example.com",
		Password:  "securepassword",
		FirstName: nil,
		LastName:  nil,
		Role:      "user",
		Active:    true,
	}

	// Create the user first
	err := store.Create(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Delete the user
	err = store.Delete(user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify the user is deleted
	deletedUser, err := store.Get(user.ID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if deletedUser != nil {
		t.Fatal("expected user to be nil, got non-nil")
	}
}

func TestStore_List(t *testing.T) {
	store := user.NewMockStore() // Use the mock store

	// Create multiple users
	for i := 0; i < 5; i++ {
		user := &common.User{
			Email:     fmt.Sprintf("user%d@example.com", i),
			Password:  "securepassword",
			FirstName: nil,
			LastName:  nil,
			Role:      "user",
			Active:    true,
		}
		err := store.Create(user)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	}

	// List all users
	users, err := store.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 5 {
		t.Errorf("expected 5 users, got %d", len(users))
	}
}
