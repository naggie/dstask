package dstask

// main task data structures

import (
	"fmt"
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
	UUID   string `yaml:"-"`
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

func (task Task) String() string {
	return fmt.Sprintf("%v: %s", task.ID, task.Summary)
}

// when refering to tasks by ID, NON_RESOLVED_STATUSES must be loaded exclusively --
// even if the filter is set to show issues that have only some statuses.
type CmdLine struct {
	Cmd           string
	IDs           []int
	Tags          []string
	AntiTags      []string
	Project       string
	AntiProjects  []string
	Priority      string
	Text          string
	IgnoreContext bool
}

// used for applying a context to a new task
func (cmdLine *CmdLine) MergeContext(_tl CmdLine) {
	for _, tag := range _tl.Tags {
		if !StrSliceContains(cmdLine.Tags, tag) {
			cmdLine.Tags = append(cmdLine.Tags, tag)
		}
	}

	for _, tag := range _tl.AntiTags {
		if !StrSliceContains(cmdLine.AntiTags, tag) {
			cmdLine.AntiTags = append(cmdLine.AntiTags, tag)
		}
	}

	// TODO same for antitags
	if _tl.Project != "" {
		if cmdLine.Project != "" {
			ExitFail("Could not apply context, project conflict")
		} else {
			cmdLine.Project = _tl.Project
		}
	}

	if _tl.Priority != "" {
		if cmdLine.Priority != "" {
			ExitFail("Could not apply context, priority conflict")
		} else {
			cmdLine.Priority = _tl.Priority
		}
	}
}

// reconstruct args string
func (cmdLine CmdLine) String() string {
	var args []string
	var annotatedTags []string

	for _, id := range cmdLine.IDs {
		args = append(args, strconv.Itoa(id))
	}

	for _, tag := range cmdLine.Tags {
		annotatedTags = append(annotatedTags, "+"+tag)
	}
	for _, tag := range cmdLine.AntiTags {
		annotatedTags = append(annotatedTags, "-"+tag)
	}
	args = append(args, annotatedTags...)

	if cmdLine.Project != "" {
		args = append(args, "project:"+cmdLine.Project)
	}

	if cmdLine.Priority != "" {
		args = append(args, cmdLine.Priority)
	}

	if cmdLine.Text != "" {
		args = append(args, "\""+cmdLine.Text+"\"")
	}

	return strings.Join(args, " ")
}

func ParseCmdLine(args ...string) CmdLine {
	var cmd string
	var ids []int
	var tags []string
	var antiTags []string
	var project string
	var antiProjects []string
	var priority string
	var words []string
	var ignoreContext bool

	// something other than an ID has been parsed -- accept no more IDs
	var idsExhausted bool

	for _, item := range args {
		lcItem := strings.ToLower(item)
		if !idsExhausted && StrSliceContains(ALL_CMDS, lcItem) {
			cmd = lcItem
			continue
		}

		if s, err := strconv.ParseInt(item, 10, 64); !idsExhausted && err == nil {
			ids = append(ids, int(s))
			continue
		}

		idsExhausted = true

		if strings.HasPrefix(lcItem, "project:") {
			project = lcItem[8:]
		} else if strings.HasPrefix(lcItem, "-project:") {
			antiProjects = append(antiProjects, lcItem[9:])
		} else if len(item) > 2 && lcItem[0:1] == "+" {
			tags = append(tags, lcItem[1:])
		} else if len(item) > 2 && lcItem[0:1] == "-" {
			antiTags = append(antiTags, lcItem[1:])
		} else if IsValidPriority(item) {
			priority = item
		} else if item == IGNORE_CONTEXT_KEYWORD {
			ignoreContext = true
		} else {
			words = append(words, item)
		}
	}

	return CmdLine{
		Cmd:           cmd,
		IDs:           ids,
		Tags:          tags,
		AntiTags:      antiTags,
		Project:       project,
		AntiProjects:  antiProjects,
		Priority:      priority,
		Text:          strings.Join(words, " "),
		IgnoreContext: ignoreContext,
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

// make tags + projects are lowercase
func (task *Task) Normalise() {
	task.Project = strings.ToLower(task.Project)

	for i, tag := range task.Tags {
		task.Tags[i] = strings.ToLower(tag)
	}
}
