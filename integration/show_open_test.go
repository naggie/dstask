package integration

import (
	"testing"

	"github.com/naggie/dstask"
	"gotest.tools/assert"
)

func TestShowOpen(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "+one", "one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "+two", "two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("show-open")
	assertProgramResult(t, output, exiterr, success)

	var tasks []dstask.Task

	// Oldest tasks come first
	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, "one", tasks[0].Summary, "one should be sorted first because it is older")
	assert.Equal(t, "two", tasks[1].Summary, "two should be sorted last")

	output, exiterr, success = program("context", "-one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("show-open")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, "two", tasks[0].Summary, "setting -one as a context")

	output, exiterr, success = program("2", "done")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("show-open")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, 0, len(tasks), "no tasks open in this context")

}
