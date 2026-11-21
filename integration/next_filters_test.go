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

func TestFilterTagWildcard(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "+team:hive", "hive")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "+team:reuss", "reuss")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "+teamwork", "teamwork")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "untagged")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "+team:*")
	assertProgramResult(t, output, exiterr, success)
	assert.ElementsMatch(t, []string{"hive", "reuss"}, taskSummaries(unmarshalTaskArray(t, output)))

	output, exiterr, success = program("next", "+team:h*")
	assertProgramResult(t, output, exiterr, success)
	assert.ElementsMatch(t, []string{"hive"}, taskSummaries(unmarshalTaskArray(t, output)))

	// +* matches tasks with any tag
	output, exiterr, success = program("next", "+*")
	assertProgramResult(t, output, exiterr, success)
	assert.ElementsMatch(t, []string{"hive", "reuss", "teamwork"},
		taskSummaries(unmarshalTaskArray(t, output)))
}

func TestFilterAntiTagWildcard(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "+team:hive", "hive")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "+teamwork", "teamwork")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "untagged")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "-team:*")
	assertProgramResult(t, output, exiterr, success)
	assert.ElementsMatch(t, []string{"teamwork", "untagged"},
		taskSummaries(unmarshalTaskArray(t, output)))

	// -* matches tasks without any tag
	output, exiterr, success = program("next", "-*")
	assertProgramResult(t, output, exiterr, success)
	assert.ElementsMatch(t, []string{"untagged"}, taskSummaries(unmarshalTaskArray(t, output)))
}

func TestTagWildcardNotStored(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	_, _, success := program("add", "foo", "+team:*")
	assert.False(t, success)

	output, exiterr, success := program("add", "+team:hive", "hive")
	assertProgramResult(t, output, exiterr, success)

	_, _, success = program("modify", "1", "+team:*")
	assert.False(t, success)

	_, _, success = program("modify", "1", "-team:*")
	assert.False(t, success)

	// wildcard in context filters, but is not added to new tasks
	output, exiterr, success = program("context", "+team:*")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "+team:reuss", "reuss")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next")
	assertProgramResult(t, output, exiterr, success)

	tasks := unmarshalTaskArray(t, output)
	assert.ElementsMatch(t, []string{"hive", "reuss"}, taskSummaries(tasks))
	for _, task := range tasks {
		assert.NotContains(t, task.Tags, "team:*")
	}
}

func TestWildcardTagsRejected(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	_, _, success := program("log", "foo", "+team:*")
	assert.False(t, success)

	_, _, success = program("template", "foo", "+team:*")
	assert.False(t, success)

	_, _, success = program("start", "foo", "+team:*")
	assert.False(t, success)
}

func TestModifyInWildcardContext(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	program := testCmd(repo)

	output, exiterr, success := program("add", "+team:hive", "hive")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("add", "untagged")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("context", "+team:*")
	assertProgramResult(t, output, exiterr, success)

	// without IDs, modify applies to all tasks in the context
	output, exiterr, success = program("modify", "+foo")
	assertProgramResult(t, output, exiterr, success)

	output, exiterr, success = program("next", "--", "+foo")
	assertProgramResult(t, output, exiterr, success)
	assert.Equal(t, []string{"hive"}, taskSummaries(unmarshalTaskArray(t, output)))
}

func taskSummaries(tasks []dstask.Task) []string {
	var result []string
	for _, task := range tasks {
		result = append(result, task.Summary)
	}
	return result
}
