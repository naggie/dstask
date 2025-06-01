package integration

import (
	"testing"

	"github.com/naggie/dstask"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, "one", tasks[0].Summary)

	output, exiterr, success = program("next", "+two")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, "two", tasks[0].Summary)
}

func TestNextMultipleTagFilter(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "+one", "one-alpha")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "+one", "+beta", "one-beta")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "+two", "two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "+one", "+beta")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Equal(t, "one-beta", tasks[0].Summary)
	assert.Len(t, tasks, 1)
}

func TestNextProjectFilter(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "+one", "one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "+two", "project:house", "two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "project:house")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Equal(t, "two", tasks[0].Summary)

	output, exiterr, success = program("project:house")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, "two", tasks[0].Summary)

	output, exiterr, success = program("-project:house")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, "one", tasks[0].Summary)
}
