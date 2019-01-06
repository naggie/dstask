# dstask

A personal task tracker designed to help you focus.

Features:

 * Powerful context system
 * Git powered sync/undo/resolve (passwordstore.org style)
 * Task listing won't break with long task text
 * `open` command -- open URLs found in specified task in the browser

Non-features:

 * Collaboration. This is a personal task tracker. Use another system for
   projects that involve multiple people. Note that it can still be beneficial
   to use dstask to track what you are working on in the context of a
   multi-person project tracked elsewhere.

<p align="center">
  <img src="https://github.com/naggie/dstask/raw/master/dstask.png">
</p>

# Installation

Note: This is beta software. There may be breaking changes before the 1.0 release.

1. Copy the executable (from the [releases page][1]) to somewhere in your path, named `dstask` and mark it executable. `/usr/local/bin/` is suggested.
1. Enable bash completions
1. Set up an alias in your bashrc: `alias task=dstask`

# Moving from Taskwarrior

Before installing dstask, you may want to export your taskwarrior database:

    task export > taskwarrior.json

After uninstalling taskwarrior and installing dstask, to import the tasks to dstask:

    task import-tw < taskwarrior.json


Commands and syntax are deliberately very similar to taskwarrior. Here are the exceptions:

  * The command is (nearly) always the first argument. Eg, `task eat some add bananas` won't work, but `task add eat some bananas` will. If there's an ID, it can proceed the command.
  * Priorities are added by the keywords `P0` `P1` `P2` `P3`. Lower number is more urgent. Default is `P2`. For example `task add eat some bananas P1`
  * Action is always the first argument. Eg, `task eat some add bananas` won't work, but `task add eat some bananas` will.
  * Contexts are defined on-the-fly, and are added to all new tasks if set. Use `--` to ignore current context in any command.

[1]: https://github.com/naggie/dstask/releases/latest

# Usage

```

Usage: task <cmd> [id...] [task summary]

Where [task summary] is text with tags/project/priority specified. Tags are
specified with + (or - for filtering) eg: +work. The project is specified with
a "project:" prefix eg: project:dstask -- no quotes. Priorities run from P3
(low), P2 (default) to P1 (high) and P0 (critical). Cmd and IDs can be swapped.

run "task help <cmd>" for command specific help.

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

```
