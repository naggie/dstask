package integration

import (
	"testing"

	"github.com/naggie/dstask"
	"gotest.tools/assert"
)

func TestNextTagFilter(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "+one", "one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "+one", "one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "+two", "two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "+one")
	assertProgramResult(t, output, exiterr, success)

	var tasks []dstask.Task

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, tasks[0].Summary, "one")

	output, exiterr, success = program("next", "+two")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, tasks[0].Summary, "two")

}
