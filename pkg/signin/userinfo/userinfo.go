package userinfo

// Userinfo is the base interface for user information.
type Userinfo interface {
	// Get returns the user who is associated with the token.
	Get() (*User, error)
}

// User represents user information.
type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}
