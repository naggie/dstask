package integration

import (
	"testing"

	"github.com/naggie/dstask"
	"gotest.tools/assert"
)

func TestNextByIDIndex(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "two", "+two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("1")
	assertProgramResult(t, output, exiterr, success)

	var tasks []dstask.Task

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, tasks[0].Summary, "one", "find task 1 by ID")
}

func TestNextByIDIndexOutsideContext(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "one", "+one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "two", "+two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("context", "+one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("2")
	assertProgramResult(t, output, exiterr, success)

	var tasks []dstask.Task

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, tasks[0].Summary, "two", "find task 2 by ID (context ignored with ID based addressing)")

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, tasks[0].ID, 1, "1 is the only ID in our current context")
}
