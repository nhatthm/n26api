package n26api_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

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
		_ = os.Setenv(envUsername, currentUsername)
		_ = os.Setenv(envPassword, currentPassword)
	})

	_ = os.Setenv(envUsername, "username")
	_ = os.Setenv(envPassword, "password")

	p := n26api.CredentialsFromEnv()

	assert.Equal(t, "username", p.Username())
	assert.Equal(t, "password", p.Password())
}
