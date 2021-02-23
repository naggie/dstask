package dstask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModifyChangesProperties(t *testing.T) {
	type testCase struct {
		task     Task
		query    Query
		expected Task
	}

	var testCases = []testCase{
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
