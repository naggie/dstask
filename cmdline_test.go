package dstask

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	type testCase struct {
		input    []string
		expected Query
	}

	var tests = []testCase{
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
		// first priority should have precedence, subsequent priorities should
		// just be part of the description
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
				Template:      0,
				Text:          "P2 P3",
				IgnoreContext: false,
				Note:          "",
			},
		},
		// same for projects, for consistency
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
