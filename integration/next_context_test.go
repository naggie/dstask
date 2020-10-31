package integration

import (
	"os"
	"testing"

	"github.com/naggie/dstask"
	"github.com/stretchr/testify/assert"
)

func TestSettingTagContext(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "+one", "one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "+one", "one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "+two", "two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("context", "+two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	var tasks []dstask.Task

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, tasks[0].Summary, "two", "setting +two as a context")

	output, exiterr, success = program("context", "-one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, tasks[0].Summary, "two", "setting -one as a context")
}

func TestSettingTagAndProjectContext(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "+one", "+alpha", "one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "project:beta", "+two", "two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("context", "project:beta")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	var tasks []dstask.Task

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, tasks[0].Summary, "two", "setting project:beta as a context")

	output, exiterr, success = program("context", "project:beta", "+one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, 0, len(tasks), "no tasks within context project:beta +one")
}

func TestContextFromEnvVar(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "+one", "+alpha", "one")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "project:beta", "+two", "two")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("context", "project:beta")
	assertProgramResult(t, output, exiterr, success)

	// override context with an env var
	unsetEnv := setEnv("DSTASK_CONTEXT", "+one +alpha")
	t.Logf("DSTASK_CONTEXT=%s", os.Getenv("DSTASK_CONTEXT"))

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	var tasks []dstask.Task

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, tasks[0].Summary, "one", "'+one +alpha' context set by DSTASK_CONTEXT ")

	// unset the context override, so we expect to use the on-disk context
	unsetEnv()

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks = unmarshalTaskArray(t, output)
	assert.Equal(t, tasks[0].Summary, "two", "project:beta is on-disk context")
}
