package user

// InMemoryTokenRepository is a simple in-memory implementation of TokenRepository
type InMemoryTokenRepository struct {
	tokens map[string]Token
}

// NewInMemoryTokenRepository creates a new InMemoryTokenRepository
func NewInMemoryTokenRepository() *InMemoryTokenRepository {
	return &InMemoryTokenRepository{tokens: make(map[string]Token)}
}

// Implement TokenRepository methods
func (repo *InMemoryTokenRepository) IsTokenBlacklisted(token string) bool {
	// Check if the token is blacklisted
	_, exists := repo.tokens[token]
	return exists
}

func (repo *InMemoryTokenRepository) BlacklistToken(token string) error {
	// Add the token to the blacklist
	repo.tokens[token] = Token{Value: token} // Assuming Token has a Value field
	return nil
}

// Implement other methods as needed...
