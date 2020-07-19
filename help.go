package dstask

import (
	"fmt"
	"os"
)

func Help(cmd string) {
	var helpStr string

	switch cmd {
	case CMD_NEXT:
		helpStr = `Usage: dstask next [filter] [--]
Usage: dstask [filter] [--]
Example: dstask +work +bug --

Display list of non-resolved tasks in the current context, most recent last,
optional filter. It is the default command, so "next" is unnecessary.

Bypass the current context with --.
`
	case CMD_ADD:
		helpStr = `Usage: dstask add [template:<id>] [task summary] [--]
Example: dstask add Fix main web page 500 error +bug P1 project:website

Add a task, returning the git commit output which contains the task ID, used
later to reference the task.

Tags, project and priority can be added anywhere within the task summary.

Add -- to ignore the current context. / can be used when adding tasks to note
any words after.

A copy of an existing task can be made by including "template:<id>". See
"dstask help template" for more information on templates.

`
	case CMD_TEMPLATE:
		helpStr = `Usage dstask template <id> [task summary] [--]
Example: dstask template Fix main web page 500 error +bug P1 project:website
Example: dstask template 34 project:

If valid task ID is supplied, a copy of the task is created as a template. If
no ID is given, a new task template is created.

Tags, project and priority can be added anywhere within the task summary.

Add -- to ignore the current context. / can be used when adding tasks to note
any words after

Template tasks are not displayed with "show-open" or "show-next" commands.
Their intent is to act as a readily available task template for commonly used
or repeated tasks.

To create a new task from a template use the command:
"dstask add template:<id> [task summary] [--]"
The template task <id> remains unchanged, but a new task is created as a copy
with any modifications made in the task summary.

Github-style task lists (checklists) are recommended for templates, useful for
performing procedures. Example:

- [ ] buy bananas
- [ ] eat bananas
- [ ] make coffee

`
	case CMD_RM, CMD_REMOVE:
		helpStr = `Usage: dstask remove <id...>
Example: dstask 15 remove

Remove a task.

The task is deleted from the filesystem, and the change is committed.

`

	case CMD_LOG:
		helpStr = `Usage: dstask log [task summary] [--]
Example: dstask log Fix main web page 500 error +bug P1 project:website

Add an immediately resolved task. Syntax identical to add command.

Tags, project and priority can be added anywhere within the task summary.

Add -- to ignore the current context.

`
	case CMD_START:
		helpStr = `Usage: dstask <id...> start
Usage: dstask start [task summary] [--]
Example: dstask 15 start
Example: dstask start Fix main web page 500 error +bug P1 project:website

Mark a task as active, meaning you're currently at work on the task.

Alternatively, "start" can add a task and start it immediately with the same
syntax is the "add" command.  Tags, project and priority can be added anywhere
within the task summary.

Add -- to ignore the current context.
`
	case CMD_NOTE:
		fallthrough
	case CMD_NOTES:
		helpStr = `Usage: dstask note <id>
Usage: dstask note <id> <text>
Example task 13 note problem is faulty hardware

Edit or append text to the markdown notes attached to a particular task.
`
	case CMD_STOP:
		helpStr = `Usage: dstask <id...> stop [text]
Example: dstask 15 stop
Example: dstask 15 stop replaced some hardware

Set a task as inactive, meaning you've stopped work on the task. Optional text
may be added, which will be appended to the note.
`
	case CMD_RESOLVE:
		fallthrough
	case CMD_DONE:
		helpStr = `Usage: dstask <id...> done [text]
Example: dstask 15 done
Example: dstask 15 done replaced some hardware

Resolve a task. Optional text may be added, which will be appended to the note.
`
	case CMD_CONTEXT:
		helpStr = `Usage: dstask context <filter>
Example: dstask context +work -bug
Example: dstask context none

Set a global filter consisting of a project, tags or antitags. Subsequent new
tasks and most commands will then have this filter applied automatically.

For example, if you were to run "task add fix the webserver," the given task
would then have the tag "work" applied automatically.

To reset to no context, run: dstask context none
`
	case CMD_MODIFY:
		helpStr = `Usage: dstask <id...> modify <filter>
Usage: dstask modify <filter>
Example: dstask 34 modify -work +home project:workbench -project:website

Modify the attributes of the given tasks, specified by ID. If no ID is given,
the operation will be performed to all tasks in the current context subject to
confirmation.

Modifiable attributes: tags, project and priority.
`
	case CMD_EDIT:
		helpStr = `Usage: dstask <id...> edit

Edit a task in your text editor.
`
	case CMD_UNDO:
		helpStr = `Usage: dstask undo
Usage: dstask undo <n>

Undo the last <n> commits on the repository. Default is 1. Use

	dstask git log

To see commit history. For more complicated history manipulation it may be best
to revert/rebase/merge on the dstask repository itself. The dstask repository
is at ~/.dstask by default.
`
	case CMD_SYNC:
		helpStr = `Usage: dstask sync

Synchronise with the remote git server. Runs git pull then git push. If there
are conflicts that cannot be automatically resolved, it is necessary to
manually resolve them in  ~/.dstask or with the "task git" command.
`
	case CMD_GIT:
		helpStr = `Usage: dstask git <args...>
Example: dstask git status

Run the given git command inside ~/.dstask
`
	case CMD_SHOW_RESOLVED:
		helpStr = `Usage: dstask resolved

Show a report of last 1000 resolved tasks.
`
	case CMD_SHOW_TEMPLATES:
		helpStr = `Usage: dtask show-templates [filter] [--]

Show a report of stored template tasks with an optional filter.

Bypass the current context with --`
	case CMD_OPEN:
		helpStr = `Usage: dstask <id...> open

Open all URLs found within the task summary and notes. If you commonly have
dozens of tabs open to later action, convert them into tasks to open later with
this command.
`
	case CMD_SHOW_PROJECTS:
		helpStr = `Usage: dstask show-projects

Show a breakdown of projects with progress information
`
	case CMD_IMPORT_TW:
		helpStr = `Usage: cat export.json | task import-tw

Import tasks from a taskwarrior json dump. The "task export" taskwarrior
command can be used for this.
`
	default:
		helpStr = `Usage: dstask [id...] <cmd> [task summary/filter]

Where [task summary] is text with tags/project/priority specified. Tags are
specified with + (or - for filtering) eg: +work. The project is specified with
a project:g prefix eg: project:dstask -- no quotes. Priorities run from P3
(low), P2 (default) to P1 (high) and P0 (critical). Text can also be specified
for a substring search of description and notes.

Cmd and IDs can be swapped, multiple IDs can be specified for batch
operations.

run "dstask help <cmd>" for command specific help.

Add -- to ignore the current context. / can be used when adding tasks to note
any words after.

Available commands:

next              : Show most important tasks (priority, creation date -- truncated and default)
add               : Add a task
template          : Add a task template
log               : Log a task (already resolved)
start             : Change task status to active
note              : Append to or edit note for a task
stop              : Change task status to pending
done              : Resolve a task
context           : Set global context for task list and new tasks (use "none" to set no context)
modify            : Change task attributes specified on command line
edit              : Edit task with text editor
undo              : Undo last n commits
sync              : Pull then push to git repository, automatic merge commit.
open              : Open all URLs found in summary/annotations
git               : Pass a command to git in the repository. Used for push/pull.
remove            : Remove a task (use to remove tasks added by mistake)
show-projects     : List projects with completion status
show-tags         : List tags in use
show-active       : Show tasks that have been started
show-paused       : Show tasks that have been started then stopped
show-open         : Show all non-resolved tasks (without truncation)
show-resolved     : Show resolved tasks
show-templates    : Show task templates
show-unorganised  : Show untagged tasks with no projects (global context)
import-tw         : Import tasks from taskwarrior via stdin
help              : Get help on any command or show this message
version           : Show dstask version information

Task table key:

`
	}
	fmt.Fprintf(os.Stderr, helpStr)

	colourPrintln(0, FG_PRIORITY_CRITICAL, BG_DEFAULT_2, "Critical priority")
	colourPrintln(0, FG_PRIORITY_HIGH, BG_DEFAULT_2, "High priority")
	colourPrintln(0, FG_DEFAULT, BG_DEFAULT_1, "Normal priority")
	colourPrintln(0, FG_PRIORITY_LOW, BG_DEFAULT_2, "Low priority")
	colourPrintln(0, FG_ACTIVE, BG_ACTIVE, "Active")
	colourPrintln(0, FG_DEFAULT, BG_PAUSED, "Paused")

	os.Exit(0)
}

func colourPrintln(mode, fg, bg int, line string) {
	line = FixStr(line, 25)
	fmt.Fprintf(os.Stderr, "\033[%d;38;5;%d;48;5;%dm%s\033[0m\n", mode, fg, bg, line)
}
