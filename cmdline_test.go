package dstask

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCmdLine(t *testing.T) {

	type testCase struct {
		input    []string
		expected CmdLine
	}

	var tests = []testCase{
		{
			[]string{"add", "have", "an", "adventure"},
			CmdLine{
				Cmd:           "add",
				IDs:           nil,
				Tags:          nil,
				AntiTags:      nil,
				Project:       "",
				AntiProjects:  nil,
				Template:      0,
				Text:          "have an adventure",
				UUID:          "",
				IgnoreContext: false,
				IDsExhausted:  true,
				Note:          "",
			},
		},
		{
			[]string{"add", "+x", "-y", "have", "an", "adventure"},
			CmdLine{
				Cmd:           "add",
				IDs:           nil,
				Tags:          []string{"x"},
				AntiTags:      []string{"y"},
				Project:       "",
				AntiProjects:  nil,
				Template:      0,
				Text:          "have an adventure",
				UUID:          "",
				IgnoreContext: false,
				IDsExhausted:  true,
				Note:          "",
			},
		},
		{
			[]string{"add", "smile", "/"},
			CmdLine{
				Cmd:           "add",
				IDs:           nil,
				Tags:          nil,
				AntiTags:      nil,
				Project:       "",
				AntiProjects:  nil,
				Template:      0,
				Text:          "smile",
				UUID:          "",
				IgnoreContext: false,
				IDsExhausted:  true,
				Note:          "",
			},
		},

		{
			[]string{"add", "floss", "project:p", "+health", "/", "every  day"},
			CmdLine{
				Cmd:           "add",
				IDs:           nil,
				Tags:          []string{"health"},
				AntiTags:      nil,
				Project:       "p",
				AntiProjects:  nil,
				Template:      0,
				Text:          "floss",
				UUID:          "",
				IgnoreContext: false,
				IDsExhausted:  true,
				Note:          "every  day",
			},
		},
		{
			[]string{"16", "modify", "+project:p", "-project:x", "-fun"},
			CmdLine{
				Cmd:           "modify",
				IDs:           []int{16},
				Tags:          nil,
				AntiTags:      []string{"fun"},
				Project:       "p",
				AntiProjects:  []string{"x"},
				Template:      0,
				Text:          "",
				UUID:          "",
				IgnoreContext: false,
				IDsExhausted:  true,
				Note:          "",
			},
		},

		{
			[]string{"modify", "uuid:0023e003-3453-9348-3452-898729283abc"},
			CmdLine{
				Cmd:           "modify",
				IDs:           nil,
				Tags:          nil,
				AntiTags:      nil,
				Project:       "",
				AntiProjects:  nil,
				Template:      0,
				Text:          "",
				UUID:          "0023e003-3453-9348-3452-898729283abc",
				IgnoreContext: false,
				IDsExhausted:  true,
				Note:          "",
			},
		},
	} // end test cases

	for i, tc := range tests {

		description := strings.Join(tc.input, " ")

		t.Run(fmt.Sprintf("test %v: %s", i, description), func(t *testing.T) {
			t.Parallel()

			parsed := ParseCmdLine(tc.input...)
			assert.Equal(t, parsed, tc.expected)

		})
	}
}

func TestMustGetIdentifiers(t *testing.T) {

	type Result struct {
		idents  []interface{}
		taskSet []string
		err     error
	}

	type testCase struct {
		input    CmdLine
		expected Result
	}

	var tests = []testCase{
		{
			CmdLine{
				Cmd:           "",
				IDs:           []int{16, 43, 11},
				Tags:          nil,
				AntiTags:      nil,
				Project:       "",
				AntiProjects:  nil,
				Template:      0,
				Text:          "",
				UUID:          "",
				IgnoreContext: false,
				IDsExhausted:  true,
				Note:          "",
			},
			Result{
				idents:  []interface{}{16, 43, 11},
				taskSet: NON_RESOLVED_STATUSES,
				err:     nil,
			},
		},
		{
			CmdLine{
				Cmd:           "",
				IDs:           []int{16, 43, 11},
				Tags:          nil,
				AntiTags:      nil,
				Project:       "",
				AntiProjects:  nil,
				Template:      0,
				Text:          "",
				UUID:          "0023e003-3453-9348-3452-898729283abc",
				IgnoreContext: false,
				IDsExhausted:  true,
				Note:          "",
			},
			Result{
				idents:  []interface{}{16, 43, 11},
				taskSet: NON_RESOLVED_STATUSES,
				err:     nil,
			},
		},
		{
			CmdLine{
				Cmd:           "",
				IDs:           nil,
				Tags:          nil,
				AntiTags:      nil,
				Project:       "",
				AntiProjects:  nil,
				Template:      0,
				Text:          "",
				UUID:          "0023e003-3453-9348-3452-898729283abc",
				IgnoreContext: false,
				IDsExhausted:  true,
				Note:          "",
			},
			Result{
				idents:  []interface{}{"0023e003-3453-9348-3452-898729283abc"},
				taskSet: []string{STATUS_RESOLVED},
				err:     nil,
			},
		},
		{
			CmdLine{
				Cmd:           "",
				IDs:           nil,
				Tags:          nil,
				AntiTags:      nil,
				Project:       "",
				AntiProjects:  nil,
				Template:      0,
				Text:          "",
				UUID:          "",
				IgnoreContext: false,
				IDsExhausted:  true,
				Note:          "",
			},
			Result{
				idents:  nil,
				taskSet: nil,
				err:     errors.New("MustGetIdentifiers() did not find any UUIDs or IDs in the command line."),
			},
		},
	}
	for i, tc := range tests {

		t.Run(fmt.Sprintf("test %v: %s", i, tc.input), func(t *testing.T) {
			t.Parallel()

			idents, taskSet, err := tc.input.MustGetIdentifiers()
			assert.Equal(t, tc.expected.idents, idents)
			assert.Equal(t, tc.expected.taskSet, taskSet)
			assert.Equal(t, tc.expected.err, err)

		})
	}
}
