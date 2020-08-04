package dstask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidPartialUUID4String(t *testing.T) {

	type testCase struct {
		input          string
		expectedOutput bool
	}

	var tests = []testCase{
		{ // Valid full UUID4 string. Should return true.
			input:          "43217ff9-efa5-4fd4-b81d-3d8e40fad1c6",
			expectedOutput: true,
		},
		{ // Invalid UUID4 String. One extra character
			input:          "43217ff9-efa5-4fd4-b81d-3d8e40fad1c64",
			expectedOutput: false,
		},
		{ // Valid Partial String.
			input:          "43217ff9-efa5-4fd4-b81d-3d",
			expectedOutput: true,
		},
		{ // Invalid Partial String. Non-Hex Character
			input:          "43217ff9-efa5-4fd4-b81d-3dx",
			expectedOutput: false,
		},
		{ // Invalid Partial String. Incorrect format
			input:          "43217ff9-efa5-4fd4-b81d3d",
			expectedOutput: false,
		},
		{ // Invalid Partial String. Incorrect format.
			input:          "43217ff9-efa5-4fd4b81d-3dx",
			expectedOutput: false,
		},
		{ // Invalid Partial String. Incorrect format.
			input:          "43217ff9-efa54fd4-b81d-3dx",
			expectedOutput: false,
		},
		{ // Invalid Partial String. Incorrect format.
			input:          "43217ff9efa5-4fd4-b81d-3dx",
			expectedOutput: false,
		},
	}

	for _, tc := range tests {

		t.Logf("Testing %s in IsValidPartialUUID4String()", tc.input)
		output := IsValidPartialUUID4String(tc.input)
		assert.Equal(t, tc.expectedOutput, output)

	}

}
