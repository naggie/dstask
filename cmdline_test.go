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
