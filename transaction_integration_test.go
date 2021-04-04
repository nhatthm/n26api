// +build integration

package n26api

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegrationFindAllTransactionsInRange tests transactions functionalities.
//
// In order to run this tests, there env vars should be set:
// - N26_USERNAME: The username to login to N26, it should be an email address.
// - N26_PASSWORD: The password to login to n26api.
// - N26_DEVICE: The device ID in UUID format (optional).
// - N26_FROM: Transaction starting time. Default to 1 day ago (optional).
// - N26_TO: Transaction ending time. Default to now (optional).
func TestIntegrationFindAllTransactionsInRange(t *testing.T) {
	var from time.Time
	var to time.Time
	var err error

	if v, ok := os.LookupEnv("N26_TO"); ok {
		to, err = time.Parse(time.RFC3339, v)
		require.NoError(t, err)
	} else {
		to = time.Now()
	}

	if v, ok := os.LookupEnv("N26_FROM"); ok {
		from, err = time.Parse(time.RFC3339, v)
		require.NoError(t, err)
	} else {
		from = to.AddDate(0, 0, -1)
	}

	deviceID := deviceID(uuid.UUID{})
	baseURL := os.Getenv("N26_BASE_URL")

	if baseURL == "" {
		baseURL = BaseURL
	}

	c := NewClient(
		WithBaseURL(baseURL),
		WithDeviceID(deviceID),
	)

	done := make(chan struct{}, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go func() {
		transactions, err := c.FindAllTransactionsInRange(ctx, from, to)

		assert.NotEmpty(t, transactions)
		assert.NoError(t, err)

		close(done)
	}()

	select {
	case <-done:
		return

	case <-ctx.Done():
		t.Fatal("timeout while getting transactions")
	}
}
