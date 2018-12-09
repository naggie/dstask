package dstask

import (
	"time"
)

type SubTask struct {
	summary  string
	resolved bool
}

type Task struct {
	// not stored in file -- rather filename and directory
	uuid   string
	status string
	// used to determine if an unlink should happen if status changes
	originalFilepath string

	// concise representation of task
	Summary     string
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
	tasks          []Task
	currentContext string
}

type TaskFilter struct {
	text     string
	tags     []string
	antitags []string
	project  string
	priority int
}

//func (ts *TaskSet) filter(filter *TaskFilter) TaskSet {
//
//}
//
//func (t *Task) Save() error {
//
//}
