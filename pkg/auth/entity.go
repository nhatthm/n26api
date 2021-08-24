package auth

import "time"

// Token is an oauth2 token.
type Token string

// OAuthToken contains all relevant information to access to the service and refresh the token.
// nolint:tagliatelle
type OAuthToken struct {
	AccessToken      Token     `json:"access_token"`
	RefreshToken     Token     `json:"refresh_token"`
	ExpiresAt        time.Time `json:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

// IsExpired checks whether the access token is expired or not.
func (t OAuthToken) IsExpired(timestamp time.Time) bool {
	return t.ExpiresAt.Before(timestamp)
}

// IsRefreshable checks whether the refresh token is alive or not.
func (t OAuthToken) IsRefreshable(timestamp time.Time) bool {
	return t.RefreshExpiresAt.After(timestamp)
}
