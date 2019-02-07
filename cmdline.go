package dstask

// main task data structures

import (
	"strconv"
	"strings"
)

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
	IDsExhausted  bool
	// any words after the note operator: /
	Note string
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
	var notesModeActivated bool
	var notes []string
	var ignoreContext bool

	// something other than an ID has been parsed -- accept no more IDs
	var IDsExhausted bool

	for _, item := range args {
		lcItem := strings.ToLower(item)
		if !IDsExhausted && cmd == "" && StrSliceContains(ALL_CMDS, lcItem) {
			cmd = lcItem
			continue
		}

		if s, err := strconv.ParseInt(item, 10, 64); !IDsExhausted && err == nil {
			ids = append(ids, int(s))
			continue
		}

		IDsExhausted = true

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
		} else if item == NOTE_MODE_KEYWORD {
			notesModeActivated = true
		} else if notesModeActivated {
			notes = append(notes, item)
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
		Note:          strings.Join(notes, " "),
		IgnoreContext: ignoreContext,
		IDsExhausted:  IDsExhausted,
	}
}
