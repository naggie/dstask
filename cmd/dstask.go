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

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "next")
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			dstask.Help()
		}

		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		tl := dstask.ParseTaskLine(os.Args[2:]...)
		tl.MergeContext(context)
		task := dstask.Task{
			WritePending: true,
			Status:       dstask.STATUS_PENDING,
			Summary:      tl.Text,
			Tags:         tl.Tags,
			Project:      tl.Project,
			Priority:     tl.Priority,
		}
		task = ts.AddTask(task)
		ts.SaveToDisk("Added %s", task)

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
		ts.SaveToDisk("Started: %s", task)

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
		ts.SaveToDisk("Stopped %s", task)

	case "resolve":
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
			task.Comments = append(task.Comments, strings.Join(os.Args[3:], " "))
		}

		task.Resolved = time.Now() // could move to MustUpdateTask
		ts.MustUpdateTask(task)
		ts.SaveToDisk("Resolved %s", task)

	case "context":
		if len(os.Args) < 3 {
			dstask.Help()
		}

		if os.Args[2] == "none" {
			dstask.SaveContext()
		} else {
			dstask.SaveContext(os.Args[2:]...)
		}

	case "modify":
	case "edit":
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

	case "annotate":
		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		id, _ := strconv.Atoi(os.Args[2])
		task := ts.MustGetByID(id)

		if len(os.Args) == 3 {
			task.Description = string(dstask.MustEditBytes([]byte(task.Description), "md"))
		} else {
			task.Description += "\n" + strings.Join(os.Args[3:], " ")
		}

		ts.MustUpdateTask(task)
		ts.SaveToDisk("Describe %s", task)

	case "undo":
		dstask.MustRunGitCmd("revert", "--no-edit", "HEAD")

	case "git":
		dstask.MustRunGitCmd(os.Args[2:]...)

	case "day":
	case "week":

	case "import-tw":
		ts := dstask.LoadTaskSetFromDisk(dstask.ALL_STATUSES)
		ts.ImportFromTaskwarrior()
		ts.SaveToDisk("Import from taskwarrior")

	case "projects":

	case "help":
		dstask.Help()

	default:
		var args []string
		// next, or just a filter which is effectively an alias for next
		if os.Args[1] == "next" {
			args = os.Args[2:]
		} else {
			args = os.Args[1:]
		}

		ts := dstask.LoadTaskSetFromDisk(dstask.NON_RESOLVED_STATUSES)
		tl := dstask.ParseTaskLine(args...)
		ts.Filter(context)
		ts.Filter(tl)
		ts.SortTaskList()
		if context.String() != "" {
			fmt.Printf("\n\n\033[33mActive context: %s\033[0m\n", context)
		} else {
			fmt.Printf("\n\n\n")
		}
		ts.Display()
	}
}
