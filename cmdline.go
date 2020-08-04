package dstask

// main task data structures

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// when referring to tasks by ID, NON_RESOLVED_STATUSES must be loaded exclusively --
// even if the filter is set to show issues that have only some statuses.
type CmdLine struct {
	Cmd           string
	IDs           []int
	Tags          []string
	AntiTags      []string
	Project       string
	AntiProjects  []string
	Priority      string
	Template      int
	Text          string
	UUID          string
	IgnoreContext bool
	IDsExhausted  bool
	// any words after the note operator: /
	Note string
}

// reconstruct args string
func (cmdLine CmdLine) String() string {
	var args []string

	for _, id := range cmdLine.IDs {
		args = append(args, strconv.Itoa(id))
	}

	for _, tag := range cmdLine.Tags {
		args = append(args, "+"+tag)
	}
	for _, tag := range cmdLine.AntiTags {
		args = append(args, "-"+tag)
	}

	if cmdLine.Project != "" {
		args = append(args, "project:"+cmdLine.Project)
	}

	for _, project := range cmdLine.AntiProjects {
		args = append(args, "-project:"+project)
	}

	if cmdLine.Priority != "" {
		args = append(args, cmdLine.Priority)
	}

	if cmdLine.Template > 0 {
		args = append(args, "template:"+string(cmdLine.Template))
	}

	if cmdLine.Text != "" {
		args = append(args, "\""+cmdLine.Text+"\"")
	}

	if cmdLine.UUID != "" {
		args = append(args, "\""+cmdLine.UUID+"\"")
	}

	return strings.Join(args, " ")
}

func (cmdLine CmdLine) PrintContextDescription() {
	if cmdLine.String() != "" {
		fmt.Printf("\033[33mActive context: %s\033[0m\n", cmdLine)
	}
}

// MustGetIdentifiers() determines if IDs or UUIDs are selected from the CmdLine.
// It should be used within functions that accept either IDs or UUIDs such as edit, note(s), or modify.
// The function returns three types:
// []interface{}: 	containins the list of task identifiers
// []string: 		the name of taskSet to load
//					NON_RESOLVED_STATUSES should be loaded when IDs are present
//					Tasks with STATUS_RESOLVED should be loaded only when UUIDs are present.
// error: 			The function will throw an error if no IDs are UUIDs are present in CmdLine.
//

func (cmdLine CmdLine) MustGetIdentifiers() ([]interface{}, []string, error) {
	if len(cmdLine.IDs) > 0 {
		identifiers := make([]interface{}, len(cmdLine.IDs))
		for i := range cmdLine.IDs {
			identifiers[i] = cmdLine.IDs[i]
		}
		return identifiers, NON_RESOLVED_STATUSES, nil
	} else if cmdLine.UUID != "" {
		identifiers := make([]interface{}, 1)
		identifiers[0] = cmdLine.UUID
		return identifiers, []string{STATUS_RESOLVED}, nil
	} else {
		return nil, nil, errors.New("MustGetIdentifiers() did not find any UUIDs or IDs in the command line.")
	}
}

// ParseCmdLine parses the raw command line typed by the user.
func ParseCmdLine(args ...string) CmdLine {
	var cmd string
	var ids []int
	var tags []string
	var antiTags []string
	var project string
	var antiProjects []string
	var priority string
	var template int
	var words []string
	var uuid string
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

		if item == IGNORE_CONTEXT_KEYWORD {
			// must be checked before negated tags, as -- is otherwise a valid tag
			ignoreContext = true
		} else if item == NOTE_MODE_KEYWORD {
			notesModeActivated = true
		} else if strings.HasPrefix(lcItem, "project:") {
			project = lcItem[8:]
		} else if strings.HasPrefix(lcItem, "+project:") {
			project = lcItem[9:]
		} else if strings.HasPrefix(lcItem, "-project:") {
			antiProjects = append(antiProjects, lcItem[9:])
		} else if strings.HasPrefix(lcItem, "template:") {
			if s, err := strconv.ParseInt(lcItem[9:], 10, 64); err == nil {
				template = int(s)
			}
		} else if len(item) > 1 && lcItem[0:1] == "+" {
			tags = append(tags, lcItem[1:])
		} else if len(item) > 1 && lcItem[0:1] == "-" {
			antiTags = append(antiTags, lcItem[1:])
		} else if IsValidPriority(item) {
			priority = item
		} else if strings.HasPrefix(lcItem, "uuid:") {
			uuid = lcItem[5:]
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
		Template:      template,
		Text:          strings.Join(words, " "),
		UUID:          uuid,
		Note:          strings.Join(notes, " "),
		IgnoreContext: ignoreContext,
		IDsExhausted:  IDsExhausted,
	}
}
