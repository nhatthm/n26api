package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api/pkg/util"
)

func TestBase64Credentials(t *testing.T) {
	username := "foo"
	password := "bar"
	expected := "Zm9vOmJhcg=="

	assert.Equal(t, expected, util.Base64Credentials(username, password))
}
