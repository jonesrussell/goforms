package user

// Token represents a user's authentication token
type Token struct {
	ID     string // Unique identifier for the token
	UserID string // ID of the user associated with the token
	Value  string // The actual token value
}
