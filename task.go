package dstask

// main task data structures

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

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
	ID int `yaml:",omitempty"`

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
}

type TaskSet struct {
	tasks []*Task

	// indices
	tasksByID   map[int]*Task
	tasksByUuid map[string]*Task

	CurrentContext string
}

func (task Task) String() string {
	return fmt.Sprintf("%v: %s", task.ID, task.Summary)
}

func (ts *TaskSet) SortTaskList() {
	sort.SliceStable(ts.tasks, func(i, j int) bool { return ts.tasks[i].Created.Before(ts.tasks[j].Created) })
	sort.SliceStable(ts.tasks, func(i, j int) bool { return ts.tasks[i].Priority < ts.tasks[j].Priority })
	sort.SliceStable(ts.tasks, func(i, j int) bool { return STATUS_ORDER[ts.tasks[i].Status] < STATUS_ORDER[ts.tasks[j].Status] })
}

// add a task, but only if it has a new uuid or no uuid. Return annotated task.
func (ts *TaskSet) AddTask(task Task) Task {
	if task.Uuid == "" {
		task.Uuid = MustGetUuid4String()
	}

	if ts.tasksByUuid[task.Uuid] != nil {
		// load tasks, do not overwrite
		return Task{}
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

	if task.Priority == "" {
		task.Priority = PRIORITY_NORMAL
	}

	if task.Created.IsZero() {
		task.Created = time.Now()
		task.WritePending = true
	}

	ts.tasks = append(ts.tasks, &task)
	ts.tasksByUuid[task.Uuid] = &task
	ts.tasksByID[task.ID] = &task
	return task
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
type CmdLine struct {
	Cmd      string
	IDs      []int
	Tags     []string
	AntiTags []string
	Project  string
	Priority string
	Text     string
}

// used for applying a context to a new task
func (tl *CmdLine) MergeContext(_tl CmdLine) {
	for _, tag := range _tl.Tags {
		if !StrSliceContains(tl.Tags, tag) {
			tl.Tags = append(tl.Tags, tag)
		}
	}

	for _, tag := range _tl.AntiTags {
		if !StrSliceContains(tl.AntiTags, tag) {
			tl.AntiTags = append(tl.AntiTags, tag)
		}
	}

	// TODO same for antitags
	if _tl.Project != "" {
		if tl.Project != "" {
			ExitFail("Could not apply context, project conflict")
		} else {
			tl.Project = _tl.Project
		}
	}

	if _tl.Priority != "" {
		if tl.Priority != "" {
			ExitFail("Could not apply context, priority conflict")
		} else {
			tl.Priority = _tl.Priority
		}
	}
}

// reconstruct args string
func (tl CmdLine) String() string {
	var args []string
	var annotatedTags []string

	if tl.ID > 0 {
		args = append(args, strconv.Itoa(tl.ID))
	}

	for _, tag := range tl.Tags {
		annotatedTags = append(annotatedTags, "+"+tag)
	}
	for _, tag := range tl.AntiTags {
		annotatedTags = append(annotatedTags, "-"+tag)
	}
	args = append(args, annotatedTags...)

	if tl.Project != "" {
		args = append(args, "project:"+tl.Project)
	}

	if tl.Priority != "" {
		args = append(args, tl.Priority)
	}

	if tl.Text != "" {
		args = append(args, "\""+tl.Text+"\"")
	}

	return strings.Join(args, " ")
}

func ParseCmdLine(args ...string) CmdLine {
	var cmd string
	var ids []int
	var tags []string
	var antiTags []string
	var project string
	var priority string
	var words []string

	// something other than an ID has been parsed -- accept no more IDs
	var idsExhausted bool

	for i, item := range args {
		if i == 0 && StrSliceContains(ALL_CMDS, item) {
			cmd = item
			continue
		}

		if s, err := strconv.ParseInt(item, 10, 64); !idsExhausted && err == nil {
			ids = append(ids, int(s))
			continue
		}

		idsExhausted = true

		if strings.HasPrefix(item, "project:") {
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

	if cmd == "" {
		cmd = CMD_NEXT
	}

	return CmdLine{
		Cmd:      cmd,
		IDs:      ids,
		Tags:     tags,
		AntiTags: antiTags,
		Project:  project,
		Priority: priority,
		Text:     strings.Join(words, " "),
	}
}

func (ts *TaskSet) Filter(tl CmdLine) {
	var tasks []*Task

	for _, task := range ts.tasks {
		if task.MatchesFilter(tl) {
			tasks = append(tasks, task)
		}
	}

	ts.tasks = tasks
}

func (t *Task) MatchesFilter(tl CmdLine) bool {
	if tl.ID != 0 && t.ID == tl.ID {
		return true
	}

	for _, tag := range tl.Tags {
		if !StrSliceContains(t.Tags, tag) {
			return false
		}
	}

	for _, tag := range tl.AntiTags {
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

	if tl.Text != "" && !strings.Contains(strings.ToLower(t.Summary+t.Notes), strings.ToLower(tl.Text)) {
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
