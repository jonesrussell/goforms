package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryTokenRepository_SaveToken(t *testing.T) {
	repo := NewInMemoryTokenRepository()

	err := repo.SaveToken("user1", "token123")
	assert.NoError(t, err)

	// Verify that the token was saved
	token, err := repo.GetToken("user1")
	assert.NoError(t, err)
	assert.Equal(t, "token123", token)
}

func TestInMemoryTokenRepository_GetToken_NotFound(t *testing.T) {
	repo := NewInMemoryTokenRepository()

	// Attempt to get a token that doesn't exist
	token, err := repo.GetToken("nonexistent_user")
	assert.NoError(t, err)
	assert.Empty(t, token) // Expecting an empty string since the token does not exist
}

func TestInMemoryTokenRepository_IsTokenBlacklisted(t *testing.T) {
	repo := NewInMemoryTokenRepository()

	// Save a token and then check if it's blacklisted
	repo.SaveToken("user1", "token123")
	assert.False(t, repo.IsTokenBlacklisted("token123")) // Should not be blacklisted initially

	// Blacklist the token
	repo.BlacklistToken("token123")
	assert.True(t, repo.IsTokenBlacklisted("token123")) // Should be blacklisted now
}

func TestInMemoryTokenRepository_BlacklistToken(t *testing.T) {
	repo := NewInMemoryTokenRepository()

	// Blacklist a token
	err := repo.BlacklistToken("token123")
	assert.NoError(t, err)

	// Verify that the token is blacklisted
	assert.True(t, repo.IsTokenBlacklisted("token123"))
}
