package n26api

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api/pkg/auth"
	"github.com/nhatthm/n26api/pkg/testkit"
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

	provider := testkit.NoMockCredentialsProvider(t)
	c := NewClient(WithCredentialsProvider(provider))

	expected := []CredentialsProvider{provider}

	assert.Equal(t, expected, c.config.credentials)
}

func TestWithTokenProvider(t *testing.T) {
	t.Parallel()

	provider := authMock.NoMockTokenProvider(t)
	c := NewClient(WithTokenProvider(provider))

	expected := []auth.TokenProvider{provider}

	assert.Equal(t, expected, c.config.tokens)
}

func TestWithTokenStorage(t *testing.T) {
	t.Parallel()

	expected := authMock.NoMockTokenStorage(t)
	c := NewClient(WithTokenStorage(expected))

	assert.Equal(t, expected, c.config.tokenStorage)
}

func TestWithClock(t *testing.T) {
	t.Parallel()

	clock := testkit.NoMockClock(t)
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
