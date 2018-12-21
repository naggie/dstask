package dstask

import (
	"fmt"
	"os"
)

func Help() {
	fmt.Fprintf(os.Stderr, `
Usage: task add <filter> <summary>
Example: task add +work Fix CI building P2

	Add a task with a summary and context. Current global context is
	added.


Usage: task <id>

	Show detailed information about a task


Usage: task
Usage: task <filter>
Example: task P1

	List available tasks.


Usage: task context <context>
Usage: task context none
Example: task context project:dstask
Example: task context +work +bug

	Set (or clear) a global context for all queries and inserts.



Usage: taskwarrior export | task import-from-taskwarrior
	Import tasks from taskwarrior. Note that existing tasks will not be
	updated. This is to avoid dealing with conflicts.


Usage: task help

	Show this help dialog


Usage: task modify <id> <attributes...>

Usage: task edit <id>

Usage: task describe <id>


Usage: task week

	Show tasks completed in the last week, rolling


Usage: task day

	Show tasks completed since midnight


Usage: task projects

	List project status (percentage done, estimated completion time)
`)
	os.Exit(1)
}
