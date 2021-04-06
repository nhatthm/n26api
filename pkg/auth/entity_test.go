package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOAuthToken_IsExpired(t *testing.T) {
	t.Parallel()

	timestampBefore := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	timestamp := timestampBefore.Add(time.Second)
	timestampAfter := timestamp.Add(time.Second)

	testCases := []struct {
		scenario  string
		expiresAt time.Time
		expected  bool
	}{
		{
			scenario:  "access token is expires when ExpiresAt is before the timestamp",
			expiresAt: timestampBefore,
			expected:  true,
		},
		{
			scenario:  "access token is alive when ExpiresAt is equal to the timestamp",
			expiresAt: timestamp,
			expected:  false,
		},
		{
			scenario:  "access token is alive when ExpiresAt is after the timestamp",
			expiresAt: timestampAfter,
			expected:  false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			token := OAuthToken{ExpiresAt: tc.expiresAt}

			assert.Equal(t, tc.expected, token.IsExpired(timestamp))
		})
	}
}

func TestOAuthToken_IsRefreshable(t *testing.T) {
	t.Parallel()

	timestampBefore := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	timestamp := timestampBefore.Add(time.Second)
	timestampAfter := timestamp.Add(time.Second)

	testCases := []struct {
		scenario  string
		expiresAt time.Time
		expected  bool
	}{
		{
			scenario:  "refresh token is expired when RefreshExpiresAt is before the timestamp",
			expiresAt: timestampBefore,
			expected:  false,
		},
		{
			scenario:  "refresh token is expired when RefreshExpiresAt is equal to the timestamp",
			expiresAt: timestamp,
			expected:  false,
		},
		{
			scenario:  "refresh token is alive when RefreshExpiresAt is after the timestamp",
			expiresAt: timestampAfter,
			expected:  true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			token := OAuthToken{RefreshExpiresAt: tc.expiresAt}

			assert.Equal(t, tc.expected, token.IsRefreshable(timestamp))
		})
	}
}
