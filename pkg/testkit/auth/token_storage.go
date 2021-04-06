package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nhatthm/n26api/pkg/auth"
)

// TokenStorageMocker is TokenStorage mocker.
type TokenStorageMocker func(tb testing.TB) *TokenStorage

// NoMockTokenStorage is no mock TokenStorage.
var NoMockTokenStorage = MockTokenStorage()

var _ auth.TokenStorage = (*TokenStorage)(nil)

// TokenStorage is a auth.TokenStorage.
type TokenStorage struct {
	mock.Mock
}

// Get satisfies auth.TokenStorage interface.
func (t *TokenStorage) Get(ctx context.Context, key string) (auth.OAuthToken, error) {
	ret := t.Called(ctx, key)

	token := ret.Get(0)
	err := ret.Error(1)

	return token.(auth.OAuthToken), err
}

// Set satisfies auth.TokenStorage interface.
func (t *TokenStorage) Set(ctx context.Context, key string, token auth.OAuthToken) error {
	return t.Called(ctx, key, token).Error(0)
}

// mockTokenStorage mocks auth.TokenStorage interface.
func mockTokenStorage(mocks ...func(s *TokenStorage)) *TokenStorage {
	s := &TokenStorage{}

	for _, m := range mocks {
		m(s)
	}

	return s
}

// MockTokenStorage creates TokenStorage mock with cleanup to ensure all the expectations are met.
func MockTokenStorage(mocks ...func(s *TokenStorage)) TokenStorageMocker {
	return func(tb testing.TB) *TokenStorage {
		tb.Helper()

		s := mockTokenStorage(mocks...)

		tb.Cleanup(func() {
			assert.True(tb, s.Mock.AssertExpectations(tb))
		})

		return s
	}
}
