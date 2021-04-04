package testkit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api/pkg/testkit"
)

func TestCredentialsProvider(t *testing.T) {
	p := testkit.MockCredentialsProvider(func(p *testkit.CredentialsProvider) {
		p.On("Username").Return("username")
		p.On("Password").Return("password")
	})(t)

	expectedUsername := "username"
	expectedPassword := "password"

	assert.Equal(t, expectedUsername, p.Username())
	assert.Equal(t, expectedPassword, p.Password())
}
