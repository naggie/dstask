<p align="center">
<img align="center" src="etc/icon.png" alt="icon" height="64" />
</p>

<h1 align="center">dstask</h1>

<p align="center">
<i> Single binary terminal-based todo manager: git-based sync + markdown notes for each task.  </i>
</p>

<p align="center">
<a href="https://cloud.drone.io/naggie/dstask"><img src="https://cloud.drone.io/api/badges/naggie/dstask/status.svg" /></a>
<a href="https://goreportcard.com/report/github.com/naggie/dstask"><img src="https://goreportcard.com/badge/github.com/naggie/dstask" /></a>
<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/license-MIT-blue.svg" /></a>
</p>

<br>
<br>
<br>


Dstask is a personal task tracker designed to help you focus. It is similar to
[taskwarrior](https://taskwarrior.org/) but uses git to synchronise instead of
a proprietary protocol.

Dstask is mature enough for daily use. I use dstask dozens of times a day,
synchronised across 4 computers.

Features:

 * Powerful context system (automatically applies filter/tags to queries and new tasks)
 * **Git powered sync**/undo/resolve (passwordstore.org style) which means no need to set up a sync server, and sync between devices is easy!
 * Task listing won't break with long task text ([unlike taskwarrior, currently](https://github.com/GothenburgBitFactory/taskwarrior/issues/2023))
 * `note` command -- edit a **full markdown note** for each task. Checklists are useful here.
 * `open` command -- **open URLs found in specified task** (including notes) in the browser
 * zsh/bash completion for speed
 * A single statically-linked binary

Non-features:

 * Collaboration. This is a personal task tracker. Use another system for
   projects that involve multiple people. Note that it can still be beneficial
   to use dstask to track what you are working on in the context of a
   multi-person project tracked elsewhere.

Requirements:

* Git
* A 256-color capable terminal

# Screenshots

<p align="center">
  <img src="https://github.com/naggie/dstask/raw/master/etc/dstask.png">
  <em>Next command (default when no command is specified)</em>
</p>

<p align="center">
  <img src="https://github.com/naggie/dstask/raw/master/etc/show-resolved.png">
  <em>Show-resolved command to review completed tasks by week. Useful for meetings.</em>
</p>

<p align="center">
  <img src="https://github.com/naggie/dstask/raw/master/etc/edit.png">
  <em>Editing a task with $EDITOR (which happens to be vim)</em>
</p>

<p align="center">
  <img src="https://github.com/naggie/dstask/raw/master/etc/add.png">
  <em>Adding a task</em>
</p>

<p align="center">
  <img src="https://github.com/naggie/dstask/raw/master/etc/sync.png">
  <em>Sync command (which uses git)</em>
</p>

# Installation

1. Copy the executable (from the [releases page][1]) to somewhere in your path, named `dstask` and mark it executable. `/usr/local/bin/` is suggested.
1. Enable bash completions by copying `.bash-completion.sh` into your home directory and sourcing it from your `.bashrc`. There's also a zsh completion script.
1. Set up an alias in your `.bashrc`: `alias task=dstask` or `alias t=dstask` to make task management slightly faster.
1. Create or clone a ~/.dstask git repository for the data, if you haven't already: `mkdir ~/.dstask && git -C ~/.dstask init`.


There is also an unofficial
[Nix](https://nixos.org/nixos/packages.html?attr=dstask&channel=nixpkgs-unstable&query=dstask)
and [Arch AUR](https://aur.archlinux.org/packages/dstask/) package!

# Moving from Taskwarrior

See [etc/MIGRATION.md](etc/MIGRATION.md)

# Future of dstask

See [etc/FUTURE.md](etc/FUTURE.md)

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
template          : Add a task template
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

| Symbol      | Syntax               | Description                                          | Example                                       |
|-------------|----------------------|------------------------------------------------------|-----------------------------------------------|
| `+`         | `+<tag>`             | Include tag. Filter/context, or when adding task.    | `dstask add fix server +work`                 |
| `-`         | `-<tag>`             | Exclude tag. Filter/context only.                    | `dstask next -feature`                        |
| `--`        | `--`                 | Ignore context. When listing or adding tasks.        | `dstask --`, `task add -- +home do guttering` |
| `/`         | `/`                  | When adding a task, everything after will be a note. | `dstask add check out ipfs / https://ipfs.io` |
| `project:`  | `project:<project>`  | Set project. Filter/context, or when adding task.    | `dstask context project:dstask`               |
| `-project:` | `-project:<project>` | Exclude project, filter/context only.                | `dstask next -project:dstask -work`           |
| `template:` | `template:<id>`      | Base new task on a template.                         | `dstask add template:24`                      |


## State

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

# Performance

See [etc/PERFORMANCE.md](etc/PERFORMANCE.md)

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
* Keep a [github-style check list](https://help.github.com/en/articles/about-task-lists) in the markdown note of complex or procedural tasks
* Failing to get started working? Start with the smallest task
* Record only required tasks. Track ideas separately, else your task list will grow unboundedly! I keep an `ideas.md` for various projects for this reason.

# Database

See [etc/DATABASE_FORMAT.md](etc/DATABASE_FORMAT.md)

The default database location is `~/.dstask/`, but can be configured by the
environment variable `DSTASK_GIT_REPO`.

# Alternatives

Alternatives listed must be capable of running in the terminal.

* [TaskLite](https://github.com/ad-si/TaskLite) -- The CLI task manager for power users, written in Haskell
* [Taskwarrior](https://taskwarrior.org/) -- the closest analogue
* [Taskbook](https://github.com/klaussinani/taskbook) -- board metaphor, note support
* [todo.txt-cli](https://github.com/todotxt/todo.txt-cli)


# FAQ

> Does dstask encrypt tasks?

Encryption is not a design goal of dstask. If you want to have your remote
repository encrypted, you may consider
[git-remote-gcrypt](https://spwhitton.name/tech/code/git-remote-gcrypt/) or
[git-crypt](https://github.com/AGWA/git-crypt). Note that dstask has not been
tested with these tools, nor can any claims be made about the security of the
tools themselves.

> Is it possible to modify more than one task at once with a filter?

Yes.

1. Set a context:
2. Run a modify command without and ID
3. Hit y to confirm to modify all tasks in context

This means it's natural to review the tasks that would be modified before
modifying by listing all tasks in the current context first, instead of
potentially operating blindly by matching tags or numbers.

You can also specify multiple task numbers at one time, as with any other command.
