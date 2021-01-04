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

	// Urgency is generated on the fly
	Urgency int `json:"urgency" yaml:"-"`

	// TaskSet uses this to indicate if a given task is excluded by a filter
	// (context etc)
	filtered bool `json:"-"`
}

// Equals returns whether t2 equals task.
// for equality, we only consider "core properties", we ignore WritePending, ID, Deleted and filtered
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
		return Task{}, fmt.Errorf("failed to read %s", finfo.Name())
	}
	err = yaml.Unmarshal(data, &t)
	if err != nil {
		return Task{}, fmt.Errorf("failed to unmarshal %s", finfo.Name())
	}

	t.Status = status
	t.Urgency = computeUrgency(t)
	return t, nil
}

func (t Task) String() string {
	if t.ID > 0 {
		return fmt.Sprintf("%v: %s", t.ID, t.Summary)
	}
	return t.Summary
}

func (t *Task) MatchesFilter(query Query) bool {
	// IDs were specified but none match (OR logic)
	if len(query.IDs) > 0 && !IntSliceContains(query.IDs, t.ID) {
		return false
	}

	for _, tag := range query.Tags {
		if !StrSliceContains(t.Tags, tag) {
			return false
		}
	}

	for _, tag := range query.AntiTags {
		if StrSliceContains(t.Tags, tag) {
			return false
		}
	}

	if StrSliceContains(query.AntiProjects, t.Project) {
		return false
	}

	if query.Project != "" && t.Project != query.Project {
		return false
	}

	if query.Priority != "" && t.Priority != query.Priority {
		return false
	}

	if query.Text != "" && !strings.Contains(strings.ToLower(t.Summary+t.Notes), strings.ToLower(query.Text)) {
		return false
	}

	return true
}

// Normalise mutates and sorts some of a task object's fields into a consistent
// format. This should make git diffs more useful.
func (t *Task) Normalise() {
	t.Project = strings.ToLower(t.Project)

	// tags must be lowercase
	for i, tag := range t.Tags {
		t.Tags[i] = strings.ToLower(tag)
	}

	// tags must be sorted
	sort.Strings(t.Tags)

	// tags must be unique
	t.Tags = DeduplicateStrings(t.Tags)

	if t.Status == STATUS_RESOLVED {
		// resolved task should not have ID as it's meaningless
		t.ID = 0
	}

	if t.Priority == "" {
		t.Priority = PRIORITY_NORMAL
	}
}

// normalise the task before validating!
func (t *Task) Validate() error {
	if !IsValidUUID4String(t.UUID) {
		return errors.New("invalid task UUID4")
	}

	if !IsValidStatus(t.Status) {
		return errors.New("invalid status specified on task")
	}

	if !IsValidPriority(t.Priority) {
		return errors.New("invalid priority specified")
	}

	for _, uuid := range t.Dependencies {
		if !IsValidUUID4String(uuid) {
			return errors.New("invalid dependency UUID4")
		}
	}

	return nil
}

// provides Summary + Last note if available
func (t *Task) LongSummary() string {
	noteLines := strings.Split(t.Notes, "\n")
	lastNote := noteLines[len(noteLines)-1]

	if len(lastNote) > 0 {
		return t.Summary + " " + NOTE_MODE_KEYWORD + " " + lastNote
	}
	return t.Summary
}

func (t *Task) Modify(query Query) {
	for _, tag := range query.Tags {
		if !StrSliceContains(t.Tags, tag) {
			t.Tags = append(t.Tags, tag)
		}
	}

	for i, tag := range t.Tags {
		if StrSliceContains(query.AntiTags, tag) {
			// delete item
			t.Tags = append(t.Tags[:i], t.Tags[i+1:]...)
		}
	}

	if query.Project != "" {
		t.Project = query.Project
	}

	if StrSliceContains(query.AntiProjects, t.Project) {
		t.Project = ""
	}

	if query.Priority != "" {
		t.Priority = query.Priority
	}
}

func (task *Task) IsResolved() bool {
	for _, nonResolvedStatus := range NON_RESOLVED_STATUSES {
		if task.Status == nonResolvedStatus {
			return false
		}
	}
	return true
}

// If you make changes to this function, make sure you update doc/urgency.md to
// reflect the changes.
func computeUrgency(task Task) int {
	if task.IsResolved() {
		return 0
	}

	urgency := 1

	priorityModifier := 5
	switch task.Priority {
	case PRIORITY_LOW:
		urgency += priorityModifier * 1
	case PRIORITY_NORMAL:
		urgency += priorityModifier * 2
	case PRIORITY_HIGH:
		urgency += priorityModifier * 3
	case PRIORITY_CRITICAL:
		urgency += priorityModifier * 5
	}

	if task.Status == STATUS_ACTIVE {
		urgency += 5
	}

	if len(task.Project) > 0 {
		urgency += 3
	}

	if len(task.Tags) > 0 {
		urgency += 3
	}

	ageModifier := 0.05
	ageInDays := int(time.Since(task.Created).Hours() / 24)
	urgency += int(float64(ageInDays) * ageModifier)

	return urgency
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
