package integration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeRepoWithPriorities(t *testing.T) (func(args ...string) []byte, func()) {
	t.Helper()

	repo, cleanup := makeDstaskRepo(t)
	program := testCmd(repo)

	run := func(args ...string) []byte {
		t.Helper()

		output, exiterr, success := program(args...)
		assertProgramResult(t, output, exiterr, success)

		return output
	}

	for _, priority := range []string{"P0", "P1", "P2", "P3"} {
		run("add", "task-"+priority, priority)
	}

	return run, cleanup
}

func TestNextMultiplePriorities(t *testing.T) {
	run, cleanup := makeRepoWithPriorities(t)
	defer cleanup()

	output := run("next", "P1", "P0")
	assert.Equal(t, []string{"task-P0", "task-P1"}, taskSummaries(unmarshalTaskArray(t, output)))

	output = run("P3", "P1")
	assert.Equal(t, []string{"task-P1", "task-P3"}, taskSummaries(unmarshalTaskArray(t, output)))
}

func TestMultiplePrioritiesContext(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	for _, priority := range []string{"P0", "P1", "P2"} {
		output, exiterr, success := program("add", "task-"+priority, priority)
		assertProgramResult(t, output, exiterr, success)
	}

	output, exiterr, success := program("context", "P0", "P1")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("context")
	assertProgramResult(t, output, exiterr, success)
	assert.Equal(t, "P0 P1", strings.TrimSpace(string(output)))

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)
	assert.Equal(t, []string{"task-P0", "task-P1"}, taskSummaries(unmarshalTaskArray(t, output)))

	// the query narrows down the priorities of the context
	output, exiterr, success = program("next", "P1")
	assertProgramResult(t, output, exiterr, success)
	assert.Equal(t, []string{"task-P1"}, taskSummaries(unmarshalTaskArray(t, output)))

	// a priority outside of the context is a conflict
	for _, args := range [][]string{
		{"next", "P2"},
		{"next", "P1", "P2"},
		{"add", "new-P3", "P3"},
	} {
		_, _, success = program(args...)
		assert.False(t, success, args)
	}

	// several priorities in the context are a filter only
	output, exiterr, success = program("add", "new")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "new-P1", "P1")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("log", "logged")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("template", "tmpl")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("start", "started")
	assertProgramResult(t, output, exiterr, success)

	priorities := map[string]string{}

	for _, cmd := range []string{"show-open", "show-resolved", "show-templates"} {
		output, exiterr, success = program(cmd, "--")
		assertProgramResult(t, output, exiterr, success)

		for _, task := range unmarshalTaskArray(t, output) {
			priorities[task.Summary] = task.Priority
		}
	}

	assert.Equal(t, "P2", priorities["new"])
	assert.Equal(t, "P1", priorities["new-P1"])
	assert.Equal(t, "P2", priorities["logged"])
	assert.Equal(t, "P2", priorities["tmpl"])
	assert.Equal(t, "P2", priorities["started"])
}

func TestSinglePriorityContextSetsPriority(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("context", "P1")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "new")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "P1", tasks[0].Priority)

	// unchanged: a priority in context and query is a conflict
	_, _, success = program("add", "other", "P1")
	assert.False(t, success)

	_, _, success = program("next", "P0", "P1")
	assert.False(t, success)
}

func TestMultiplePrioritiesRejected(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "task", "P2")
	assertProgramResult(t, output, exiterr, success)

	for _, args := range [][]string{
		{"add", "foo", "P0", "P1"},
		{"log", "foo", "P0", "P1"},
		{"template", "foo", "P0", "P1"},
		{"start", "foo", "P0", "P1"},
		{"1", "modify", "P0", "P1"},
		{"modify", "P0", "P1"},
	} {
		_, exiterr, success = program(args...)
		assert.False(t, success, args)
		if assert.NotNil(t, exiterr, args) {
			assert.Contains(t, string(exiterr.Stderr), "more than one priority", args)
		}
	}

	_, _, success = program("1", "P0", "P1")
	assert.False(t, success)

	// quoted text is not parsed as priority
	output, exiterr, success = program("add", "level P1", "P0")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("show-open", "--")
	assertProgramResult(t, output, exiterr, success)

	priorities := map[string]string{}
	for _, task := range unmarshalTaskArray(t, output) {
		priorities[task.Summary] = task.Priority
	}

	assert.Equal(t, map[string]string{"task": "P2", "level P1": "P0"}, priorities)
}

func TestModifyInMultiplePrioritiesContext(t *testing.T) {
	run, cleanup := makeRepoWithPriorities(t)
	defer cleanup()

	run("context", "P0", "P3")

	// without IDs, modify applies to all tasks in the context
	run("modify", "+foo")

	output := run("next", "--", "+foo")
	assert.Equal(t, []string{"task-P0", "task-P3"}, taskSummaries(unmarshalTaskArray(t, output)))
}

func TestMultiplePrioritiesContextFromEnvVar(t *testing.T) {
	run, cleanup := makeRepoWithPriorities(t)
	defer cleanup()

	unsetEnv := setEnv("DSTASK_CONTEXT", "P3 P0")
	defer unsetEnv()

	output := run("next")
	assert.Equal(t, []string{"task-P0", "task-P3"}, taskSummaries(unmarshalTaskArray(t, output)))
}
