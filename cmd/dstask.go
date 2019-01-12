package main

import (
	"fmt"
	"github.com/mvdan/xurls"
	"github.com/naggie/dstask"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"time"
)

func main() {
	context := dstask.LoadContext()
	cmdLine := dstask.ParseCmdLine(os.Args[1:]...)

	if cmdLine.IgnoreContext {
		context = dstask.CmdLine{}
	}

	switch cmdLine.Cmd {
	case "":
		// default command is CMD_NEXT if not specified
		fallthrough
	case dstask.CMD_NEXT:
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		ts.Filter(context)
		ts.Filter(cmdLine)
		ts.SortTaskList()
		if context.String() != "" {
			fmt.Printf("\n\n\033[33mActive context: %s\033[0m\n", context)
		} else {
			fmt.Printf("\n\n\n")
		}
		ts.Display()

	case dstask.CMD_ADD:
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)

		if len(cmdLine.Text) != 0 {
			cmdLine.MergeContext(context)
			task := dstask.Task{
				WritePending: true,
				Status:       dstask.STATUS_PENDING,
				Summary:      cmdLine.Text,
				Tags:         cmdLine.Tags,
				Project:      cmdLine.Project,
				Priority:     cmdLine.Priority,
			}
			task = ts.AddTask(task)
			ts.SaveToDisk("Added %s", task)
		}

	case dstask.CMD_LOG:
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)

		if len(cmdLine.Text) != 0 {
			cmdLine.MergeContext(context)
			task := dstask.Task{
				WritePending: true,
				Status:       dstask.STATUS_RESOLVED,
				Summary:      cmdLine.Text,
				Tags:         cmdLine.Tags,
				Project:      cmdLine.Project,
				Priority:     cmdLine.Priority,
			}
			task = ts.AddTask(task)
			ts.SaveToDisk("Logged %s", task)
		}

	case dstask.CMD_START:
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		if len(cmdLine.IDs) > 0 {
			// start given tasks by IDs
			for _, id := range cmdLine.IDs {
				task := ts.MustGetByID(id)
				task.Status = dstask.STATUS_ACTIVE
				ts.MustUpdateTask(task)
				ts.SaveToDisk("Started %s", task)
			}
		} else if len(cmdLine.Text) != 0 {
			// create a new task that is already active (started)
			cmdLine.MergeContext(context)
			task := dstask.Task{
				WritePending: true,
				Status:       dstask.STATUS_ACTIVE,
				Summary:      cmdLine.Text,
				Tags:         cmdLine.Tags,
				Project:      cmdLine.Project,
				Priority:     cmdLine.Priority,
			}
			task = ts.AddTask(task)
			ts.SaveToDisk("Added and started %s", task)
		}

	case dstask.CMD_STOP:
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		for _, id := range cmdLine.IDs {
			task := ts.MustGetByID(id)
			task.Status = dstask.STATUS_PENDING
			ts.MustUpdateTask(task)
			ts.SaveToDisk("Stopped %s", task)
		}

	case dstask.CMD_DONE:
		fallthrough
	case dstask.CMD_RESOLVE:
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		for _, id := range cmdLine.IDs {
			task := ts.MustGetByID(id)
			task.Status = dstask.STATUS_RESOLVED

			if cmdLine.Text != "" {
				task.Notes += "\n" + cmdLine.Text
			}
			ts.MustUpdateTask(task)
			ts.SaveToDisk("Resolved %s", task)
		}

	case dstask.CMD_CONTEXT:
		if os.Args[2] == "none" {
			dstask.SaveContext(dstask.CmdLine{})
		} else {
			dstask.SaveContext(cmdLine)
		}

	case dstask.CMD_MODIFY:
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		for _, id := range cmdLine.IDs {
			task := ts.MustGetByID(id)

			for _, tag := range cmdLine.Tags {
				if !dstask.StrSliceContains(task.Tags, tag) {
					task.Tags = append(task.Tags, tag)
				}
			}

			for i, tag := range task.Tags {
				if dstask.StrSliceContains(cmdLine.AntiTags, tag) {
					// delete item
					task.Tags = append(task.Tags[:i], task.Tags[i+1:]...)
				}
			}

			if cmdLine.Project != "" {
				task.Project = cmdLine.Project
			}

			if dstask.StrSliceContains(cmdLine.AntiProjects, task.Project) {
				task.Project = ""
			}

			if cmdLine.Priority != "" {
				task.Priority = cmdLine.Priority
			}

			ts.MustUpdateTask(task)
			ts.SaveToDisk("Modified %s", task)
		}

	case dstask.CMD_EDIT:
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		for _, id := range cmdLine.IDs {
			task := ts.MustGetByID(id)

			// hide ID
			task.ID = 0

			data, err := yaml.Marshal(&task)
			if err != nil {
				// TODO present error to user, specific error message is important
				dstask.ExitFail("Failed to marshal task %s", task)
			}

			data = dstask.MustEditBytes(data, "yml")

			err = yaml.Unmarshal(data, &task)
			if err != nil {
				// TODO present error to user, specific error message is important
				// TODO reattempt mechansim
				dstask.ExitFail("Failed to unmarshal yml")
			}

			// re-add ID
			task.ID = id

			ts.MustUpdateTask(task)
			ts.SaveToDisk("Edited %s", task)
		}

	case dstask.CMD_NOTES:
		fallthrough
	case dstask.CMD_NOTE:
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		for _, id := range cmdLine.IDs {
			task := ts.MustGetByID(id)
			if cmdLine.Text == "" {
				task.Notes = string(dstask.MustEditBytes([]byte(task.Notes), "md"))
			} else {
				if task.Notes == "" {
					task.Notes = cmdLine.Text
				} else {
					task.Notes += "\n" + cmdLine.Text
				}
			}

			ts.MustUpdateTask(task)
			ts.SaveToDisk("Edit note %s", task)
		}

	case dstask.CMD_UNDO:
		dstask.MustRunGitCmd("revert", "--no-edit", "HEAD")

	case dstask.CMD_SYNC:
		dstask.MustRunGitCmd("pull", "--no-edit", "--commit", "origin", "master")
		dstask.MustRunGitCmd("push", "origin", "master")

	case dstask.CMD_GIT:
		dstask.MustRunGitCmd(os.Args[2:]...)

	case dstask.CMD_RESOLVED_TODAY:
		t := time.Now()
		year, month, day := t.Date()
		bod := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
		ts := dstask.LoadTaskSetFromDisk(dstask.ALL_STATUSES)
		ts.Filter(context)
		ts.FilterResolvedSince(bod)
		ts.Display()

	case dstask.CMD_RESOLVED_WEEK:
		ts := dstask.LoadTaskSetFromDisk(dstask.ALL_STATUSES)
		ts.Filter(context)
		ts.FilterResolvedSince(time.Now().AddDate(0, 0, -7))
		ts.Display()

	case dstask.CMD_OPEN:
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		for _, id := range cmdLine.IDs {
			task := ts.MustGetByID(id)
			urls := xurls.Relaxed().FindAllString(task.Summary+" "+task.Notes, -1)

			if len(urls) == 0 {
				dstask.ExitFail("No URLs found in task %v", task.ID)
			}

			for _, url := range urls {
				dstask.MustOpenBrowser(url)
			}
		}

	case dstask.CMD_IMPORT_TW:
		ts := dstask.LoadTaskSetFromDisk(dstask.ALL_STATUSES)
		ts.ImportFromTaskwarrior()
		ts.SaveToDisk("Import from taskwarrior")

	case dstask.CMD_SHOW_PROJECTS:
		ts := dstask.LoadTaskSetFromDisk(dstask.ALL_STATUSES)
		cmdLine.MergeContext(context)
		ts.Filter(context)
		projects := ts.GetProjects()
		for name := range projects {
			fmt.Printf("%+v\n", projects[name])
		}

	case dstask.CMD_HELP:
		dstask.Help()

	case dstask.CMD_COMPLETIONS:
		// given the entire user's command line arguments as the arguments for
		// this cmd, suggest possible candidates for the last arg.
		// see the relevant shell completion bindings in this repository for
		// integration. Note there are various idiosyncrasies with bash
		// involving arg separation.
		var completions []string
		var originalArgs []string
		var prefix string

		if len(os.Args) > 3 {
			originalArgs = os.Args[3:]
		}

		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)

		// args are dstask _completions <user command line>
		// parse command line as normal to set rules
		cmdLine := dstask.ParseCmdLine(originalArgs...)

		// limit completions to available context, but not if the user is
		// trying to change context or context ignore is on
		if !cmdLine.IgnoreContext && cmdLine.Cmd != dstask.CMD_CONTEXT {
			ts.Filter(context)
		}

		// no command specified, default given
		if !cmdLine.IDsExhausted {
			for _, cmd := range dstask.ALL_CMDS {
				if !strings.HasPrefix(cmd, "_") {
					completions = append(completions, cmd)
				}
			}
		}

		// priorities
		completions = append(completions, dstask.PRIORITY_CRITICAL)
		completions = append(completions, dstask.PRIORITY_HIGH)
		completions = append(completions, dstask.PRIORITY_NORMAL)
		completions = append(completions, dstask.PRIORITY_LOW)

		// projects
		for project := range ts.GetProjects() {
			completions = append(completions, "project:"+project)
			completions = append(completions, "-project:"+project)
		}

		// tags
		for tag := range ts.GetTags() {
			completions = append(completions, "+"+tag)
			completions = append(completions, "-"+tag)
		}

		if len(originalArgs) > 0 {
			prefix = originalArgs[len(originalArgs)-1]
		}

		for _, completion := range completions {
			if strings.HasPrefix(completion, prefix) && !dstask.StrSliceContains(originalArgs, completion) {
				fmt.Println(completion)
			}
		}
	}
}
