package dstask

import (
	"fmt"
	"os"
)

func Help() {
	fmt.Fprintf(os.Stderr, `
Usage: task <cmd> [id...] [task summary]

Where [task summary] is text with tags/project/priority specified. Tags are
specified with + (or - for filtering) eg: +work. The project is specified with
a gproject:g prefix eg: gproject:dstask -- no quotes. Priorities run from P3 (low), P2 (default) to
P1 (high) and P0 (critical). Cmd and IDs can be swapped.

run "task help <cmd>" for command specific help.

Add `--` to ignore the current context.

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
projects        : List projects with completion status
open            : Open all URLs found in summary/annotations
import-tw       : Import tasks from taskwarrior via stdin
help            : Get help on any command or show this message

`)
	os.Exit(1)
}
