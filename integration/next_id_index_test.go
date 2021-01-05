package integration

import (
	"testing"

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

	tasks := unmarshalTaskArray(t, output)
	assert.Equal(t, "one", tasks[0].Summary, "find task 1 by ID")
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

	tasks := unmarshalTaskArray(t, output)
	assert.Equal(t, "two", tasks[0].Summary, "find task 2 by ID (context ignored with ID based addressing)")

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, 1, tasks[0].ID, "1 is the only ID in our current context")
}

func TestNextWithurgency(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "two", "+two", "project:two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "three", "P0")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "four", "project:four")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)
	tasks := unmarshalTaskArray(t, output)
	assert.Equal(t, "three", tasks[0].Summary)
	assert.Equal(t, "two", tasks[1].Summary)
	assert.Equal(t, "four", tasks[2].Summary)
	assert.Equal(t, "one", tasks[3].Summary)
}
