package token

// TokenRepository defines the interface for token management
type TokenRepository interface {
	SaveToken(userID string, token string) error
	GetToken(userID string) (string, error)
	IsTokenBlacklisted(token string) bool
	BlacklistToken(token string) error
}
