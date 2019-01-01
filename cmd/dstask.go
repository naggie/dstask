package main

import (
	"fmt"
	"github.com/mvdan/xurls"
	"github.com/naggie/dstask"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

func main() {
	context := dstask.LoadContext()
	cmdLine := dstask.ParseCmdLine(os.Args[1:]...)

	switch cmdLine.Cmd {
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

	case dstask.CMD_START:
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		for _, id := range cmdLine.IDs {
			task := ts.MustGetByID(id)
			task.Status = dstask.STATUS_ACTIVE
			ts.MustUpdateTask(task)
			ts.SaveToDisk("Started: %s", task)
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

	case dstask.CMD_ANNOTATE:
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
			ts.SaveToDisk("Annotate %s", task)
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
		ts.FilterResolvedSince(time.Now().AddDate(0,0,-7))
		ts.Display()

	case dstask.CMD_OPEN:
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		for _, id := range cmdLine.IDs {
			task := ts.MustGetByID(id)
			url := xurls.Relaxed().FindString(task.Summary + " " + task.Notes)

			if url == "" {
				dstask.ExitFail("No URL found in task %v", task.ID)
			}

			dstask.MustOpenBrowser(url)
		}

	case dstask.CMD_IMPORT_TW:
		ts := dstask.LoadTaskSetFromDisk(dstask.ALL_STATUSES)
		ts.ImportFromTaskwarrior()
		ts.SaveToDisk("Import from taskwarrior")

	case dstask.CMD_PROJECTS:

	case dstask.CMD_HELP:
		dstask.Help()
	}
}
