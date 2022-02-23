package util

import (
	"testing"
	"time"
)

func TestDurationTillNoon(t *testing.T) {
	duration := DurationToNextNoon(time.Date(2022, 02, 23, 10, 0, 00, 00, time.Local))
	expected := 2.0
	if duration.Hours() != expected {
		t.Errorf("Expected %f but got: %f\n", expected, duration.Hours())
	}

	duration = DurationToNextNoon(time.Date(2022, 02, 23, 22, 0, 00, 00, time.Local))
	expected = 14.0
	if duration.Hours() != expected {
		t.Errorf("Expected %f but got: %f\n", expected, duration.Hours())
	}
}
