package dstask

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeTempFilename(t *testing.T) {
	type testCase struct {
		id       int
		summary  string
		expected string
	}
	var testCases = []testCase{
		{
			1,
			`& &`,
			`dstask.*.1-.md`,
		},
		{
			2147483647, // max int32
			`J's $100, != €100`,
			`dstask.*.2147483647-js-100-100.md`,
		},
		{
			-2147483648, // min int32
			`J's $100, != €100`,
			`dstask.*.-2147483648-js-100-100.md`,
		},
		{
			99,
			`A simple summary!`,
			`dstask.*.99-a-simple-summary.md`,
		},
		{
			1,
			`& that's that.`,
			`dstask.*.1-thats-that.md`,
		},
	}

	for _, tc := range testCases {
		tf := MakeTempFilename(tc.id, tc.summary, "md")

		assert.Equal(t, tc.expected, tf)

		f, err := os.CreateTemp("", tf)
		assert.Nil(t, err)
		assert.Nil(t, f.Close())
		assert.Nil(t, os.Remove(f.Name()))
	}
}

func TestStrSliceContainsAll(t *testing.T) {

	type testCase struct {
		subset   []string
		superset []string
		expected bool
	}

	var testCases = []testCase{
		{
			[]string{},
			[]string{},
			true,
		},
		{
			[]string{"one"},
			[]string{"one"},
			true,
		},
		{
			[]string{"one"},
			[]string{"two"},
			false,
		},
		{
			[]string{"one"},
			[]string{},
			false,
		},
		{
			[]string{"one"},
			[]string{"one", "two"},
			true,
		},
		{
			[]string{"one", "two"},
			[]string{"one", "two"},
			true,
		},
		{
			[]string{"two", "one"},
			[]string{"three", "one", "two"},
			true,
		},
		{
			[]string{"apple", "two", "one"},
			[]string{"three", "one", "two"},
			false,
		},
		{
			[]string{},
			[]string{"three", "one", "two"},
			true,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, StrSliceContainsAll(tc.subset, tc.superset))
	}

}
