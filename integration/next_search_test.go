package integration

import (
	"testing"

	"github.com/naggie/dstask"
	"github.com/stretchr/testify/assert"
)

func TestNextSearchWord(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("somethingRandom")
	assertProgramResult(t, output, exiterr, success)

	var tasks []dstask.Task

	tasks = unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 0, "no tasks should be returned for a missing search term")

	output, exiterr, success = program("two")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, tasks[0].Summary, "two", "search term should find a task")
}
