package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnixTimestampMS(t *testing.T) {
	input := time.Date(2020, 2, 2, 0, 0, 0, 1222222, time.UTC)
	expected := int64(1580601600001)

	assert.Equal(t, expected, UnixTimestampMS(input))
}
