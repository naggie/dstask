package dstask

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestModifyChangesProperties(t *testing.T) {
	type testCase struct {
		task     Task
		query    Query
		expected Task
	}

	testCases := []testCase{
		{ // Add a note
			Task{},
			Query{
				Template: 1,
				Note:     "Test Note",
			},
			Task{
				Notes: "Test Note",
			},
		},
		{ // Append a note
			Task{
				Notes: "Start Note",
			},
			Query{
				Note: "Query Note",
			},
			Task{
				Notes: "Start Note\nQuery Note",
			},
		},
		{ // Priority when not set
			Task{},
			Query{
				Priority: "P1",
			},
			Task{
				Priority: "P1",
			},
		},
		{ // Priority Overridden
			Task{
				Priority: "P3",
			},
			Query{
				Priority: "P1",
			},
			Task{
				Priority: "P1",
			},
		},
		{ // Removing projects
			Task{
				Project: "myproject",
			},
			Query{
				AntiProjects: []string{"myproject"},
			},
			Task{},
		},
	}

	for _, tc := range testCases {
		tc.task.Modify(tc.query)
		assert.Equal(t, tc.expected, tc.task)
	}
}

func TestMatchesFilterDue(t *testing.T) {
	day := func(day, hour int) time.Time {
		return time.Date(2025, 7, day, hour, 0, 0, 0, time.Local)
	}

	tests := []struct {
		due      time.Time
		query    Query
		expected bool
	}{
		{time.Time{}, Query{DateFilter: "none"}, true},
		{day(1, 0), Query{DateFilter: "none"}, false},
		{time.Time{}, Query{DateFilter: "any"}, false},
		{day(1, 0), Query{DateFilter: "any"}, true},
		// before and the end of a range include the whole day
		{day(31, 15), Query{DateFilter: "before", Due: day(31, 0)}, true},
		{day(32, 0), Query{DateFilter: "before", Due: day(31, 0)}, false},
		{day(1, 0), Query{DateFilter: "after", Due: day(1, 0), DueEnd: day(31, 0)}, true},
		{day(31, 15), Query{DateFilter: "after", Due: day(1, 0), DueEnd: day(31, 0)}, true},
		{day(32, 0), Query{DateFilter: "after", Due: day(1, 0), DueEnd: day(31, 0)}, false},
		{day(0, 15), Query{DateFilter: "after", Due: day(1, 0), DueEnd: day(31, 0)}, false},
		{time.Time{}, Query{DateFilter: "after", Due: day(1, 0), DueEnd: day(31, 0)}, false},
	}

	for i, tc := range tests {
		task := Task{Due: tc.due}
		assert.Equal(t, tc.expected, task.MatchesFilter(tc.query), "test %d", i)
	}
}

func TestMatchesFilterDueDST(t *testing.T) {
	berlin, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Skip("time zone Europe/Berlin not available")
	}

	// 2026-10-25 has 25 hours in Berlin, 2026-03-29 has 23 hours
	for _, day := range []int{25, 29} {
		month := time.October
		if day == 29 {
			month = time.March
		}

		end := time.Date(2026, month, day, 0, 0, 0, 0, berlin)
		query := Query{DateFilter: "after", Due: end.AddDate(0, 0, -6), DueEnd: end}

		lastMinute := Task{Due: time.Date(2026, month, day, 23, 59, 0, 0, berlin)}
		nextDay := Task{Due: time.Date(2026, month, day+1, 0, 0, 0, 0, berlin)}
		assert.True(t, lastMinute.MatchesFilter(query), end)
		assert.False(t, nextDay.MatchesFilter(query), end)
	}
}

func TestMatchesFilterWildcardTags(t *testing.T) {
	tests := []struct {
		tags     []string
		query    Query
		expected bool
	}{
		{[]string{"team:a"}, Query{Tags: []string{"team:*"}}, true},
		{[]string{"teamwork"}, Query{Tags: []string{"team:*"}}, false},
		{[]string{"teamwork"}, Query{AntiTags: []string{"team:*"}}, true},
		{[]string{"x"}, Query{Tags: []string{"*"}}, true},
		{nil, Query{Tags: []string{"*"}}, false},
		{nil, Query{AntiTags: []string{"*"}}, true},
		{[]string{"x"}, Query{AntiTags: []string{"*"}}, false},
	}

	for i, tc := range tests {
		task := Task{Tags: tc.tags}
		assert.Equal(t, tc.expected, task.MatchesFilter(tc.query), "test %d", i)
	}
}

func TestModifyDueNoneAndWildcardTags(t *testing.T) {
	task := Task{Due: time.Now(), Tags: []string{"a"}}
	task.Modify(Query{DateFilter: "none", Tags: []string{"b", "team:*"}})

	assert.True(t, task.Due.IsZero())
	assert.Equal(t, []string{"a", "b"}, task.Tags)
}

func TestMatchesFilterPriorities(t *testing.T) {
	tests := []struct {
		priority string
		query    Query
		expected bool
	}{
		{"P1", Query{}, true},
		{"P1", Query{Priority: "P1"}, true},
		{"P2", Query{Priority: "P1"}, false},
		{"P0", ParseQuery("P0", "P1"), true},
		{"P1", ParseQuery("P0", "P1"), true},
		{"P2", ParseQuery("P0", "P1"), false},
	}

	for i, tc := range tests {
		task := Task{Priority: tc.priority}
		assert.Equal(t, tc.expected, task.MatchesFilter(tc.query), "test %d", i)
	}
}
