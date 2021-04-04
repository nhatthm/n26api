package n26api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api/pkg/auth"
	authMock "github.com/nhatthm/n26api/pkg/testkit/auth"
)

func TestChainTokenProvider_Token(t *testing.T) {
	testCases := []struct {
		scenario      string
		mockProviders []authMock.TokenProviderMocker
		expectedToken auth.Token
		expectedError string
	}{
		{
			scenario: "no provider",
		},
		{
			scenario:      "has error",
			mockProviders: provideTokenProviders("", errors.New("failed to get token")),
			expectedError: "failed to get token",
		},
		{
			scenario:      "success",
			mockProviders: provideTokenProviders("token", nil),
			expectedToken: "token",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			p := newChainTokenProvider()

			for _, mockProvider := range tc.mockProviders {
				p.chain(mockProvider(t))
			}

			token, err := p.Token(context.Background())

			assert.Equal(t, tc.expectedToken, token)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func provideTokenProviders(token auth.Token, err error) []authMock.TokenProviderMocker {
	return []authMock.TokenProviderMocker{
		// This provider should not be called.
		authMock.NoMockTokenProvider,
		// Has Token.
		authMock.MockTokenProvider(func(p *authMock.TokenProvider) {
			p.On("Token", context.Background()).Return(token, err)
		}),
		// No Token.
		authMock.MockTokenProvider(func(p *authMock.TokenProvider) {
			p.On("Token", context.Background()).Return(auth.Token(""), nil)
		}),
	}
}
