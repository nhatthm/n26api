package util

import "time"

// UnixTimestampMS returns unix timestamp in milliseconds.
func UnixTimestampMS(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
