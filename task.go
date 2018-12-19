package dstask

// main task data structures

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	STATUS_PENDING   = "pending"
	STATUS_ACTIVE    = "active"
	STATUS_RESOLVED  = "resolved"
	STATUS_DELEGATED = "delegated"
	STATUS_DEFERRED  = "deferred"
	STATUS_SOMEDAY   = "someday"
	STATUS_RECURRING = "recurring" // tentative

	// filter: P1 P2 etc
	PRIORITY_CRITICAL = "P1"
	PRIORITY_HIGH     = "P2"
	PRIORITY_NORMAL   = "P3"
	PRIORITY_LOW      = "P4"

	MAX_TASKS_OPEN = 10000
)

// for import (etc) it's necessary to have full context
var ALL_STATUSES = []string{
	STATUS_ACTIVE,
	STATUS_PENDING,
	STATUS_DELEGATED,
	STATUS_DEFERRED,
	STATUS_SOMEDAY,
	STATUS_RECURRING,
	STATUS_RESOLVED,
}

// for most operations, it's not necessary or desirable to load the expensive resolved tasks
var NON_RESOLVED_STATUSES = []string{
	STATUS_ACTIVE,
	STATUS_PENDING,
	STATUS_DELEGATED,
	STATUS_DEFERRED,
	STATUS_SOMEDAY,
	STATUS_RECURRING,
}

// TODO consider using iota enum for statuses, with custom marshaller
// https://gist.github.com/lummie/7f5c237a17853c031a57277371528e87
// though this seems simpler
var STATUS_ORDER = map[string]int{
	STATUS_ACTIVE:    1,
	STATUS_PENDING:   2,
	STATUS_DELEGATED: 3,
	STATUS_DEFERRED:  4,
	STATUS_SOMEDAY:   5,
	STATUS_RECURRING: 6,
	STATUS_RESOLVED:  7,
}

type SubTask struct {
	Summary  string
	Resolved bool
}

type Task struct {
	// not stored in file -- rather filename and directory
	Uuid   string `yaml:"-"`
	Status string `yaml:"-"`
	// is new or has changed. Need to write to disk.
	WritePending bool `yaml:"-"`

	// ephemeral, used to address tasks quickly. Non-resolved only.
	ID int

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
	tasks []*Task

	// indices
	tasksByID   map[int]*Task
	tasksByUuid map[string]*Task

	CurrentContext string
}

func (ts *TaskSet) SortTaskList() {
	sort.SliceStable(ts.tasks, func(i, j int) bool { return ts.tasks[i].Created.Before(ts.tasks[j].Created) })
	sort.SliceStable(ts.tasks, func(i, j int) bool { return ts.tasks[i].Priority < ts.tasks[j].Priority })
	sort.SliceStable(ts.tasks, func(i, j int) bool { return STATUS_ORDER[ts.tasks[i].Status] < STATUS_ORDER[ts.tasks[j].Status] })
}

// add a task, but only if it has a new uuid or no uuid. Return true if task
// was added.
func (ts *TaskSet) AddTask(task Task) bool {
	if task.Uuid == "" {
		task.Uuid = MustGetUuid4String()
	}

	if ts.tasksByUuid[task.Uuid] != nil {
		// load tasks, do not overwrite
		return false
	}

	// resolved task should not have ID
	if task.Status == STATUS_RESOLVED {
		task.ID = 0
	}

	// check ID is unique if there is one
	if task.ID > 0 && ts.tasksByID[task.ID] != nil {
		task.ID = 0
	}

	// pick one if task isn't resolved and ID isn't there
	if task.ID == 0 && task.Status != STATUS_RESOLVED {
		for id := 1; id <= MAX_TASKS_OPEN; id++ {
			if ts.tasksByID[id] == nil {
				task.ID = id
				task.WritePending = true
				break
			}
		}
	}

	if task.Created.IsZero() {
		task.Created = time.Now()
		task.WritePending = true
	}

	ts.tasks = append(ts.tasks, &task)
	ts.tasksByUuid[task.Uuid] = &task
	ts.tasksByID[task.ID] = &task
	return true
}

// TODO maybe this is the place to check for invalid state transitions instead
// of the main switch statement. Though, a future 3rdparty sync system could
// need this to work regardless.
func (ts *TaskSet) MustUpdateTask(task Task) {
	if ts.tasksByUuid[task.Uuid] == nil {
		ExitFail("Could not find given task to update by UUID")
	}

	task.WritePending = true

	if task.Status == STATUS_RESOLVED {
		task.ID = 0
	}

	// existing pointer must point to address of new task copied
	*ts.tasksByUuid[task.Uuid] = task
}

// when refering to tasks by ID, NON_RESOLVED_STATUSES must be loaded exclusively --
// even if the filter is set to show issues that have only some statuses.
type TaskLine struct {
	ID       int
	Tags     []string
	AntiTags []string
	Project  string
	Priority string
	Text     string
}

func ParseTaskLine(args ...string) TaskLine {
	var id int
	var tags []string
	var antiTags []string
	var project string
	var priority string
	var words []string

	for i, item := range args {
		if s, err := strconv.ParseInt(item, 10, 64); i == 0 && err == nil {
			id = int(s)
		} else if strings.HasPrefix(item, "project:") {
			project = item[8:]
		} else if len(item) > 2 && item[0:1] == "+" {
			tags = append(tags, item[1:])
		} else if len(item) > 2 && item[0:1] == "-" {
			antiTags = append(antiTags, item[1:])
		} else if IsValidPriority(item) {
			priority = item
		} else {
			words = append(words, item)
		}
	}

	return TaskLine{
		ID:       id,
		Tags:     tags,
		AntiTags: antiTags,
		Project:  project,
		Priority: priority,
		Text:     strings.Join(words, " "),
	}
}

func (ts *TaskSet) Filter(tl TaskLine) {
	var tasks []*Task

	for _, task := range ts.tasks {
		if task.MatchesFilter(tl) {
			tasks = append(tasks, task)
		}
	}

	ts.tasks = tasks
}

func (t *Task) MatchesFilter(tl TaskLine) bool {
	if tl.ID != 0  && t.ID == tl.ID {
		return true
	}

	for _, tag := range(tl.Tags) {
		if !StrSliceContains(t.Tags, tag) {
			return false
		}
	}

	for _, tag := range(tl.AntiTags) {
		if StrSliceContains(t.Tags, tag) {
			return false
		}
	}

	if tl.Project != "" && t.Project != tl.Project {
		return false
	}

	if tl.Priority != "" && t.Priority != tl.Priority {
		return false
	}

	if tl.Text != "" && !strings.Contains(strings.ToLower(t.Summary + t.Description), strings.ToLower(tl.Text)) {
		return false
	}

	return true
}

func (ts *TaskSet) MustGetByID(id int) Task {
	if ts.tasksByID[id] == nil {
		ExitFail("No open task with that ID exists.")
	}

	return *ts.tasksByID[id]
}
