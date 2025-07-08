package util

import "time"

// returns the current time in the UTC timezone
func NowUTC() time.Time {
	return time.Now().UTC()
}

// returns the "zero" time (0001-01-01 00:00). This is a placeholder for an unitialized
// time
func ZeroTime() time.Time {
	return time.Time{}
}
