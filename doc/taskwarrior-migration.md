# Moving from Taskwarrior

## Importing data

Please see the [dstask-import tw instructions](dstask-import.md#taskwarrior) to import taskwarrior tasks.

## Using dstask

Commands and syntax are deliberately very similar to taskwarrior. Here are the exceptions:

  * The command is (nearly) always the first argument. Eg, `task eat some add bananas` won't work, but `task add eat some bananas` will. If there's an ID, it can proceed the command but doesn't have to.
  * Priorities are added by the keywords `P0` `P1` `P2` `P3`. Lower number is more urgent. Default is `P2`. For example `task add eat some bananas P1`. The keyword can be anywhere after the command.
  * Action is always the first argument. Eg, `task eat some add bananas` won't work, but `task add eat some bananas` will.
  * Contexts are defined on-the-fly, and are added to all new tasks if set. Use `--` to ignore current context in any command.

[1]: https://github.com/naggie/dstask/releases/latest

