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
