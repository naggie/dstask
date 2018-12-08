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
	assignedTo string
	subtasks []SubTask
	comments []string

	created time.Time
	modified time.Time
	resolved time.Time
	due time.Time
}
