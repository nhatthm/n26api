package n26api

import (
	"fmt"
	"net/http"

	"github.com/nhatthm/n26api/pkg/auth"
	"github.com/nhatthm/n26api/pkg/util"
)

// RoundTripperFunc is an inline http.RoundTripper.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip satisfies RoundTripperFunc.
func (fn RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

// BasicAuthRoundTripper sets Basic Authorization header to the given request.
func BasicAuthRoundTripper(username, password string, tripper http.RoundTripper) RoundTripperFunc {
	value := fmt.Sprintf("Basic %s", util.Base64Credentials(username, password))

	return func(req *http.Request) (*http.Response, error) {
		req.Header.Add("Authorization", value)

		return tripper.RoundTrip(req)
	}
}

// BearerAuthRoundTripper sets Bearer Authorization header to the given request.
func BearerAuthRoundTripper(token string, tripper http.RoundTripper) RoundTripperFunc {
	value := fmt.Sprintf("Bearer %s", token)

	return func(req *http.Request) (*http.Response, error) {
		req.Header.Add("Authorization", value)

		return tripper.RoundTrip(req)
	}
}

// TokenRoundTripper sets Bearer Authorization header to the given request with a token given by a auth.TokenProvider.
func TokenRoundTripper(p auth.TokenProvider, tripper http.RoundTripper) RoundTripperFunc {
	return func(r *http.Request) (*http.Response, error) {
		token, err := p.Token(r.Context())
		if err != nil {
			return nil, err
		}

		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		return tripper.RoundTrip(r)
	}
}
