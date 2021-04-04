package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api/pkg/auth"
	authMock "github.com/nhatthm/n26api/pkg/testkit/auth"
)

func TestTokenProvider_Token(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockProvider  authMock.TokenProviderMocker
		expectedToken auth.Token
		expectedError string
	}{
		{
			scenario: "token is string",
			mockProvider: authMock.MockTokenProvider(func(p *authMock.TokenProvider) {
				p.On("Token", context.Background()).
					Return("token", nil)
			}),
			expectedToken: auth.Token("token"),
		},
		{
			scenario: "token is auth.Token",
			mockProvider: authMock.MockTokenProvider(func(p *authMock.TokenProvider) {
				p.On("Token", context.Background()).
					Return(auth.Token("token"), nil)
			}),
			expectedToken: auth.Token("token"),
		},
		{
			scenario: "error",
			mockProvider: authMock.MockTokenProvider(func(p *authMock.TokenProvider) {
				p.On("Token", context.Background()).
					Return("", errors.New("token error"))
			}),
			expectedError: "token error",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			token, err := tc.mockProvider(t).Token(context.Background())

			assert.Equal(t, tc.expectedToken, token)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
