package user

// InMemoryTokenRepository is an in-memory implementation of TokenRepository
type InMemoryTokenRepository struct {
	tokens map[string]string
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
