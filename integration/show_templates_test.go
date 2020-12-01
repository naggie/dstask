package integration

import (
	"testing"

	"gotest.tools/assert"
)

// TODO

func TestTaskShowTemplates(t *testing.T) {

	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("template", "template1")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("show-templates")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Equal(t, "template1", tasks[0].Summary, "should be a template")
}
