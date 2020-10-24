package dstask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
