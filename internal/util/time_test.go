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
