package dstask

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDuePeriod(t *testing.T) {
	// use a time zone with DST so that the DST cases below test something
	berlin, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Skip("time zone Europe/Berlin not available")
	}

	local := time.Local
	time.Local = berlin
	defer func() { time.Local = local }()

	date := func(year int, month time.Month, day int) time.Time {
		return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	}
	saturdayNight := date(2026, 11, 21).Add(23 * time.Hour)

	tests := []struct {
		now        time.Time
		period     string
		start, end time.Time
	}{
		// Monday, Sunday and Saturday
		{date(2026, 11, 16), "this-week", date(2026, 11, 16), date(2026, 11, 22)},
		{date(2026, 11, 22), "this-week", date(2026, 11, 16), date(2026, 11, 22)},
		{saturdayNight, "this-week", date(2026, 11, 16), date(2026, 11, 22)},
		{date(2026, 11, 16), "next-week", date(2026, 11, 23), date(2026, 11, 29)},
		{date(2026, 11, 22), "next-week", date(2026, 11, 23), date(2026, 11, 29)},
		// week across the turn of the year
		{date(2026, 12, 31), "this-week", date(2026, 12, 28), date(2027, 1, 3)},
		{date(2026, 12, 31), "next-week", date(2027, 1, 4), date(2027, 1, 10)},
		// months
		{date(2026, 11, 30), "this-month", date(2026, 11, 1), date(2026, 11, 30)},
		{date(2026, 12, 15), "this-month", date(2026, 12, 1), date(2026, 12, 31)},
		{date(2026, 12, 31), "next-month", date(2027, 1, 1), date(2027, 1, 31)},
		{date(2027, 1, 31), "next-month", date(2027, 2, 1), date(2027, 2, 28)},
		{date(2028, 1, 31), "next-month", date(2028, 2, 1), date(2028, 2, 29)},
		// DST changes in Europe
		{date(2026, 3, 29), "this-week", date(2026, 3, 23), date(2026, 3, 29)},
		{date(2026, 10, 25), "this-week", date(2026, 10, 19), date(2026, 10, 25)},
	}

	for _, tc := range tests {
		start, end, ok := parseDuePeriod("due:"+tc.period, tc.now)
		assert.True(t, ok)
		assert.Equal(t, tc.start, start, "%s at %s", tc.period, tc.now)
		assert.Equal(t, tc.end, end, "%s at %s", tc.period, tc.now)
	}

	for _, dueStr := range []string{"due:today", "due.before:this-week", "due:week"} {
		_, _, ok := parseDuePeriod(dueStr, date(2026, 11, 16))
		assert.False(t, ok, dueStr)
	}
}
