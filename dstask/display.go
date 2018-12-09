package dstask

import (
	"fmt"
)

// should use a better console library after first POC

/// display list of filtered tasks with context and filter
func (ts *TaskSet) Display() {
	for n, t := range(ts.Tasks) {
		fmt.Println(n, t)
	}
}

// display a single task in detail, with numbered subtasks
func (t *Task) Display() {

}
