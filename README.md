# dstask

CI status of master: [![CircleCI](https://circleci.com/gh/naggie/dstask.svg?style=svg)](https://circleci.com/gh/naggie/dstask)

A personal task tracker designed to help you focus.

Features:

 * Powerful context system
 * Git powered sync/undo/resolve (passwordstore.org style) which means no need to set up a sync server, and sync between devices is easy!
 * Task listing won't break with long task text
 * `open` command -- open URLs found in specified task in the browser
 * `note` command -- edit a full markdown note for a task
 * zsh/bash completion for speed

Non-features:

 * Collaboration. This is a personal task tracker. Use another system for
   projects that involve multiple people. Note that it can still be beneficial
   to use dstask to track what you are working on in the context of a
   multi-person project tracked elsewhere.

Requirements:

* Git
* A 256-color capable terminal

<p align="center">
  <img src="https://github.com/naggie/dstask/raw/master/dstask.png">
</p>

# Installation

Note: This is beta software. There may be breaking changes before the 1.0 release.

1. Copy the executable (from the [releases page][1]) to somewhere in your path, named `dstask` and mark it executable. `/usr/local/bin/` is suggested.
1. Enable bash completions
1. Set up an alias in your bashrc: `alias task=dstask`, and source the relevant completion script

# Moving from Taskwarrior

Before installing dstask, you may want to export your taskwarrior database:

    task export > taskwarrior.json

After uninstalling taskwarrior and installing dstask, to import the tasks to dstask:

    dstask import-tw < taskwarrior.json


Commands and syntax are deliberately very similar to taskwarrior. Here are the exceptions:

  * The command is (nearly) always the first argument. Eg, `task eat some add bananas` won't work, but `task add eat some bananas` will. If there's an ID, it can proceed the command.
  * Priorities are added by the keywords `P0` `P1` `P2` `P3`. Lower number is more urgent. Default is `P2`. For example `task add eat some bananas P1`
  * Action is always the first argument. Eg, `task eat some add bananas` won't work, but `task add eat some bananas` will.
  * Contexts are defined on-the-fly, and are added to all new tasks if set. Use `--` to ignore current context in any command.

[1]: https://github.com/naggie/dstask/releases/latest

# Major things missing

There are a few things missing at the moment. That said I use dstask day to day and trust it with my work.

* Recurring tasks
* Subtask implementation (github issue style or otherwise)
* Deferring tasks
* Due dates
* Advanced reports
* Task dependencies

# Usage

```

Usage: dstask [id...] <cmd> [task summary/filter]

Where [task summary] is text with tags/project/priority specified. Tags are
specified with + (or - for filtering) eg: +work. The project is specified with
a project:g prefix eg: project:dstask -- no quotes. Priorities run from P3
(low), P2 (default) to P1 (high) and P0 (critical). Text can also be specified
for a substring search of description and notes.

Cmd and IDs can be swapped, multiple IDs can be specified for batch
operations.

run "task help <cmd>" for command specific help.

Add -- to ignore the current context. / can be used when adding tasks to note
any words after.

Available commands:

add             : Add a task
log             : Log a task (already resolved)
start           : Change task status to active
note            : Append to or edit note for a task
stop            : Change task status to pending
done            : Resolve a task
context         : Set global context for task list and new tasks
modify          : Set attributes for a task
edit            : Edit task with text editor
undo            : Undo last action with git revert
pull            : Pull then push to git repository, automatic merge commit.
git             : Pass a command to git in the repository. Used for push/pull.
resolved        : Show completed tasks
show-projects   : List projects with completion status
open            : Open all URLs found in summary/annotations
import-tw       : Import tasks from taskwarrior via stdin
help            : Get help on any command or show this message

```


# A note on performance

Currently I'm using dstask to manage thousands of tasks and the interface still
appears instant.

Dstask currently loads and parses every non-resolved task, each task being a
single file. This may sound wasteful, but it allows for a simple design and is
actually performant thanks to modern OS disk caches and SSDs.

If it starts to slow down as my number of non-resolved tasks increases, I'll
look into indexing and other optimisations such as archiving really old tasks.
