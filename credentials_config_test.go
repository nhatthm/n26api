package n26api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api"
)

func TestCredentials(t *testing.T) {
	p := n26api.Credentials("username", "password")

	assert.Equal(t, "username", p.Username())
	assert.Equal(t, "password", p.Password())
}
