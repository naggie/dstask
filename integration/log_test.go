package integration

import (
	"testing"

	"github.com/naggie/dstask"
	"gotest.tools/assert"
)

func TestLog(t *testing.T) {
	t.Skip("skipping because show-resolved cannot be run non-interactively yet")
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("log", "two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("show-resolved")
	assertProgramResult(t, output, exiterr, success)

	var tasks []dstask.Task

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, tasks[0].Summary, "two", "task two should be resolved")
}
