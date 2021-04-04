package testkit_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api/pkg/testkit"
)

func TestMockStaticClock(t *testing.T) {
	timestamp := time.Now()
	clock := testkit.MockStaticClock(timestamp)(t)

	assert.Equal(t, timestamp, clock.Now())
}
