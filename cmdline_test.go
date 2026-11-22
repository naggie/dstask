package dstask

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	type testCase struct {
		input    []string
		expected Query
	}

	tests := []testCase{
		{
			[]string{"add", "have", "an", "adventure"},
			Query{
				Cmd:           "add",
				IDs:           nil,
				Tags:          nil,
				AntiTags:      nil,
				Project:       "",
				AntiProjects:  nil,
				Template:      0,
				Text:          "have an adventure",
				IgnoreContext: false,
				Note:          "",
			},
		},
		{
			[]string{"add", "+x", "-y", "have", "an", "adventure"},
			Query{
				Cmd:           "add",
				IDs:           nil,
				Tags:          []string{"x"},
				AntiTags:      []string{"y"},
				Project:       "",
				AntiProjects:  nil,
				Template:      0,
				Text:          "have an adventure",
				IgnoreContext: false,
				Note:          "",
			},
		},
		{
			[]string{"add", "smile", "/"},
			Query{
				Cmd:           "add",
				IDs:           nil,
				Tags:          nil,
				AntiTags:      nil,
				Project:       "",
				AntiProjects:  nil,
				Template:      0,
				Text:          "smile",
				IgnoreContext: false,
				Note:          "",
			},
		},

		{
			[]string{"add", "floss", "project:p", "+health", "/", "every  day"},
			Query{
				Cmd:           "add",
				IDs:           nil,
				Tags:          []string{"health"},
				AntiTags:      nil,
				Project:       "p",
				AntiProjects:  nil,
				Template:      0,
				Text:          "floss",
				IgnoreContext: false,
				Note:          "every  day",
			},
		},
		{
			[]string{"16", "modify", "+project:p", "-project:x", "-fun"},
			Query{
				Cmd:           "modify",
				IDs:           []int{16},
				Tags:          nil,
				AntiTags:      []string{"fun"},
				Project:       "p",
				AntiProjects:  []string{"x"},
				Template:      0,
				Text:          "",
				IgnoreContext: false,
				Note:          "",
			},
		},
		{
			[]string{"--", "show-resolved"},
			Query{
				Cmd:           "show-resolved",
				IDs:           nil,
				Tags:          nil,
				AntiTags:      nil,
				Project:       "",
				AntiProjects:  nil,
				Template:      0,
				Text:          "",
				IgnoreContext: true,
				Note:          "",
			},
		},
		// several priorities are a filter (any of them), the first one is kept
		// in Priority. Commands that set a priority reject them.
		// see https://github.com/naggie/dstask/issues/120 for context
		{
			[]string{"add", "P1", "P2", "P3"},
			Query{
				Cmd:           "add",
				IDs:           nil,
				Tags:          nil,
				AntiTags:      nil,
				Project:       "",
				AntiProjects:  nil,
				Priority:      PRIORITY_HIGH,
				Priorities:    []string{PRIORITY_HIGH, PRIORITY_NORMAL, PRIORITY_LOW},
				Template:      0,
				Text:          "",
				IgnoreContext: false,
				Note:          "",
			},
		},
		// the first project has precedence, a second one is part of the description
		{
			[]string{"add", "project:foo", "project:bar"},
			Query{
				Cmd:           "add",
				IDs:           nil,
				Tags:          nil,
				AntiTags:      nil,
				Project:       "foo",
				AntiProjects:  nil,
				Template:      0,
				Text:          "project:bar",
				IgnoreContext: false,
				Note:          "",
			},
		},
		{
			[]string{"add", "My", "Task", "template:1", "/", "Test", "Note"},
			Query{
				Cmd:           "add",
				IDs:           nil,
				Tags:          nil,
				AntiTags:      nil,
				Project:       "",
				AntiProjects:  nil,
				Template:      1,
				Text:          "My Task",
				IgnoreContext: false,
				Note:          "Test Note",
			},
		},
		{
			[]string{"next", "due:none"},
			Query{Cmd: "next", DateFilter: "none"},
		},
		{
			[]string{"next", "due:any"},
			Query{Cmd: "next", DateFilter: "any"},
		},
		{
			[]string{"next", "due.after:2025-07-01", "due.before:2025-07-31"},
			Query{
				Cmd:        "next",
				DateFilter: "after",
				Due:        time.Date(2025, 7, 1, 0, 0, 0, 0, time.Local),
				DueEnd:     time.Date(2025, 7, 31, 0, 0, 0, 0, time.Local),
			},
		},
		{
			[]string{"next", "due.before:2025-07-31", "due.after:2025-07-01"},
			Query{
				Cmd:        "next",
				DateFilter: "after",
				Due:        time.Date(2025, 7, 1, 0, 0, 0, 0, time.Local),
				DueEnd:     time.Date(2025, 7, 31, 0, 0, 0, 0, time.Local),
			},
		},
		{
			[]string{"next", "+team:*", "-name:*"},
			Query{Cmd: "next", Tags: []string{"team:*"}, AntiTags: []string{"name:*"}},
		},
		{
			[]string{"next", "DUE:NONE", "+Team:*"},
			Query{Cmd: "next", DateFilter: "none", Tags: []string{"team:*"}},
		},
		{
			[]string{"next", "due.after:2025-07-01", "due:overdue"},
			Query{
				Cmd:        "next",
				DateFilter: "after",
				Due:        time.Date(2025, 7, 1, 0, 0, 0, 0, time.Local),
				DueEnd:     startOfDay(time.Now()),
			},
		},
	} // end test cases

	for i, tc := range tests {
		description := strings.Join(tc.input, " ")

		t.Run(fmt.Sprintf("test %v: %s", i, description), func(t *testing.T) {
			t.Parallel()

			actual := ParseQuery(tc.input...)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryStringRoundtrip(t *testing.T) {
	tests := [][]string{
		{"due:none"},
		{"due:any"},
		{"due.after:2025-07-01", "due.before:2025-07-31"},
		{"due.before:2025-07-31"},
		{"+team:*", "-*"},
		{"P0", "P1"},
		{"P3", "P1", "+work"},
	}

	for _, args := range tests {
		query := ParseQuery(args...)
		assert.Equal(t, query, ParseQuery(strings.Fields(query.String())...), args)
	}
}

func TestMergeDue(t *testing.T) {
	tests := []struct {
		query    []string
		ctx      []string
		expected []string
	}{
		{[]string{}, []string{"due:none"}, []string{"due:none"}},
		{[]string{"due:none"}, []string{"due:none"}, []string{"due:none"}},
		{[]string{"due:2025-07-02"}, []string{"due:any"}, []string{"due:2025-07-02"}},
		{[]string{"due:any"}, []string{"due.before:2025-07-31"}, []string{"due.before:2025-07-31"}},
		{
			[]string{"due.before:2025-07-31"},
			[]string{"due.after:2025-07-01"},
			[]string{"due.after:2025-07-01", "due.before:2025-07-31"},
		},
		{
			[]string{"due.after:2025-07-01", "due.before:2025-07-31"},
			[]string{"due.after:2025-07-01", "due.before:2025-07-31"},
			[]string{"due.after:2025-07-01", "due.before:2025-07-31"},
		},
		{[]string{"due.on:2025-07-02"}, []string{"due:2025-07-02"}, []string{"due.on:2025-07-02"}},
		{[]string{"due:2025-07-02"}, []string{"due.in:2025-07-02"}, []string{"due:2025-07-02"}},
	}

	for _, tc := range tests {
		query := ParseQuery(tc.query...)
		assert.Equal(t, ParseQuery(tc.expected...), query.Merge(ParseQuery(tc.ctx...)), tc)
	}
}

func TestWithoutDueFilter(t *testing.T) {
	ctx := ParseQuery("+work", "due.after:2025-07-01", "due.before:2025-07-31")
	assert.Equal(t, ParseQuery("+work"), withoutDueFilter(ctx))

	assert.Equal(t, ParseQuery("+work"), withoutDueFilter(ParseQuery("+work", "due:none")))
	assert.Equal(t, ParseQuery("+work"), withoutDueFilter(ParseQuery("+work", "due:any")))

	ctx = ParseQuery("+work", "due.after:2025-07-01")
	assert.Equal(t, ctx, withoutDueFilter(ctx))
}

func TestParseQueryPriorities(t *testing.T) {
	query := ParseQuery("next", "P1", "foo", "P0")
	assert.Equal(t, "P1", query.Priority)
	assert.Equal(t, []string{"P1", "P0"}, query.Priorities)
	assert.Equal(t, "foo", query.Text)

	// a repeated priority is the same as a single one
	assert.Equal(t, ParseQuery("P1"), ParseQuery("P1", "P1"))
	assert.Nil(t, ParseQuery("P1", "P1").Priorities)

	assert.Equal(t, ParseQuery("P1", "P0", "P1"), ParseQuery("P1", "P0"))
}

func TestMergePriority(t *testing.T) {
	tests := []struct {
		query    []string
		ctx      []string
		expected []string
	}{
		{[]string{}, []string{"P1"}, []string{"P1"}},
		{[]string{"P1"}, []string{}, []string{"P1"}},
		{[]string{}, []string{"P0", "P1"}, []string{"P0", "P1"}},
		{[]string{"P0", "P1"}, []string{}, []string{"P0", "P1"}},
		{[]string{"P1"}, []string{"P0", "P1"}, []string{"P1"}},
		{[]string{"P1", "P0"}, []string{"P0", "P1"}, []string{"P1", "P0"}},
		{[]string{"P0", "P1"}, []string{"P0", "P1", "P2"}, []string{"P0", "P1"}},
	}

	for _, tc := range tests {
		query := ParseQuery(tc.query...)
		assert.Equal(t, ParseQuery(tc.expected...), query.Merge(ParseQuery(tc.ctx...)), tc)
	}

	// a context stored by an older version has no Priorities
	query := ParseQuery()
	assert.Equal(t, ParseQuery("P2"), query.Merge(Query{Priority: "P2"}))
}

func TestWithoutPriorityFilter(t *testing.T) {
	assert.Equal(t, ParseQuery("+work"), withoutPriorityFilter(ParseQuery("+work", "P0", "P1")))

	ctx := ParseQuery("+work", "P1")
	assert.Equal(t, ctx, withoutPriorityFilter(ctx))
}
