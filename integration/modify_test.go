package integration

import (
	"testing"

	"github.com/naggie/dstask"
	"github.com/stretchr/testify/assert"
)

func TestModifyTasksByID(t *testing.T) {

	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "one", "+one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "two", "+two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "three", "+three")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("modify", "2", "3", "+extra")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	var tasks []dstask.Task

	tasks = unmarshalTaskArray(t, output)
	assert.ElementsMatch(t, []string{"three", "extra"}, tasks[0].Tags, "extra tag added to task three")
	assert.ElementsMatch(t, []string{"two", "extra"}, tasks[1].Tags, "extra tag added to task two")
	assert.ElementsMatch(t, []string{"one"}, tasks[2].Tags, "task 1 not modified")
}

func TestModifyTasksInContext(t *testing.T) {

	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "one", "+one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "two", "+two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "three", "+three")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("context", "+three")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("modify", "+extra")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	var tasks []dstask.Task

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, []string{"extra", "three"}, tasks[0].Tags, "tags should have been modified")
}
