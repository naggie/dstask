package dstask

// main task data structures

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type SubTask struct {
	Summary  string
	Resolved bool
}

// Task is our representation of tasks added at the command line and serialized
// to the task database on disk. It is rendered in multiple ways by the TaskSet
// to which it belongs.
type Task struct {
	// not stored in file -- rather filename and directory
	UUID   string `json:"uuid" yaml:"-"` // TODO: use actual uuid.UUID type here
	Status string `json:"status" yaml:",omitempty"`
	// is new or has changed. Need to write to disk.
	WritePending bool `json:"-" yaml:"-"`

	// ephemeral, used to address tasks quickly. Non-resolved only. Populated
	// from IDCache or on-the-fly.
	ID int `json:"id" yaml:"-"`

	// Deleted, if true, marks this task for deletion
	Deleted bool `json:"-" yaml:"-"`

	// concise representation of task
	Summary string `json:"summary"`
	// more detail, or information to remember to complete the task
	Notes   string   `json:"notes"`
	Tags    []string `json:"tags"`
	Project string   `json:"project"`
	// see const.go for PRIORITY_ strings
	Priority    string    `json:"priority"`
	DelegatedTo string    `json:"-"`
	Subtasks    []SubTask `json:"-"`
	// uuids of tasks that this task depends on
	// blocked status can be derived.
	// TODO possible filter: :blocked. Also, :overdue
	Dependencies []string `json:"-"`

	Created  time.Time `json:"created"`
	Resolved time.Time `json:"resolved"`
	Due      time.Time `json:"due"`

	// TaskSet uses this to indicate if a given task is excluded by a filter
	// (context etc)
	filtered bool `json:"-"`
}

// Equals returns whether t2 equals task.
// for equality, we ignore "core properties", not WritePending, ID, Deleted and filtered
func (t Task) Equals(t2 Task) bool {
	if t2.UUID != t.UUID {
		return false
	}
	if t2.Status != t.Status {
		return false
	}
	if t2.Summary != t.Summary {
		return false
	}
	if t2.Notes != t.Notes {
		return false
	}
	if !reflect.DeepEqual(t.Tags, t2.Tags) {
		return false
	}
	if t2.Project != t.Project {
		return false
	}
	if t2.Priority != t.Priority {
		return false
	}
	if t2.DelegatedTo != t.DelegatedTo {
		return false
	}
	if !reflect.DeepEqual(t.Subtasks, t2.Subtasks) {
		return false
	}
	if !reflect.DeepEqual(t.Dependencies, t2.Dependencies) {
		return false
	}
	if !t2.Created.Equal(t.Created) || !t2.Resolved.Equal(t.Resolved) || !t2.Due.Equal(t.Due) {
		return false
	}
	return true
}

// Unmarshal a Task from disk. We explicitly pass status, because the caller
// already knows the status, and can override the status declared in yaml.
func unmarshalTask(path string, finfo os.FileInfo, ids IdsMap, status string) (Task, error) {
	if len(finfo.Name()) != TASK_FILENAME_LEN {
		return Task{}, fmt.Errorf("filename does not encode UUID %s (wrong length)", finfo.Name())
	}

	uuid := finfo.Name()[0:36]
	if !IsValidUUID4String(uuid) {
		return Task{}, fmt.Errorf("filename does not encode UUID %s", finfo.Name())
	}

	t := Task{
		UUID:   uuid,
		Status: status,
		ID:     ids[uuid],
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Task{}, fmt.Errorf("Failed to read %s", finfo.Name())
	}
	err = yaml.Unmarshal(data, &t)
	if err != nil {
		return Task{}, fmt.Errorf("Failed to unmarshal %s", finfo.Name())
	}

	t.Status = status
	return t, nil
}

func (task Task) String() string {
	if task.ID > 0 {
		return fmt.Sprintf("%v: %s", task.ID, task.Summary)
	}
	return task.Summary
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

func (t *Task) SaveToDisk(repoPath string) {
	// save should be idempotent
	t.WritePending = false

	filepath := MustGetRepoPath(repoPath, t.Status, t.UUID+".yml")

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

		filepath := MustGetRepoPath(repoPath, st, t.UUID+".yml")

		if _, err := os.Stat(filepath); !os.IsNotExist(err) {
			err := os.Remove(filepath)
			if err != nil {
				ExitFail("Could not remove task %s: %v", filepath, err)
			}
		}
	}
}
