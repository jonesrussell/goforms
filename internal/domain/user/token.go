package user

// Token represents a user's authentication token
type Token struct {
	ID     string // Unique identifier for the token
	UserID string // ID of the user associated with the token
	Value  string // The actual token value
}

// TokenPair represents a pair of access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
