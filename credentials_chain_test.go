package n26api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api/pkg/testkit"
)

func TestChainCredentialsProvider(t *testing.T) {
	testCases := []struct {
		scenario         string
		mockProviders    []testkit.CredentialsProviderMocker
		expectedUsername string
		expectedPassword string
	}{
		{
			scenario: "no provider",
		},
		{
			scenario:         "chained providers",
			mockProviders:    provideCredentialsProviders(),
			expectedUsername: "username",
			expectedPassword: "password",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			p := newChainCredentialsProvider()

			for _, mockProvider := range tc.mockProviders {
				p.chain(mockProvider(t))
			}

			assert.Equal(t, tc.expectedUsername, p.Username())
			assert.Equal(t, tc.expectedPassword, p.Password())
		})
	}
}

func TestChainCredentialsProviders(t *testing.T) {
	mocks := provideCredentialsProviders()
	providers := make([]CredentialsProvider, 0, len(mocks))

	for _, m := range mocks {
		providers = append(providers, m(t))
	}

	p := chainCredentialsProviders(providers...)

	assert.Equal(t, "username", p.Username())
	assert.Equal(t, "password", p.Password())
}

func provideCredentialsProviders() []testkit.CredentialsProviderMocker {
	return []testkit.CredentialsProviderMocker{
		// This provider should not be called.
		testkit.NoMockCredentialsProvider,
		// Has Username, no password.
		testkit.MockCredentialsProvider(func(p *testkit.CredentialsProvider) {
			p.On("Username").Return("username")
			// Password() should not be called because the 2nd provider gets called first.
		}),
		// No Username, has password.
		testkit.MockCredentialsProvider(func(p *testkit.CredentialsProvider) {
			p.On("Username").Return("")
			p.On("Password").Return("password")
		}),
		// No Username, no password.
		testkit.MockCredentialsProvider(func(p *testkit.CredentialsProvider) {
			p.On("Username").Return("")
			p.On("Password").Return("")
		}),
	}
}
