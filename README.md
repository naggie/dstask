[![Build Status](https://cloud.drone.io/api/badges/naggie/dstask/status.svg)](https://cloud.drone.io/naggie/dstask)

[![Go Report Card](https://goreportcard.com/badge/github.com/naggie/dstask)](https://goreportcard.com/report/github.com/naggie/dstask)

# dstask

A personal task tracker designed to help you focus. It is similar to
[taskwarrior](https://taskwarrior.org/) but uses git to synchronise instead of
a proprietary protocol.

Dstask is currently in beta -- the interface, data format and commands may
change before version 1.0. That said, it's unlikely that there will be a
breaking change as things are nearly finalised.

It's mature enough for daily use. I use dstask dozens of times a day, synchronised across 4 computers.

Features:

 * Powerful context system (automatically applies filter/tags to queries and new tasks)
 * **Git powered sync**/undo/resolve (passwordstore.org style) which means no need to set up a sync server, and sync between devices is easy!
 * Task listing won't break with long task text (unlike taskwarrior, currently)
 * `open` command -- **open URLs found in specified task** in the browser
 * `note` command -- edit a **full markdown note** for each task. Checklists are useful here.
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
  <img src="https://github.com/naggie/dstask/raw/master/etc/dstask.png">
</p>

# Installation

1. Copy the executable (from the [releases page][1]) to somewhere in your path, named `dstask` and mark it executable. `/usr/local/bin/` is suggested.
1. Enable bash completions by copying `.bash-completion.sh` into your home directory and sourcing it from your `.bashrc`. There's also a zsh completion script.
1. Set up an alias in your `.bashrc`: `alias task=dstask` or `alias t=dstask` to make task management slightly faster.
1. Create or clone a ~/.dstask git repository for the data, if you haven't already: `mkdir ~/.dstask && git -C ~/.dstask init`.

# Moving from Taskwarrior

Before installing dstask, you may want to export your taskwarrior database:

    task export > taskwarrior.json

After un-installing taskwarrior and installing dstask, to import the tasks to dstask:

    dstask import-tw < taskwarrior.json


Commands and syntax are deliberately very similar to taskwarrior. Here are the exceptions:

  * The command is (nearly) always the first argument. Eg, `task eat some add bananas` won't work, but `task add eat some bananas` will. If there's an ID, it can proceed the command but doesn't have to.
  * Priorities are added by the keywords `P0` `P1` `P2` `P3`. Lower number is more urgent. Default is `P2`. For example `task add eat some bananas P1`. The keyword can be anywhere after the command.
  * Action is always the first argument. Eg, `task eat some add bananas` won't work, but `task add eat some bananas` will.
  * Contexts are defined on-the-fly, and are added to all new tasks if set. Use `--` to ignore current context in any command.

[1]: https://github.com/naggie/dstask/releases/latest

# Future of dstask

There are a few things missing at the moment. That said I use dstask day to day and trust it with my work.

* Subtask/checklist implementation (github check list style)
* Task dependencies
* Recurring tasks
* Deferred/scheduled tasks with duration
* Due dates

After these features are implemented, I intend on adding CalDav integration. dstask should be able to act as a CalDav server with the following features:

* Display of recurring tasks
* Display of appointments/meetings
* Display of deadlines
* Display of resolved tasks (maybe, separate calendar)
* Possible creation of time based tasks from calendar

Running a caldav server would enable synchronisation with an iPhone/Android
phone, MacOS Calendar, Outlook calendar etc. It would provide similar
functionality to [timewarrior](https://timewarrior.net/).

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

run "dstask help <cmd>" for command specific help.

Add -- to ignore the current context. / can be used when adding tasks to note
any words after.

Available commands:

next              : Show most important tasks (priority, creation date -- truncated and default)
add               : Add a task
log               : Log a task (already resolved)
start             : Change task status to active
note              : Append to or edit note for a task
stop              : Change task status to pending
done              : Resolve a task
context           : Set global context for task list and new tasks (use "none" to set no context)
modify            : Set attributes for a task
edit              : Edit task with text editor
undo              : Undo last action with git revert
sync              : Pull then push to git repository, automatic merge commit.
open              : Open all URLs found in summary/annotations
git               : Pass a command to git in the repository. Used for push/pull.
show-projects     : List projects with completion status
show-tags         : List tags in use
show-active       : Show tasks that have been started
show-paused       : Show tasks that have been started then stopped
show-open         : Show all non-resolved tasks (without truncation)
show-resolved     : Show resolved tasks
show-unorganised  : Show untagged tasks with no projects (global context)
import-tw         : Import tasks from taskwarrior via stdin
help              : Get help on any command or show this message
version           : Show dstask version information
```

# Syntax


## Priority

| Symbol | Name      | Note                                                                 |
|--------|-----------|----------------------------------------------------------------------|
| `P0`   | Critical  | Must be resolved immediately. May appear in all contexts in future.  |
| `P1`   | High      |                                                                      |
| `P2`   | Normal    | Default priority                                                     |
| `P3`   | Low       | Shown at bottom and faded.                                           |


## Operators

| Symbol      | Syntax               | Description                                          | Example                                     |
|-------------|----------------------|------------------------------------------------------|---------------------------------------------|
| `+`         | `+<tag>`             | Include tag. Filter/context, or when adding task.    | `dstask add fix server +work`                 |
| `-`         | `-<tag>`             | Exclude tag. Filter/context only.                    | `dstask next -feature`                        |
| `--`        | `--`                 | Ignore context. When listing or adding tasks.        | `dstask --`, `task add -- +home do guttering` |
| `/`         | `/`                  | When adding a task, everything after will be a note. | `dstask add check out ipfs / https://ipfs.io` |
| `project:`  | `project:<project>`  | Set project. Filter/context, or when adding task.    | `dstask context project:dstask`               |
| `-project:` | `-project:<project>` | Exclude project, filter/context only.                | `dstask next -project:dstask -work`           |


# State

| State    | Description                                   |
|----------| ----------------------------------------------|
| Pending  | Tasks that have never been started            |
| Active   | Tasks that have been started                  |
| Paused   | Tasks that have been started but then stopped |
| Resolved | Tasks that have been done/close/completed     |


# Dealing with merge conflicts

Dstask is written in such a way that merge conflicts should not happen, unless
a task is edited independently on 2 or more machines without synchronising. In
practice this happens rarely; however when it does happen dstask will fail to
commit and warn you. You'll then need to go to the underlying `~/.dstask` git
repository and resolve manually before committing and running `dstask sync`. In
some rare cases the ID can conflict. This is something dstask will soon be
equipped to handle automatically when the `sync` command runs.

# A note on performance

Currently I'm using dstask to manage thousands of tasks and the interface still
appears instant.

Dstask currently loads and parses every non-resolved task, each task being a
single file. This may sound wasteful, but it allows git to track history
natively and is actually performant thanks to modern OS disk caches and SSDs.

If it starts to slow down as my number of non-resolved tasks increases, I'll
look into indexing and other optimisations such as archiving really old tasks.
I don't believe that this will be necessary, as the number of open tasks is
(hopefully) bounded.

# Issues

As you've probably noticed, I don't use the github issues. Currently I use
dstask itself to track dstask bugs in my personal dstask repository. I've left
the issues system enabled to allow people to report bugs or request features.
As soon as dstask is used by more than a handful of people, I'll probably
import the dstask issues to github.


# General tips

* Overwhelmed by tasks? Try focussing by prioritising (set priorities) or narrowing the context. The `show-tags` and `show-projects` commands are useful for creating a context.
* Use dstask to track things you might forget, rather than everything. SNR is important. Don't track tasks for the sake of it.
* Spend regular time reviewing tasks. You'll probably find some you've already resolved, and many you've forgotten. The `show-unorganised` command is good for this.
* Try to work through tasks from the top of the list. Dstask sorts by priority then creation date -- the most important tasks are at the top.
* Use `start`/`stop` to mark what you're genuinely working on right now; it makes resuming work faster. Paused tasks will be slightly highlighted, so you won't lose track of them. `show-paused` helps if they start to pile up.

# Database format

The format on disk stores the tasks in a directory according to the task
status, with each task stored under a yaml file with a UUID4 as the filename.
UUIDs are used to avoid conflicts when synchronising. The yaml schema is
defined by this Go struct:
https://github.com/naggie/dstask/blob/c00bc97c3f0132f1d291fdbe33dfb06e02ca6ef6/task.go#L18

This way only non-resolved tasks are actually loaded for most commands, so
performance is stable even with a large task history.

The ID presented to the user is simply a sequential ID. IDs are re-used when
tasks are resolved; tasks store their preferred ID for consistency across
different systems.

