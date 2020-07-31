package dstask

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNote(t *testing.T) {

	type testCase struct {
		cmdLine       CmdLine
		task          Task
		openEditor    bool
		errorExpected bool
		notes         string
	}

	var tests = []testCase{

		{
			CmdLine{
				Cmd:  "note",
				Text: "I should make a note",
			},
			Task{
				Notes: "",
			},
			true,
			false,
			"I should make a note",
		},
		{
			CmdLine{
				Cmd:  "note",
				Text: "I should make another note",
			},
			Task{
				Notes: "There is a note here already.",
			},
			true,
			false,
			"There is a note here already.\nI should make another note",
		},
		{
			CmdLine{
				Cmd:  "note",
				Text: "I should make another note",
			},
			Task{
				Notes: "openEditor is false. This note should not change.",
			},
			false,
			false,
			"openEditor is false. This note should not change.",
		},
	}

	for i, tc := range tests {

		description := fmt.Sprintf("Running cmd: %s with Text: %s on Task: %v", tc.cmdLine.Cmd, tc.cmdLine.Text, tc.task)

		t.Run(fmt.Sprintf("test %v: %s", i, description), func(t *testing.T) {
			//t.Parallel() - This causes the test to fail by adding appending tc.cmdLine.Text twice.

			err := tc.task.AddNote(tc.cmdLine, tc.openEditor)
			notes := tc.task.Notes

			assert := assert.New(t)
			assert.Equal(tc.notes, notes)
			if tc.errorExpected {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
			}
		})

	}
}
