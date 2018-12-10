package dstask

import (
	"time"
)

type SubTask struct {
	Summary  string
	Resolved bool
}

type Task struct {
	// not stored in file -- rather filename and directory
	uuid   string
	status string
	// used to determine if an unlink should happen if status changes
	originalFilepath string

	// concise representation of task
	Summary string
	// task in more detail, only if necessary
	Description string
	Tags        []string
	Project     string
	// see const.go for PRIORITY_ strings
	Priority    string
	DelegatedTo string
	Subtasks    []SubTask
	Comments    []string
	// uuids of tasks that this task depends on
	// blocked status can be derived.
	// TODO possible filter: :blocked. Also, :overdue
	Dependencies []string

	Created  time.Time
	Modified time.Time
	Resolved time.Time
	Due      time.Time
}

type TaskSet struct {
	Tasks          []Task
	CurrentContext string
	knownUuids     map[string]bool
}

func NewTaskSet() *TaskSet {
	return &TaskSet{
		knownUuids: make(map[string]bool),
	}
}

// add a task, but only if it has a new uuid. Return true if task was added.
func (ts *TaskSet) MaybeAddTask(task Task) bool {
	if ts.knownUuids[task.uuid] {
		// load tasks, do not overwrite
		return false
	}

	ts.knownUuids[task.uuid] = true
	ts.Tasks = append(ts.Tasks, task)
	return true
}

// filter should be set before loading any data. The filter can be used to
// optimise a bit -- eg when listing, completed tasks should not be shown so we
// can avoid loading them. However when importing, it is important to load all
// tasks for full context.
type TaskFilter struct {
	Text     string
	Tags     []string
	Antitags []string
	Project  string
	Priority int
	Statuses []string
}

//func (ts *TaskSet) filter(filter *TaskFilter) TaskSet {
//
//}
//
//func (t *Task) Save() error {
//
//}
