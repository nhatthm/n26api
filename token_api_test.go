package n26api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nhatthm/n26api/pkg/auth"
	"github.com/nhatthm/n26api/pkg/testkit"
	authMock "github.com/nhatthm/n26api/pkg/testkit/auth"
)

func TestApiTokenProvider_GetToken(t *testing.T) {
	t.Parallel()

	username := "john.doe"
	password := "jane.doe"
	cred := Credentials(username, password)
	deviceID := uuid.New()

	configureTimeout := func(_ *testing.T, p *apiTokenProvider) {
		// Timeout: 175ms.
		// Wait time: 50ms.
		// Max calls: 3.
		p.WithMFATimeout(175 * time.Millisecond).
			WithMFAWait(50 * time.Millisecond)
	}

	testCases := []struct {
		scenario      string
		mockServer    testkit.ServerMocker
		configure     func(t *testing.T, p *apiTokenProvider)
		expectedError string
	}{
		{
			scenario:   "could not get token from token",
			mockServer: testkit.MockEmptyServer(),
			configure: func(t *testing.T, p *apiTokenProvider) { // nolint: thelper
				s := authMock.MockTokenStorage(func(s *authMock.TokenStorage) {
					s.On("Get", context.Background(), "john.doe").
						Return(auth.OAuthToken{}, errors.New("get token error"))
				})(t)

				p.WithStorage(s)
			},
			expectedError: "could not get token from storage: get token error",
		},
		{
			scenario: "unexpected response",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthPasswordLoginUnexpectedResponse(username, password, deviceID),
			),
			expectedError: "unexpected response: EOF",
		},
		{
			scenario: "wrong credentials",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthPasswordLoginFailureWrongCredentials(username, password, deviceID),
			),
			expectedError: "wrong credentials",
		},
		{
			scenario: "too many login attempts",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthPasswordLoginFailureTooManyAttempts(username, password, deviceID),
			),
			expectedError: "too many login attempts",
		},
		{
			scenario: "internal error",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthPasswordLoginFailure(username, password, deviceID),
			),
			expectedError: "unexpected response: unexpected response status: 500 Internal Server Error",
		},
		{
			scenario: "mfa challenge with invalid token",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthPasswordLoginSuccess(username, password, deviceID),
				testkit.WithAuthMFAChallengeFailureInvalidToken(),
			),
			expectedError: "could not challenge mfa",
		},
		{
			scenario: "mfa challenge failure",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthPasswordLoginSuccess(username, password, deviceID),
				testkit.WithAuthMFAChallengeFailure(),
			),
			expectedError: "failed to challenge mfa: unexpected response status: 500 Internal Server Error",
		},
		{
			scenario: "mfa timeout",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthPasswordLoginSuccess(username, password, deviceID),
				testkit.WithAuthMFAChallengeSuccess(),
				testkit.WithAuthConfirmLoginFailure(),
				testkit.WithAuthConfirmLoginFailureInvalidToken(2),
			),
			configure:     configureTimeout,
			expectedError: "could not confirm login",
		},
		{
			scenario: "could not set token",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthPasswordLoginSuccess(username, password, deviceID),
				testkit.WithAuthMFAChallengeSuccess(),
				testkit.WithAuthConfirmLoginFailureInvalidToken(2),
				testkit.WithAuthConfirmLoginSuccess(),
			),
			configure: func(t *testing.T, p *apiTokenProvider) { // nolint: thelper
				configureTimeout(t, p)

				s := authMock.MockTokenStorage(func(s *authMock.TokenStorage) {
					s.On("Get", context.Background(), "john.doe").
						Return(auth.OAuthToken{}, nil)

					s.On("Set", context.Background(), "john.doe", mock.Anything).
						Return(errors.New("set token error"))
				})(t)

				p.WithStorage(s)
			},
			expectedError: "could not persist token to storage: set token error",
		},
		{
			scenario: "success",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthPasswordLoginSuccess(username, password, deviceID),
				testkit.WithAuthMFAChallengeSuccess(),
				testkit.WithAuthConfirmLoginFailureInvalidToken(2),
				testkit.WithAuthConfirmLoginSuccess(),
			),
			configure: configureTimeout,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := tc.mockServer(t)
			p := newAPITokenProvider(cred, deviceID).
				WithBaseURL(s.URL()).
				WithTimeout(time.Second)

			if tc.configure != nil {
				tc.configure(t, p)
			}

			token, err := p.Token(context.Background())

			if tc.expectedError == "" {
				assert.Equal(t, s.AccessToken(), token)
				assert.NoError(t, err)
			} else {
				assert.Empty(t, token)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestApiTokenProvider_GetTokenFromCache(t *testing.T) {
	t.Parallel()

	username := "john.doe"
	password := "jane.doe"
	cred := Credentials(username, password)
	deviceID := uuid.New()
	timestamp := time.Now()

	s := testkit.MockEmptyServer(
		testkit.WithAuthSuccess(username, password, deviceID),
	)(t)

	c := testkit.MockClock(func(c *testkit.Clock) {
		c.On("Now").Return(timestamp).Once()
		// 2nd is after 4 minutes to check TTL.
		c.On("Now").Return(timestamp.Add(4 * time.Minute)).Once()
	})(t)

	p := newAPITokenProvider(cred, deviceID).
		WithBaseURL(s.URL()).
		WithTimeout(time.Second).
		WithMFAWait(time.Millisecond).
		WithClock(c)

	// 1st try.
	token, err := p.Token(context.Background())

	assert.Equal(t, s.AccessToken(), token)
	assert.NotEmpty(t, string(token))
	assert.NoError(t, err)

	// 2nd try.
	token, err = p.Token(context.Background())

	assert.Equal(t, s.AccessToken(), token)
	assert.NotEmpty(t, string(token))
	assert.NoError(t, err)
}

func TestApiTokenProvider_RefreshToken(t *testing.T) {
	t.Parallel()

	username := "john.doe"
	password := "jane.doe"
	cred := Credentials(username, password)
	deviceID := uuid.New()
	timestamp := time.Now()
	refreshTTL := time.Hour

	mockClock := testkit.MockClock(func(c *testkit.Clock) {
		// 1st step: Get token.
		c.On("Now").Return(timestamp).Once()
		// 2nd step: Refresh token.
		c.On("Now").Return(timestamp.Add(refreshTTL - time.Minute)).Once()
	})

	testCases := []struct {
		scenario       string
		mockServer     testkit.ServerMocker
		configureStep2 func(t *testing.T, p *apiTokenProvider, s *testkit.Server)
		expectedError  string
	}{
		{
			scenario: "failed to refresh token",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthSuccess(username, password, deviceID),
				testkit.WithAuthRefreshTokenFailure(),
			),
			expectedError: "failed to refresh token: unexpected response status: 500 Internal Server Error",
		},
		{
			scenario: "invalid token",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthSuccess(username, password, deviceID),
				testkit.WithAuthRefreshTokenFailureInvalidToken(),
				testkit.WithAuthSuccess(username, password, deviceID),
			),
		},
		{
			scenario: "could not persist token",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthSuccess(username, password, deviceID),
				testkit.WithAuthRefreshTokenSuccess(),
			),
			configureStep2: func(t *testing.T, p *apiTokenProvider, svr *testkit.Server) { // nolint: thelper
				s := authMock.MockTokenStorage(func(s *authMock.TokenStorage) {
					s.On("Get", context.Background(), "john.doe").
						Return(
							auth.OAuthToken{
								AccessToken:      svr.AccessToken(),
								RefreshToken:     svr.RefreshToken(),
								ExpiresAt:        timestamp.Add(time.Minute),
								RefreshExpiresAt: timestamp.Add(time.Hour),
							},
							nil,
						)

					s.On("Set", context.Background(), "john.doe", mock.Anything).
						Return(errors.New("set token error"))
				})(t)

				p.WithStorage(s)
			},
			expectedError: "could not persist token to storage: set token error",
		},
		{
			scenario: "success",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthSuccess(username, password, deviceID),
				testkit.WithAuthRefreshTokenSuccess(),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := tc.mockServer(t)
			p := newAPITokenProvider(cred, deviceID).
				WithBaseURL(s.URL()).
				WithTimeout(time.Second).
				WithMFAWait(time.Millisecond).
				WithRefreshTTL(refreshTTL).
				WithClock(mockClock(t))

			// 1st step: get the token.
			token1, err := p.Token(context.Background())

			assert.Equal(t, s.AccessToken(), token1)
			assert.NotEmpty(t, string(token1))
			assert.NoError(t, err)

			// 2nd step: refresh the token.
			if tc.configureStep2 != nil {
				tc.configureStep2(t, p, s)
			}

			token2, err := p.Token(context.Background())

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, s.AccessToken(), token2)
				assert.NotEqual(t, token1, token2)
			} else {
				assert.Empty(t, string(token2))
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestApiTokenProvider_TokenExpired(t *testing.T) {
	t.Parallel()

	username := "john.doe"
	password := "jane.doe"
	cred := Credentials(username, password)
	deviceID := uuid.New()
	timestamp := time.Now()
	refreshTTL := time.Hour

	s := testkit.MockEmptyServer(
		testkit.WithAuthSuccess(username, password, deviceID),
		testkit.WithAuthSuccess(username, password, deviceID),
	)(t)

	c := testkit.MockClock(func(c *testkit.Clock) {
		// 1st step: Get token.
		c.On("Now").Return(timestamp).Once()
		// 2nd step: Refresh token.
		c.On("Now").Return(timestamp.Add(refreshTTL + time.Minute)).Once()
	})(t)

	p := newAPITokenProvider(cred, deviceID).
		WithBaseURL(s.URL()).
		WithTimeout(time.Second).
		WithMFAWait(time.Millisecond).
		WithRefreshTTL(refreshTTL).
		WithClock(c)

	// 1st try.
	token1, err := p.Token(context.Background())

	assert.Equal(t, s.AccessToken(), token1)
	assert.NotEmpty(t, string(token1))
	assert.NoError(t, err)

	// 2nd try.
	token2, err := p.Token(context.Background())

	assert.Equal(t, s.AccessToken(), token2)
	assert.NotEqual(t, token1, token2)
	assert.NotEmpty(t, string(token2))
	assert.NoError(t, err)
}
