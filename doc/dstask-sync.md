# dstask-sync

dstask-sync is a tool to synchronize between external services and dstask.
At this point it only supports syncing from GitHub.
The goal currently is to have tasks in dstask that represent tasks in GitHub,
such that a dstask-based workflow (tracking, managing and prioritizing tasks)
can take into account work that is defined in GitHub, although the goal is not to "replace" GitHub.

Specifically:
* The sync is one-way (from Github to dstask).
* We only sync key properties (summary etc), not GH comments or dstask notes
* You are expected to close issues in GitHub and then sync to get the task closed in dstask.

## configuration

Create a file `$HOME/.tasksync.toml` with one or more github sections, like this:

```
[[github]]
token = "<Github API token>"
user = "<Github org/user>"
repo = "<Github repository>"   
get_closed = true             # get closed tickets in addition to open ones?
assignee = ""                 # if set, only import tickets that have this assignee
milestone = ""                # if set, select only tickets that have this milestone
label = ""                    # if set, only select tickets that have this label
template = "default"          # must be set to a valid task file in ~/.dstask/templates-github/<filename>
```

Note:
* selection by Github project not supported yet
* you may have multiple sections that import an overlapping set of tasks, with conflicting directives.
  we resolve this as "last directive wins". This is simply the result of executing each section in sequence,
  such that tasks re-imported by later sections may overwrite tasks that were imported by earlier sections.

## templates

Put a file like this in "~/.dstask/templates-github/default.yaml"

TODO fill in

These variables can be used:
TODO
we need tags and project also

## Properties mapping in detail

As a reminder, here are how issues/PR's and tasks, are modeled on Github and within dstask respectively.

| Github properties | dstask properties |
-----------------------------------------
| org/user          |                                            |
| repo              |                                            |
| number            | uuid                                       |
| title             | summary                                    |
| state open/closed | status (pending, active, paused, resolved) |
| pr vs issue       |                                            |
| author            |                                            |
| created timestamp | created timestamp                          |
|                   | resolved timestamp                         |
|                   | due timestamp                              |
| labels            | tags                                       |
| projects          | project                                    |
| milestones        |                                            |
| review state      |                                            |
| reviewers         | delegatedto (?)                            |
| assignees         | subtasks                                   |
| link issue        | dependencies                               |
| comments          | notes                                      |

The mapping works as follows:

* UUID: auto-generated based on org/user, repo and ticket number
* summary: template gets expanded with data from Github
* status:
  - if open in GH, default to pending but if task already exists with status active or paused, we leave that intact
  - if closed in GH, set to resolved.
* notes: default empty.  local/pre-existing notes honored.
* created: created timestamp from GitHub.
* priority, tags, project: template gets expanded with data from Github

Other fields (subtasks, due timestamp, etc) are left default/empty.

