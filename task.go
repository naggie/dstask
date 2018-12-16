package dstask

// main task data structures

import (
	"sort"
	"time"
	"strconv"
	"strings"
)

const (
	STATUS_PENDING   = "pending"
	STATUS_ACTIVE    = "active"
	STATUS_RESOLVED  = "resolved"
	STATUS_DELEGATED = "delegated"
	STATUS_DEFERRED  = "deferred"
	STATUS_SOMEDAY   = "someday"
	STATUS_RECURRING = "recurring" // tentative

	GIT_REPO = "~/.dstask/"
	// space delimited keyword file for compgen
	COMPLETION_FILE = "~/.cache/dstask/completions"

	// filter: P1 P2 etc
	PRIORITY_CRITICAL = "P1"
	PRIORITY_HIGH     = "P2"
	PRIORITY_NORMAL   = "P3"
	PRIORITY_LOW      = "P4"

	MAX_TASKS_OPEN    = 10000
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
var NORMAL_STATUSES = []string{
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
	uuid   string
	status string
	// ephemeral, used to address tasks quickly. Non-resolved only.
	id int

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
	tasks          []*Task

	// indices
	tasksByID      map[int]*Task
	tasksByUuid    map[string]*Task

	CurrentContext string
}

func (ts *TaskSet) SortTaskList() {
	sort.SliceStable(ts.tasks, func(i, j int) bool { return ts.tasks[i].Created.Before(ts.tasks[j].Created) })
	sort.SliceStable(ts.tasks, func(i, j int) bool { return ts.tasks[i].Priority < ts.tasks[j].Priority })
	sort.SliceStable(ts.tasks, func(i, j int) bool { return STATUS_ORDER[ts.tasks[i].status] < STATUS_ORDER[ts.tasks[j].status] })
}

// add a task, but only if it has a new uuid or no uuid. Return true if task
// was added.
func (ts *TaskSet) AddTask(task *Task) bool {
	if ts.tasksByUuid[task.uuid] != nil {
		// load tasks, do not overwrite
		return false
	}

	// resolved task should not have ID
	if task.status != STATUS_RESOLVED {
		task.id = 0
	}

	// check ID is unique if there is one
	if task.id > 0 && ts.tasksByID[task.id] != nil {
		task.id = 0
	}

	// pick one if task isn't resolved and ID isn't there
	if task.status != STATUS_RESOLVED {
		for id:=1; id <= MAX_TASKS_OPEN; id++ {
			if ts.tasksByID[id] == nil {
				task.id = id
				break
			}
		}
	}

	ts.tasks = append(ts.tasks, task)
	ts.tasksByUuid[task.uuid] = task
	ts.tasksByID[task.id] = task
	return true
}

// when refering to tasks by ID, NORMAL_STATUSES must be loaded exclusively --
// even if the filter is set to show issues that have only some statuses.
type TaskLine struct {
	// operation, taken from first occurance. Mapped to a method to operate on
	// tasks.
	Action string
	Id       int
	Tags     []string
	AntiTags []string
	Project  string
	Priority string
	Text     string
}

func ParseTaskLine(args []string) *TaskLine {
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

	return &TaskLine{
		Id:       id,
		Tags:     tags,
		AntiTags: antiTags,
		Project:  project,
		Priority: priority,
		Text:     strings.Join(words, " "),
	}
}

func (ts *TaskSet) Filter(tl *TaskLine) {
	var tasks []*Task

	for _, t := range(ts.tasks) {
		if t.Project != tl.Project {
			return
		}

		tasks = append(tasks, t)
	}

	ts.tasks = tasks
}
