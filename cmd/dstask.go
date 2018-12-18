package main

import (
	"github.com/naggie/dstask"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "next")
	}

	switch os.Args[1] {
	case "next":
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		ts.SortTaskList()
		ts.Display()

	case "add":
		if len(os.Args) < 3 {
			dstask.Help()
		}

		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		tl := dstask.ParseTaskLine(os.Args[2:])
		ts.AddTask(dstask.Task{
			WritePending: true,
			Status:       dstask.STATUS_PENDING,
			Summary:      tl.Text,
			Tags:         tl.Tags,
			Project:      tl.Project,
			Priority:     tl.Priority,
		})
		ts.SaveToDisk("Added: " + tl.Text)

	case "start":
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
		ts.SaveToDisk("Started: " + task.Summary)

	case "stop":
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
		ts.SaveToDisk("Stopped: " + task.Summary)

	case "done":
	case "context":
	case "modify":
	case "edit":
	case "describe":
	case "projects":
	case "day":
	case "week":
	case "import-tw":
		ts := dstask.LoadTaskSetFromDisk(dstask.ALL_STATUSES)
		ts.ImportFromTaskwarrior()
		ts.SaveToDisk("Import from taskwarrior")

	case "git":
		dstask.MustRunGitCmd(os.Args[2:]...)

	case "undo":
		dstask.MustRunGitCmd("revert", "HEAD")

	case "help":
		dstask.Help()

	default:
		dstask.Help()
	}
}
