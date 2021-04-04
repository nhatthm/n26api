package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nhatthm/n26api/pkg/auth"
)

// TokenProviderMocker is TokenProvider mocker.
type TokenProviderMocker func(t testing.TB) *TokenProvider

// NoMockTokenProvider is no mock TokenProvider.
var NoMockTokenProvider = MockTokenProvider()

var _ auth.TokenProvider = (*TokenProvider)(nil)

// TokenProvider is a auth.TokenProvider.
type TokenProvider struct {
	mock.Mock
}

// Token satisfies auth.TokenProvider.
func (i *TokenProvider) Token(ctx context.Context) (auth.Token, error) {
	ret := i.Called(ctx)

	token := ret.Get(0)
	err := ret.Error(1)

	if token, ok := token.(string); ok {
		return auth.Token(token), err
	}

	return token.(auth.Token), err
}

// mockTokenProvider mocks auth.TokenProvider interface.
func mockTokenProvider(mocks ...func(p *TokenProvider)) *TokenProvider {
	i := &TokenProvider{}

	for _, m := range mocks {
		m(i)
	}

	return i
}

// MockTokenProvider creates TokenProvider mock with cleanup to ensure all the expectations are met.
func MockTokenProvider(mocks ...func(p *TokenProvider)) TokenProviderMocker {
	return func(t testing.TB) *TokenProvider {
		i := mockTokenProvider(mocks...)

		t.Cleanup(func() {
			assert.True(t, i.Mock.AssertExpectations(t))
		})

		return i
	}
}
