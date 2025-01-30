package user

import (
	"context"
	"errors"
)

// InMemoryTokenRepository is an in-memory implementation of TokenRepository
type InMemoryTokenRepository struct {
	tokens map[string]string
}

// Delete implements TokenRepository.
func (repo *InMemoryTokenRepository) Delete(tokenID string) error {
	if _, exists := repo.tokens[tokenID]; !exists {
		return errors.New("token not found")
	}
	delete(repo.tokens, tokenID)
	return nil
}

// Find implements TokenRepository.
func (repo *InMemoryTokenRepository) Find(tokenID string) (Token, error) {
	panic("unimplemented")
}

// Login implements TokenRepository.
func (repo *InMemoryTokenRepository) Login(ctx context.Context, login *Login) (*TokenPair, error) {
	panic("unimplemented")
}

// Logout implements TokenRepository.
func (repo *InMemoryTokenRepository) Logout(ctx context.Context, token string) error {
	panic("unimplemented")
}

// Save implements TokenRepository.
func (repo *InMemoryTokenRepository) Save(token Token) error {
	panic("unimplemented")
}

// ValidateToken implements TokenRepository.
func (repo *InMemoryTokenRepository) ValidateToken(token string) (string, error) {
	panic("unimplemented")
}

// NewInMemoryTokenRepository creates a new InMemoryTokenRepository
func NewInMemoryTokenRepository() *InMemoryTokenRepository {
	return &InMemoryTokenRepository{
		tokens: make(map[string]string),
	}
}

// SaveToken saves a token for a user
func (repo *InMemoryTokenRepository) SaveToken(userID string, token string) error {
	repo.tokens[userID] = token
	return nil
}

// GetToken retrieves a token for a user
func (repo *InMemoryTokenRepository) GetToken(userID string) (string, error) {
	token, exists := repo.tokens[userID]
	if !exists {
		return "", nil // or an error if preferred
	}
	return token, nil
}

// Implement TokenRepository methods
func (repo *InMemoryTokenRepository) IsTokenBlacklisted(token string) bool {
	_, exists := repo.tokens[token]
	return exists
}

func (repo *InMemoryTokenRepository) BlacklistToken(token string) error {
	repo.tokens[token] = token
	return nil
}
