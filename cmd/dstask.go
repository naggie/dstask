package main

import (
	"fmt"
	"github.com/naggie/dstask"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	context := dstask.LoadContext()
	cmdLine := dstask.ParseCmdLine(os.Args[1:])

	switch os.Args[1] {
	case CMD_NEXT:
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

	case AMD_ADD:
		if len(os.Args) < 3 {
			dstask.Help()
		}

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

	case CMD_START:
		if len(os.Args) != 3 {
			dstask.Help()
		}

		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		idStr, _ := strconv.Atoi(os.Args[2])
		task := ts.MustGetByID(idStr)

		// TODO probably allow more here
		if task.Status != dstask.STATUS_PENDING {
			dstask.ExitFail("That task is not pending")
		}

		task.Status = dstask.STATUS_ACTIVE
		ts.MustUpdateTask(task)
		ts.SaveToDisk("Started: %s", task)

	case CMD_STOP:
		if len(os.Args) != 3 {
			dstask.Help()
		}

		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		idStr, _ := strconv.Atoi(os.Args[2])
		task := ts.MustGetByID(idStr)

		if task.Status != dstask.STATUS_ACTIVE {
			dstask.ExitFail("That task is not yet started")
		}

		task.Status = dstask.STATUS_PENDING
		ts.MustUpdateTask(task)
		ts.SaveToDisk("Stopped %s", task)

	case CMD_RESOLVE:
		if len(os.Args) < 3 {
			dstask.Help()
		}

		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		idStr, _ := strconv.Atoi(os.Args[2])
		task := ts.MustGetByID(idStr)

		// TODO definitely move to MustUpdateTask
		if task.Status == dstask.STATUS_RESOLVED {
			dstask.ExitFail("That task is already resolved")
		}

		task.Status = dstask.STATUS_RESOLVED

		if len(os.Args) > 3 {
			task.Notes += "\n" + strings.Join(os.Args[3:], " ")
		}

		task.Resolved = time.Now() // could move to MustUpdateTask
		ts.MustUpdateTask(task)
		ts.SaveToDisk("Resolved %s", task)

	case CMD_CONTEXT:
		if len(os.Args) < 3 {
			dstask.Help()
		}

		if os.Args[2] == "none" {
			dstask.SaveContext()
		} else {
			dstask.SaveContext(os.Args[2:]...)
		}

	case CMD_MODIFY:
	case CMD_EDIT:
		if len(os.Args) != 3 {
			dstask.Help()
		}

		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		id, _ := strconv.Atoi(os.Args[2])
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

	case CMD_ANNOTATE:
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		id, _ := strconv.Atoi(os.Args[2])
		task := ts.MustGetByID(id)

		if len(os.Args) == 3 {
			task.Notes = string(dstask.MustEditBytes([]byte(task.Notes), "md"))
		} else {
			task.Notes += "\n" + strings.Join(os.Args[3:], " ")
		}

		ts.MustUpdateTask(task)
		ts.SaveToDisk("Describe %s", task)

	case CMD_UNDO:
		dstask.MustRunGitCmd("revert", "--no-edit", "HEAD")

	case GMD_GIT:
		dstask.MustRunGitCmd(os.Args[2:]...)

	case CMD_DAY:
	case CMD_WEEK:

	case CMD_IMPORT_TW:
		ts := dstask.LoadTaskSetFromDisk(dstask.ALL_STATUSES)
		ts.ImportFromTaskwarrior()
		ts.SaveToDisk("Import from taskwarrior")

	case CMD_PROJECTS:

	case CMD_HELP:
		dstask.Help()

	}
}
