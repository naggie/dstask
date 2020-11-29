package dstask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func countFiltered(tasks []*Task) int {
	var numFiltered int
	for _, task := range tasks {
		if task.filtered {
			numFiltered++
		}
	}
	return numFiltered
}

func TestFilterTasksByID(t *testing.T) {

	makeTestTasks := func() []*Task {
		return []*Task{
			{ID: 1},
			{ID: 2},
			{ID: 3},
			{ID: 200},
			{ID: 500},
		}
	}

	t.Run("test without IDs filter", func(t *testing.T) {
		tso := taskSetOpts{
			withIDs: nil,
		}
		testTasks := makeTestTasks()
		filterTasksByID(testTasks, &tso)
		// no tasks are filtered, since no ids were passed
		assert.Equal(t, countFiltered(testTasks), len(testTasks))
	})

	t.Run("test with non-existent ID", func(t *testing.T) {
		tso := taskSetOpts{
			withIDs: []int{999},
		}
		testTasks := makeTestTasks()
		filterTasksByID(testTasks, &tso)
		// no tasks were filtered, since a non-existent ID was passed
		assert.Equal(t, countFiltered(testTasks), len(testTasks))
	})

	t.Run("test with 2 good IDs", func(t *testing.T) {
		tso := taskSetOpts{
			withIDs: []int{1, 2},
		}
		testTasks := makeTestTasks()
		filterTasksByID(testTasks, &tso)
		// all but 2 tasks filtered
		assert.Equal(t, countFiltered(testTasks), len(testTasks)-2)
	})

}
