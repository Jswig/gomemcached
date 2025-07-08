package util

import (
	"testing"
	"time"
)

func TestNowUTC(t *testing.T) {
	result := NowUTC()
	if result.Location() != time.UTC {
		t.Errorf("got location %v, want UTC", result.Location())
	}
}

func TestZeroTime(t *testing.T) {
	result := ZeroTime()
	if !result.IsZero() {
		t.Errorf("got time %v, expected zero time", result)
	}
}
