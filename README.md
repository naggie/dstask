# dstask

A personal task tracker designed to help you focus.

Features:

 * Powerful context system
 * Git powered sync/undo/resolve (passwordstore.org style)
 * Task listing won't break with long task text

Non-features:

 * Collaboration. This is a personal task tracker. Use another system for
   projects that involve multiple people. Note that it can still be beneficial
   to use dstask to track what you are working on in the context of a
   multi-person project tracked elsewhere.

<p align="center">
  <img width="460" height="300" src="https://github.com/naggie/dstask/raw/master/dstask.png">
</p>

# Installation

1. Copy the executable (from the [releases page][1]) to somewhere in your path, named `dstask` and mark it executable. `/usr/local/bin/` is suggested.
1. Enable bash completions
1. Set up an alias in your bashrc: `alias task=dstask`

# Moving from Taskwarrior

Before installing dstask, you may want to export your taskwarrior database:

    task export > taskwarrior.json

After uninstalling taskwarrior and installing dstask, to import the tasks to dstask:

    task import-tw < taskwarrior.json


Commands and syntax are deliberately very similar to taskwarrior. Here are the exceptions:

  * Action is always the first argument. Eg, `task eat some add bananas` won't work, but `task add eat some bananas` will.
  * Priorities are added by the keywords `P1` `P2` `P3` `P4`. Lower number is more urgent. Default is `P3`. For example `task add eat some bananas P2`
  * Contexts are defined on-the-fly, and are added to all new tasks if set.


[1]: https://github.com/naggie/dstask/releases/latest
