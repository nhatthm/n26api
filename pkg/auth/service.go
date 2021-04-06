package auth

import "context"

const (
	// BasicAuthUsername is the username which is used in Authorization header while logging to n26api.
	BasicAuthUsername = "nativeweb"
	// BasicAuthPassword is the password which is used in Authorization header while logging to n26api.
	BasicAuthPassword = ""
)

// TokenProvider provides oauth2 token.
type TokenProvider interface {
	// Token provides a token.
	Token(ctx context.Context) (Token, error)
}

// TokenStorage persists or gets OAuthToken.
type TokenStorage interface {
	// Get gets OAuthToken from data source.
	Get(ctx context.Context, key string) (OAuthToken, error)
	// Set sets OAuthToken to data source.
	Set(ctx context.Context, key string, token OAuthToken) error
}
