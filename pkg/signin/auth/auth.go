package auth

import (
	"time"
)

// Auth is the base interface for authorization workflow for oauth
// api like google, microsoft, facebook etc.
type Auth interface {
	// GetURL gets the url for authorization workflow.
	GetURL(redirect, state string) string

	// Exchange converts authorization code into token.
	Exchange(code string) (*Token, error)

	// Refresh renews access token using refresh token.
	Refresh(refreshToken string) (*Token, error)

	// Revoke cancels the access token and reset the authorization workflow.
	Revoke(accessToken string) error
}

// Token represents the credentials used to authorize.
type Token struct {
	AccessToken  string    `json:"accessToken,omitempty"`
	RefreshToken string    `json:"refreshToken,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
}
