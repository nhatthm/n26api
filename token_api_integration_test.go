// +build integration

package n26api

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationApiTokenProvider tests login functionalities.
//
// In order to run this tests, there env vars should be set:
// - N26_USERNAME: The username to login to N26, it should be an email address.
// - N26_PASSWORD: The password to login to n26api.
// - N26_DEVICE: The device ID in UUID format (optional).
func TestIntegrationApiTokenProvider(t *testing.T) {
	deviceID := deviceID(uuid.UUID{})
	apiUrl := os.Getenv("N26_BASE_URL")

	if apiUrl == "" {
		apiUrl = BaseUrl
	}

	p := newAPITokenProvider(apiUrl, time.Second, CredentialsFromEnv(), deviceID, liveClock{})

	done := make(chan struct{}, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go func() {
		token, err := p.Token(context.Background())

		assert.NotEmpty(t, token)
		assert.NoError(t, err)

		close(done)
	}()

	select {
	case <-done:
		return

	case <-ctx.Done():
		t.Fatal("timeout while getting access token")
	}
}
