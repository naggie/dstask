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
			&Task{ID: 1},
			&Task{ID: 2},
			&Task{ID: 3},
			&Task{ID: 200},
			&Task{ID: 500},
		}
	}

	t.Run("test without IDs filter", func(t *testing.T) {
		tso := taskSetOpts{
			withIDs: nil,
		}
		testTasks := makeTestTasks()
		filterTasksByID(testTasks, &tso)
		// every task should be filtered
		assert.Equal(t, countFiltered(testTasks), len(testTasks))
	})

	t.Run("test with non-existent ID", func(t *testing.T) {
		tso := taskSetOpts{
			withIDs: []int{999},
		}
		testTasks := makeTestTasks()
		filterTasksByID(testTasks, &tso)
		// every task should be filtered
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
