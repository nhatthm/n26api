package n26api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api"
)

const (
	envUsername = "N26_USERNAME"
	envPassword = "N26_PASSWORD"
)

func TestCredentialsFromEnv(t *testing.T) {
	t.Setenv(envUsername, "username")
	t.Setenv(envPassword, "password")

	p := n26api.CredentialsFromEnv()

	assert.Equal(t, "username", p.Username())
	assert.Equal(t, "password", p.Password())
}
