package server

import (
	"fmt"
	"time"
)

func formatUptime(seconds int64) string {
	d := time.Duration(seconds) * time.Second

	if seconds < 60 {
		return fmt.Sprintf("%d seconds", seconds)
	}

	minutes := int(d.Minutes())
	hours := minutes / 60
	days := hours / 24

	switch {
	case days > 0 && hours%24 > 0:
		return fmt.Sprintf("%d days %d hours", days, hours%24)
	case days > 0:
		return fmt.Sprintf("%d days", days)
	case hours > 0 && minutes%60 > 0:
		return fmt.Sprintf("%d hours %d minutes", hours, minutes%60)
	case hours > 0:
		return fmt.Sprintf("%d hours", hours)
	default:
		return fmt.Sprintf("%d minutes", minutes)
	}
}

func formatRelative(t time.Time) string {
	diff := time.Since(t)

	switch {
	case diff < 2*time.Second:
		return "just now"
	case diff < time.Minute:
		return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
	case diff < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	case diff < 7*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	case diff < 30*24*time.Hour:
		return fmt.Sprintf("%d weeks ago", int(diff.Hours()/(24*7)))
	default:
		// Fallback to calendar format
		return t.Format("Jan 2, 2006 at 15:04")
	}
}
