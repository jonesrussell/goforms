package token

// Rename TokenRepository to Repository
type Repository interface {
	SaveToken(userID string, token string) error
	GetToken(userID string) (string, error)
	IsTokenBlacklisted(token string) bool
	BlacklistToken(token string) error
}
