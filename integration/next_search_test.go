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

	output, exiterr, success := program("add", "one", "/", "alpha")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "two")
	assertProgramResult(t, output, exiterr, success)

	// search something that doesn't exist
	output, exiterr, success = program("somethingRandom")
	assertProgramResult(t, output, exiterr, success)

	var tasks []dstask.Task

	tasks = unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 0, "no tasks should be returned for a missing search term")

	// search the summary of task two
	output, exiterr, success = program("two")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, "two", tasks[0].Summary, "search term should find a task")

	// search the notes field of task one
	output, exiterr, success = program("alpha")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, "one", tasks[0].Summary, "string \"alpha\" is in a note for task one")
}
