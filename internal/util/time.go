package util

import "time"

// returns the current time in the UTC timezone
func NowUTC() time.Time {
	return time.Now().UTC()
}
