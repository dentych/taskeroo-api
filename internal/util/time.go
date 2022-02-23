package util

import "time"

func DurationToNextNoon(now time.Time) time.Duration {
	hour := 21
	if now.Hour() < hour {
		return time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location()).Sub(now)
	} else {
		tomorrow := now.AddDate(0, 0, 1)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), hour, 0, 0, 0, tomorrow.Location()).Sub(now)
	}
}
