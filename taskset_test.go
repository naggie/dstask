package dstask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterStringSlice(t *testing.T) {
	type testCase struct {
		with, without []string
		expected      []string
	}

	var testCases = []testCase{
		{
			with:     []string{"apple", "banana", "orange"},
			without:  []string{"orange"},
			expected: []string{"apple", "banana"},
		},
		{
			with:     []string{"apple", "banana"},
			without:  []string{"orange"},
			expected: []string{"apple", "banana"},
		},
		{
			with:     []string{"apple", "banana", ""},
			without:  []string{"orange"},
			expected: []string{"apple", "banana", ""},
		},
		{
			with:     []string{"apple", "banana", "orange"},
			without:  []string{"banana"},
			expected: []string{"apple", "orange"},
		},
		{
			with:     []string{"apple", "banana", "orange"},
			without:  nil,
			expected: []string{"apple", "banana", "orange"},
		},
		{
			with:     nil,
			without:  nil,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		actual := filterStringSlice(tc.with, tc.without)
		if !assert.Equal(t, tc.expected, actual) {
			t.Fail()
		}
	}
}
