package dstask

import (
	"fmt"
	"os"
)

func Help() {
	fmt.Fprintf(os.Stderr, `
Usage: task <cmd> [args...]

run "task help <cmd>" for command specific help.

Available commands:

add          : Add a task
start        : Change task status to active
stop         : Change task status to pending
resolve      : Resolve a task
comment      : Add a comment to a task
context      : Set global context for task list and new tasks
modify       : set attributes for a task
edit         : Edit task with text editor
describe     : Add a detailed description in text editor
undo         : Undo last action with git revert
git          : Pass a command to git in the repository. Used for push/pull.

day          : Show tasks completed since midnight in current context
week         : Show tasks completed within the last week
projects     : List projects with completion status

import-tw    : Import tasks from taskwarrior

help         : Get help on any command or show this message


`)
	os.Exit(1)
}
