# Dstask database format

The default database location is `~/.dstask/`, but can be configured by the
environment variable `DSTASK_GIT_REPO`.

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

TODO elaborate with examples.
