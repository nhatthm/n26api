package n26api

import "time"

// Clock provides time-functionalities.
type Clock interface {
	// Now returns the current local time.
	Now() time.Time
}

type liveClock struct{}

// Now returns the current local time.
func (liveClock) Now() time.Time {
	return time.Now()
}
