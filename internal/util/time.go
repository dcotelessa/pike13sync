package util

import (
	"time"
)

// FormatDateTime formats a dateTime string for display
func FormatDateTime(dateTime string) string {
	t, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		return dateTime
	}
	return t.Format("Jan 2 at 3:04 PM")
}

// GetStartAndEndOfWeek calculates the start and end of the current week
func GetStartAndEndOfWeek() (string, string) {
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)
	
	return startOfWeek.Format(time.RFC3339), endOfWeek.Format(time.RFC3339)
}
