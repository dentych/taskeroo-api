package util

import "time"

func DurationToNextNoon(now time.Time) time.Duration {
	if now.Hour() < 12 {
		return time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location()).Sub(now)
	} else {
		tomorrow := now.AddDate(0, 0, 1)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 12, 0, 0, 0, tomorrow.Location()).Sub(now)
	}
}
