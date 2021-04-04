package n26api_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nhatthm/n26api"
)

const (
	envUsername = "N26_USERNAME"
	envPassword = "N26_PASSWORD"
)

func TestCredentialsFromEnv(t *testing.T) {
	currentUsername := os.Getenv(envUsername)
	currentPassword := os.Getenv(envPassword)

	t.Cleanup(func() {
		err := os.Setenv(envUsername, currentUsername)
		require.NoError(t, err)

		err = os.Setenv(envPassword, currentPassword)
		require.NoError(t, err)
	})

	err := os.Setenv(envUsername, "username")
	require.NoError(t, err)

	err = os.Setenv(envPassword, "password")
	require.NoError(t, err)

	p := n26api.CredentialsFromEnv()

	assert.Equal(t, "username", p.Username())
	assert.Equal(t, "password", p.Password())
}
