package dstask

// main task data structures

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type SubTask struct {
	Summary  string
	Resolved bool
}

type Task struct {
	// not stored in file -- rather filename and directory
	UUID   string `yaml:"-"`
	Status string `yaml:",omitempty"`
	// is new or has changed. Need to write to disk.
	WritePending bool `yaml:"-"`

	// ephemeral, used to address tasks quickly. Non-resolved only. Populated
	// from IDCache or on-the-fly.
	ID int `yaml:"-"`

	// Deleted, if true, marks this task for deletion
	Deleted bool `yaml:"-"`

	// concise representation of task
	Summary string
	// more detail, or information to remember to complete the task
	Notes   string
	Tags    []string
	Project string
	// see const.go for PRIORITY_ strings
	Priority    string
	DelegatedTo string
	Subtasks    []SubTask
	// uuids of tasks that this task depends on
	// blocked status can be derived.
	// TODO possible filter: :blocked. Also, :overdue
	Dependencies []string

	Created  time.Time
	Resolved time.Time
	Due      time.Time

	// TaskSet uses this to indicate if a given task is excluded by a filter
	// (context etc)
	filtered bool

	Parent string
	// only valid for recurring tasks, and resolved recurring tasks. Tasks
	// created by a recurring task should have this removed or will fail
	// validation. syntax: cron.
	Schedule string `yaml:"omitempty"`
	// Recurring task this was derived from. Used by scheduler to gate new
	// tasks.
	Parent   string `yaml:"omitempty"`
}

func (task Task) String() string {
	if task.ID > 0 {
		return fmt.Sprintf("%v: %s", task.ID, task.Summary)
	} else {
		return task.Summary
	}
}

// used for applying a context to a new task
func (cmdLine *CmdLine) MergeContext(context CmdLine) {
	for _, tag := range context.Tags {
		if !StrSliceContains(cmdLine.Tags, tag) {
			cmdLine.Tags = append(cmdLine.Tags, tag)
		}
	}

	for _, tag := range context.AntiTags {
		if !StrSliceContains(cmdLine.AntiTags, tag) {
			cmdLine.AntiTags = append(cmdLine.AntiTags, tag)
		}
	}

	// TODO same for antitags
	if context.Project != "" {
		if cmdLine.Project != "" && cmdLine.Project != context.Project {
			ExitFail("Could not apply context, project conflict")
		} else {
			cmdLine.Project = context.Project
		}
	}

	if context.Priority != "" {
		if cmdLine.Priority != "" {
			ExitFail("Could not apply context, priority conflict")
		} else {
			cmdLine.Priority = context.Priority
		}
	}
}

func (task *Task) MatchesFilter(cmdLine CmdLine) bool {
	for _, id := range cmdLine.IDs {
		if id == task.ID {
			return true
		}
	}

	// IDs were specified but no match
	if len(cmdLine.IDs) > 0 {
		return false
	}

	for _, tag := range cmdLine.Tags {
		if !StrSliceContains(task.Tags, tag) {
			return false
		}
	}

	for _, tag := range cmdLine.AntiTags {
		if StrSliceContains(task.Tags, tag) {
			return false
		}
	}

	if StrSliceContains(cmdLine.AntiProjects, task.Project) {
		return false
	}

	if cmdLine.Project != "" && task.Project != cmdLine.Project {
		return false
	}

	if cmdLine.Priority != "" && task.Priority != cmdLine.Priority {
		return false
	}

	if cmdLine.Text != "" && !strings.Contains(strings.ToLower(task.Summary+task.Notes), strings.ToLower(cmdLine.Text)) {
		return false
	}

	return true
}

// Normalise mutates and sorts some of a task object's fields into a consistent
// format. This should make git diffs more useful.
func (task *Task) Normalise() {
	task.Project = strings.ToLower(task.Project)

	// tags must be lowercase
	for i, tag := range task.Tags {
		task.Tags[i] = strings.ToLower(tag)
	}

	// tags must be sorted
	sort.Strings(task.Tags)

	// tags must be unique
	task.Tags = DeduplicateStrings(task.Tags)

	if task.Status == STATUS_RESOLVED {
		// resolved task should not have ID as it's meaningless
		task.ID = 0
	}

	if task.Priority == "" {
		task.Priority = PRIORITY_NORMAL
	}
}

// normalise the task before validating!
func (task *Task) Validate() error {
	if !IsValidUUID4String(task.UUID) {
		return errors.New("Invalid task UUID4")
	}

	if !IsValidStatus(task.Status) {
		return errors.New("Invalid status specified on task")
	}

	if !IsValidPriority(task.Priority) {
		return errors.New("Invalid priority specified")
	}

	for _, uuid := range task.Dependencies {
		if !IsValidUUID4String(uuid) {
			return errors.New("Invalid dependency UUID4")
		}
	}

	return nil
}

// provides Summary + Last note if available
func (task *Task) LongSummary() string {
	noteLines := strings.Split(task.Notes, "\n")
	lastNote := noteLines[len(noteLines)-1]

	if len(lastNote) > 0 {
		return task.Summary + " " + NOTE_MODE_KEYWORD + " " + lastNote
	} else {
		return task.Summary
	}
}

func (task *Task) Modify(cmdLine CmdLine) {
	for _, tag := range cmdLine.Tags {
		if !StrSliceContains(task.Tags, tag) {
			task.Tags = append(task.Tags, tag)
		}
	}

	for i, tag := range task.Tags {
		if StrSliceContains(cmdLine.AntiTags, tag) {
			// delete item
			task.Tags = append(task.Tags[:i], task.Tags[i+1:]...)
		}
	}

	if cmdLine.Project != "" {
		task.Project = cmdLine.Project
	}

	if StrSliceContains(cmdLine.AntiProjects, task.Project) {
		task.Project = ""
	}

	if cmdLine.Priority != "" {
		task.Priority = cmdLine.Priority
	}
}

func (t *Task) SaveToDisk() {
	// save should be idempotent
	t.WritePending = false

	filepath := MustGetRepoPath(t.Status, t.UUID+".yml")

	if t.Deleted {
		// Task is marked deleted. Delete from its current status directory.
		if err := os.Remove(filepath); err != nil {
			ExitFail("Could not remove task %s: %v", filepath, err)
		}

	} else {
		// Task is not deleted, and will be written to disk to a directory
		// that indicates its current status. We make a shallow copy first,
		// and we set Status to empty string. This shallow copy is serialised
		// to disk, with the Status field omitted. This avoids redundant data.
		taskCp := *t
		taskCp.Status = ""
		d, err := yaml.Marshal(&taskCp)
		if err != nil {
			// TODO present error to user, specific error message is important
			ExitFail("Failed to marshal task %s", t)
		}

		err = ioutil.WriteFile(filepath, d, 0600)
		if err != nil {
			ExitFail("Failed to write task %s", t)
		}
	}

	// Delete task from other status directories. Only one copy should exist, at most.
	for _, st := range ALL_STATUSES {
		if st == t.Status {
			continue
		}

		filepath := MustGetRepoPath(st, t.UUID+".yml")

		if _, err := os.Stat(filepath); !os.IsNotExist(err) {
			err := os.Remove(filepath)
			if err != nil {
				ExitFail("Could not remove task %s: %v", filepath, err)
			}
		}
	}
}
