package user

import (
	"testing"

	"github.com/jonesrussell/goforms/internal/application/repositories" // Ensure this is correct
)

// Mock database setup (you can use an in-memory SQLite database for testing)
func setupTestDB() (*repositories.MockDB, error) {
	mockDB := repositories.NewMockDB() // Use the mock database
	return mockDB, nil
}

func TestCreateUser(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}

	store := &store{db: db}

	u := &User{
		Email:          "test@example.com",
		HashedPassword: "hashed_password",
		FirstName:      nil,
		LastName:       nil,
		Role:           "user",
		Active:         true,
	}

	err = store.Create(u)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify that the user was created
	var createdUser User
	err = db.Get(&createdUser, "SELECT * FROM users WHERE email = ?", u.Email)
	if err != nil {
		t.Fatalf("Failed to retrieve created user: %v", err)
	}
	if createdUser.Email != u.Email {
		t.Errorf("Expected email %s, got %s", u.Email, createdUser.Email)
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Close()

	store := &store{db: db}

	// Create a user for testing
	user := &User{
		Email:          "test@example.com",
		HashedPassword: "hashed_password",
		FirstName:      nil,
		LastName:       nil,
		Role:           "user",
		Active:         true,
	}
	store.Create(user)

	// Test retrieving the user by email
	retrievedUser, err := store.GetByEmail(user.Email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}
	if retrievedUser == nil {
		t.Fatal("Expected user to be found, but got nil")
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrievedUser.Email)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Close()

	store := &store{db: db}

	// Test retrieving a non-existent user
	retrievedUser, err := store.GetByEmail("nonexistent@example.com")
	if err != nil {
		t.Fatalf("Unexpected error when getting non-existent user: %v", err)
	}
	if retrievedUser != nil {
		t.Fatal("Expected nil for non-existent user, but got a user")
	}
}
