package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowProjects(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "one", "project:myproject")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "two", "project:myproject")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "three", "project:myproject")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("show-projects")
	assertProgramResult(t, output, exiterr, success)

	projects := unmarshalProjectArray(t, output)
	assert.Equal(t, 0, projects[0].TasksResolved, "no tasks resolved")
	assert.Equal(t, projects[0].Tasks, 3, "three tasks created")

	output, exiterr, success = program("2", "done")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("show-projects")
	assertProgramResult(t, output, exiterr, success)

	projects = unmarshalProjectArray(t, output)
	assert.Equal(t, 1, projects[0].TasksResolved, "no tasks resolved")
	assert.Equal(t, 3, projects[0].Tasks, "three tasks created")
}
