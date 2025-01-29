package user

import (
	"github.com/jonesrussell/goforms/internal/application/logging"
	"github.com/jonesrussell/goforms/internal/application/repositories/database"
)

// tokenRepository implements the TokenRepository interface
type tokenRepository struct {
	logger logging.Logger
	db     *database.DB
}

// BlacklistToken marks a token as blacklisted in the database
func (r *tokenRepository) BlacklistToken(token string) error {
	// Implement your logic to blacklist the token in the database
	_, err := r.db.Exec("INSERT INTO blacklisted_tokens (token) VALUES (?)", token)
	if err != nil {
		r.logger.Error("Failed to blacklist token", logging.Error(err))
		return err
	}
	return nil
}
