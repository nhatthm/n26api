package n26api

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	mockClock "github.com/nhatthm/go-clock/mock"
	"github.com/stretchr/testify/assert"

	authMock "github.com/nhatthm/n26api/pkg/testkit/auth"
)

func TestWithBaseURL(t *testing.T) {
	t.Parallel()

	expected := "http://example.com"
	c := NewClient(WithBaseURL(expected))

	assert.Equal(t, expected, c.config.baseURL)
}

func TestWithTimeout(t *testing.T) {
	t.Parallel()

	expected := time.Minute
	c := NewClient(WithTimeout(expected))

	assert.Equal(t, expected, c.config.timeout)
}

func TestWithUsername(t *testing.T) {
	t.Parallel()

	expected := "username"
	c := NewClient(WithUsername(expected))

	assert.Equal(t, expected, c.config.username)
}

func TestWithPassword(t *testing.T) {
	t.Parallel()

	expected := "password"
	c := NewClient(WithPassword(expected))

	assert.Equal(t, expected, c.config.password)
}

func TestWithDeviceID(t *testing.T) {
	t.Parallel()

	expected := uuid.New()
	c := NewClient(WithDeviceID(expected))

	assert.Equal(t, expected, c.config.deviceID)
}

func TestWithCredentials(t *testing.T) {
	t.Parallel()

	expectedUsername := "username"
	expectedPassword := "password"
	c := NewClient(WithUsername(expectedUsername), WithPassword(expectedPassword))

	assert.Equal(t, expectedUsername, c.config.username)
	assert.Equal(t, expectedPassword, c.config.password)
}

func TestWithCredentialsProvider(t *testing.T) {
	t.Parallel()

	provider1 := Credentials("provider1", "")
	provider2 := Credentials("provider2", "")
	c := NewClient(
		WithCredentialsProvider(provider2),
		WithCredentialsProvider(provider1),
	)

	expected := chainCredentialsProviders(
		Credentials("", ""),
		provider1,
		provider2,
		CredentialsFromEnv(),
	)

	assert.Equal(t, expected, c.config.credentials)
}

func TestWithCredentialsProviderAtLast(t *testing.T) {
	t.Parallel()

	provider1 := Credentials("provider1", "")
	provider2 := Credentials("provider2", "")
	c := NewClient(
		WithCredentialsProviderAtLast(provider2),
		WithCredentialsProviderAtLast(provider1),
	)

	expected := chainCredentialsProviders(
		Credentials("", ""),
		CredentialsFromEnv(),
		provider2,
		provider1,
	)

	assert.Equal(t, expected, c.config.credentials)
}

func TestWithTokenProvider(t *testing.T) {
	t.Parallel()

	provider1 := authMock.MockTokenProvider(func(p *authMock.TokenProvider) {
		// provider 1.
		p.On("Token", context.Background()).Return(1).Maybe()
	})(t)

	provider2 := authMock.MockTokenProvider(func(p *authMock.TokenProvider) {
		// provider 2.
		p.On("Token", context.Background()).Return(2).Maybe()
	})(t)

	c := NewClient(
		WithTokenProvider(provider1),
		WithTokenProvider(provider2),
	)

	expected := chainTokenProviders(provider2, provider1)

	assert.Equal(t, *expected, (*c.token)[0:len(*c.token)-1])
}

func TestWithTokenStorage(t *testing.T) {
	t.Parallel()

	expected := authMock.NoMockTokenStorage(t)
	c := NewClient(WithTokenStorage(expected))

	assert.Equal(t, expected, c.config.tokenStorage)
}

func TestWithClock(t *testing.T) {
	t.Parallel()

	clock := mockClock.NoMock(t)
	c := NewClient(WithClock(clock))

	assert.Equal(t, clock, c.clock)
}

func TestWithMFATimeout(t *testing.T) {
	t.Parallel()

	expected := time.Second
	c := NewClient(WithMFATimeout(expected))

	assert.Equal(t, expected, c.config.mfaTimeout)
}

func TestWithMFAWait(t *testing.T) {
	t.Parallel()

	expected := time.Second
	c := NewClient(WithMFAWait(expected))

	assert.Equal(t, expected, c.config.mfaWait)
}
