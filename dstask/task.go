package dstask

import (
	time
)

const (
	PENDING = "pending"
	ACTIVE = "active"
	RESOLVED = "resolved"
	DELEGATED = "delegated"
	DEFERRED = "deferred"
	SOMEDAY = "someday"
)

type SubTask struct {
	summary string
	resolved bool
}

type Task struct {
	uuid string
	status string
	summary string
	description string
	tags []string
	project string
	priority int
	delegatedTo string
	subtasks []SubTask
	comments []string
	// uuids of tasks that this task depends on
	// blocked status can be derived.
	// TODO possible filter: :blocked. Also, :overdue
	dependencies []string

	created time.Time
	modified time.Time
	resolved time.Time
	due time.Time
}

type TaskSet struct {
	tasks []Task
}

type TaskFilter struct {
	text string
	tags []string
	antitags []string
	project string
	priority int
}

func (ts *TaskSet) filter(filter *TaskFilter) TaskSet {

}

func (t *Task) Save() error {

}
