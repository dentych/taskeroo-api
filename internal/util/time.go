package util

import "time"

func DurationToNextNoon(now time.Time) time.Duration {
	hour := 12
	minute := 00
	if now.Hour() <= hour && now.Minute() < minute {
		return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location()).Sub(now)
	} else {
		tomorrow := now.AddDate(0, 0, 1)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), hour, minute, 0, 0, tomorrow.Location()).Sub(now)
	}
}
