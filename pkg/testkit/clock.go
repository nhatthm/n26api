package testkit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ClockMocker is Clock mocker.
type ClockMocker func(t testing.TB) *Clock

// NoMockClock is no mock Clock.
var NoMockClock = MockClock()

// Clock is a n26api.Clock.
type Clock struct {
	mock.Mock
}

// Now returns the current local time.
func (c *Clock) Now() time.Time {
	return c.Called().Get(0).(time.Time)
}

// mockClock mocks n26api.Clock interface.
func mockClock(mocks ...func(c *Clock)) *Clock {
	c := &Clock{}

	for _, m := range mocks {
		m(c)
	}

	return c
}

// MockClock creates Clock mock with cleanup to ensure all the expectations are met.
func MockClock(mocks ...func(c *Clock)) ClockMocker {
	return func(t testing.TB) *Clock {
		c := mockClock(mocks...)

		t.Cleanup(func() {
			assert.True(t, c.Mock.AssertExpectations(t))
		})

		return c
	}
}

// MockStaticClock returns a fixed timestamp on Now().
func MockStaticClock(timestamp time.Time) ClockMocker {
	return MockClock(func(c *Clock) {
		c.On("Now").Return(timestamp)
	})
}
