package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringPtr(t *testing.T) {
	s := "foobar"
	expected := "foobar"

	assert.Equal(t, &expected, StringPtr(s))
}

func TestInt64Ptr(t *testing.T) {
	i := int64(42)
	expected := int64(42)

	assert.Equal(t, &expected, Int64Ptr(i))
}
