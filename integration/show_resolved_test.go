package integration

import (
	"testing"

	"github.com/naggie/dstask"
	"github.com/stretchr/testify/assert"
)

func TestShowResolved(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "+one", "one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "+two", "two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "three")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("1", "done")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("show-resolved")
	assertProgramResult(t, output, exiterr, success)

	var tasks []dstask.Task

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, tasks[0].Summary, "one", "one should be resolved")

	// Test the sorting of resolved tasks
	output, exiterr, success = program("3", "done")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("2", "done")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, tasks[0].Summary, "two", "two is the most-recently resolved")
}
