package dstask

import (
	"strings"
	"time"
)

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func weekDayStrToTime(dateStr string, selector string) (due time.Time) {
	weekdays := map[string]time.Weekday{
		"sun":       time.Sunday,
		"sunday":    time.Sunday,
		"mon":       time.Monday,
		"monday":    time.Monday,
		"tue":       time.Tuesday,
		"tues":      time.Tuesday,
		"tuesday":   time.Tuesday,
		"wed":       time.Wednesday,
		"wednesday": time.Wednesday,
		"thu":       time.Thursday,
		"thur":      time.Thursday,
		"thurs":     time.Thursday,
		"thursday":  time.Thursday,
		"fri":       time.Friday,
		"friday":    time.Friday,
		"sat":       time.Saturday,
		"saturday":  time.Saturday,
	}
	nowWeekday := time.Now().Weekday()
	targetWeekday, ok := weekdays[strings.ToLower(dateStr)]
	if !ok {
		return time.Time{}
	}
	daysDifference := int(targetWeekday) - int(nowWeekday)

	if selector == "next" {
		return startOfDay(time.Now().AddDate(0, 0, daysDifference+7))
	}
	if selector == "this" || selector == "" {
		if daysDifference < 0 {
			return startOfDay(time.Now().AddDate(0, 0, daysDifference+7))
		}
	}
	return startOfDay(time.Now().AddDate(0, 0, daysDifference))
}

func ParseStrToDate(dateStr string) (due time.Time) {
	now := time.Now()
	lower := strings.ToLower(strings.TrimSpace(dateStr))

	switch lower {
	case "today":
		return startOfDay(now)
	case "tomorrow":
		return startOfDay(now.AddDate(0, 0, 1))
	case "yesterday":
		return startOfDay(now.AddDate(0, 0, -1))
	}

	// Check for next-[weekday], this-[weekday]
	parts := strings.SplitN(lower, "-", 2)
	if len(parts) == 2 {
		selector, rest := parts[0], parts[1]
		if wdTime := weekDayStrToTime(rest, selector); !wdTime.IsZero() {
			return wdTime
		}
	}

	// Check for [weekday]
	if wdTime := weekDayStrToTime(lower, ""); !wdTime.IsZero() {
		return wdTime
	}

	// Try YYYY-MM-DD, MM-DD, or DD
	if t, err := time.ParseInLocation("2006-01-02", dateStr, time.Local); err == nil {
		return t
	}
	if t, err := time.ParseInLocation("01-02", dateStr, time.Local); err == nil {
		t = t.AddDate(time.Now().Year(), 0, 0)
		return t
	}
	if t, err := time.ParseInLocation("2", dateStr, time.Local); err == nil {
		year, month, _ := time.Now().Date()
		t = t.AddDate(year, int(month)-1, 0)
		return t
	}

	ExitFail("Invalid due date format: " + dateStr + "\n" +
		"Expected format: YYYY-MM-DD, MM-DD or DD, relative date like 'next-monday', 'today', etc.")
	return time.Time{}
}
