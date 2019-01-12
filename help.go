package dstask

import (
	"fmt"
	"os"
)

func Help(cmd string) {
	var helpStr string

	switch cmd {
	case CMD_NEXT:
		helpStr = `Usage: task next [filter] [--]
Usage: task [filter] [--]
Example: task +work +bug --

Display list of non-resolved tasks in the current context, most recent last,
optional filter. It is the default command, so "next" is unnecessary.

Bypass the current context with --.
`
	case CMD_ADD:
		helpStr = `Usage: task add [task summary] [--]
Example: task add Fix main web page 500 error +bug P1 project:website

Add a task, returning the git commit output which contains the task ID, used
later to reference the task.

Tags, project and priority can be added anywhere within the task summary.

Add -- to ignore the current context.

`

	case CMD_LOG:
	case CMD_START:
	case CMD_NOTE:
		fallthrough
	case CMD_NOTES:
	case CMD_STOP:
	case CMD_DONE:
	case CMD_DONE:
	case CMD_RESOLVE:
	case CMD_CONTEXT:
	case CMD_MODIFY:
	case CMD_EDIT:
	case CMD_UNDO:
	case CMD_SYNC:
	case CMD_GIT:
	case CMD_RESOLVED_TODAY:
	case CMD_RESOLVED_WEEK:
	case CMD_OPEN:
	case CMD_SHOW_PROJECTS:
	case CMD_IMPORT_TW:
	default:
		helpStr = `Usage: task <cmd> [id...] [task summary/filter]

Where [task summary] is text with tags/project/priority specified. Tags are
specified with + (or - for filtering) eg: +work. The project is specified with
a project:g prefix eg: project:dstask -- no quotes. Priorities run from P3
(low), P2 (default) to P1 (high) and P0 (critical). Cmd and IDs can be swapped.

run "task help <cmd>" for command specific help.

Add -- to ignore the current context.

Available commands:

add             : Add a task
log             : Log a task (already resolved)
start           : Change task status to active
note            : Append to or edit note for a task
stop            : Change task status to pending
resolve         : Resolve a task
context         : Set global context for task list and new tasks
modify          : Set attributes for a task
edit            : Edit task with text editor
undo            : Undo last action with git revert
pull            : Pull then push to git repository, automatic merge commit.
git             : Pass a command to git in the repository. Used for push/pull.
resolved-today  : Show tasks completed since midnight in current context
resolved-week   : Show tasks completed within the last week
show-projects   : List projects with completion status
open            : Open all URLs found in summary/annotations
import-tw       : Import tasks from taskwarrior via stdin
help            : Get help on any command or show this message

`
	}
	fmt.Fprintf(os.Stderr, helpStr)
	os.Exit(1)
}
