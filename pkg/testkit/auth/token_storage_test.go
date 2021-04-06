package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api/pkg/auth"
	authMock "github.com/nhatthm/n26api/pkg/testkit/auth"
)

func TestTokenStorage_Get(t *testing.T) {
	t.Parallel()

	timestamp := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)

	key := "username"
	token := auth.OAuthToken{
		AccessToken:      "access",
		RefreshToken:     "refresh",
		ExpiresAt:        timestamp,
		RefreshExpiresAt: timestamp,
	}

	testCases := []struct {
		scenario      string
		mockStorage   authMock.TokenStorageMocker
		expectedToken auth.OAuthToken
		expectedError string
	}{
		{
			scenario: "everything is nil",
			mockStorage: authMock.MockTokenStorage(func(s *authMock.TokenStorage) {
				s.On("Get", context.Background(), key).
					Return(auth.OAuthToken{}, nil)
			}),
		},
		{
			scenario: "token is not nil",
			mockStorage: authMock.MockTokenStorage(func(s *authMock.TokenStorage) {
				s.On("Get", context.Background(), key).
					Return(token, nil)
			}),
			expectedToken: token,
		},
		{
			scenario: "error is not nil",
			mockStorage: authMock.MockTokenStorage(func(s *authMock.TokenStorage) {
				s.On("Get", context.Background(), key).
					Return(auth.OAuthToken{}, errors.New("get error"))
			}),
			expectedError: "get error",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := tc.mockStorage(t)
			token, err := s.Get(context.Background(), key)

			assert.Equal(t, tc.expectedToken, token)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestTokenStorage_Set(t *testing.T) {
	t.Parallel()

	timestamp := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)

	key := "username"
	token := auth.OAuthToken{
		AccessToken:      "access",
		RefreshToken:     "refresh",
		ExpiresAt:        timestamp,
		RefreshExpiresAt: timestamp,
	}

	testCases := []struct {
		scenario      string
		mockStorage   authMock.TokenStorageMocker
		expectedError string
	}{
		{
			scenario: "error",
			mockStorage: authMock.MockTokenStorage(func(s *authMock.TokenStorage) {
				s.On("Set", context.Background(), key, token).
					Return(errors.New("set error"))
			}),
			expectedError: "set error",
		},
		{
			scenario: "success",
			mockStorage: authMock.MockTokenStorage(func(s *authMock.TokenStorage) {
				s.On("Set", context.Background(), key, token).
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := tc.mockStorage(t)
			err := s.Set(context.Background(), key, token)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
