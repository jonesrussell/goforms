package user

import (
	"testing"
)

func TestInMemoryTokenRepository_SaveToken(t *testing.T) {
	repo := NewInMemoryTokenRepository()

	err := repo.SaveToken("user1", "token123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify that the token was saved
	token, err := repo.GetToken("user1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token != "token123" {
		t.Errorf("expected token123, got %s", token)
	}
}

func TestInMemoryTokenRepository_GetToken_NotFound(t *testing.T) {
	repo := NewInMemoryTokenRepository()

	// Attempt to get a token that doesn't exist
	token, err := repo.GetToken("nonexistent_user")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token != "" { // Expecting an empty string since the token does not exist
		t.Errorf("expected empty token, got %s", token)
	}
}

func TestInMemoryTokenRepository_IsTokenBlacklisted(t *testing.T) {
	repo := NewInMemoryTokenRepository()

	// Save a token and then check if it's blacklisted
	err := repo.SaveToken("user1", "token123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.IsTokenBlacklisted("token123") {
		t.Fatal("expected token to not be blacklisted initially")
	}

	// Blacklist the token
	err = repo.BlacklistToken("token123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify that the token is blacklisted
	if !repo.IsTokenBlacklisted("token123") {
		t.Fatal("expected token to be blacklisted")
	}
}

func TestInMemoryTokenRepository_BlacklistToken(t *testing.T) {
	repo := NewInMemoryTokenRepository()

	err := repo.BlacklistToken("token123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify that the token is blacklisted
	if !repo.IsTokenBlacklisted("token123") {
		t.Fatal("expected token to be blacklisted")
	}
}
