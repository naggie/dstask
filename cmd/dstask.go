package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mvdan/xurls"
	"github.com/naggie/dstask"
	"gopkg.in/yaml.v2"
)

func main() {
	// Sets globals: GIT_REPO, STATE_FILE, IDS_FILE
	dstask.ParseConfig()
	dstask.EnsureRepoExists(dstask.GIT_REPO)
	// Load state for getting and setting context
	state := dstask.LoadState()
	context := state.Context
	cmdLine := dstask.ParseCmdLine(os.Args[1:]...)

	if cmdLine.IgnoreContext {
		context = dstask.CmdLine{}
	}

	switch cmdLine.Cmd {
	// Empty string is interpreted as CMD_NEXT
	case "", dstask.CMD_NEXT:
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)
		ts.Filter(context)
		ts.Filter(cmdLine)
		ts.FilterOutStatus(dstask.STATUS_TEMPLATE)
		ts.SortByPriority()
		context.PrintContextDescription()
		ts.DisplayByNext(true)
		ts.DisplayCriticalTaskWarning()

	case dstask.CMD_SHOW_OPEN:
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)
		ts.Filter(context)
		ts.Filter(cmdLine)
		ts.FilterOutStatus(dstask.STATUS_TEMPLATE)
		ts.SortByPriority()
		context.PrintContextDescription()
		ts.DisplayByNext(false)
		ts.DisplayCriticalTaskWarning()

	case dstask.CMD_ADD:
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)

		if cmdLine.Template > 0 {
			var taskSummary string
			tt, err := ts.MustGetByID(cmdLine.Template)
			if err != nil {
				dstask.ExitFail("In CMD_ADD: Unable to create Template from task %v: %v", cmdLine.Template, err)
			}

			context.PrintContextDescription()
			cmdLine.MergeContext(context)

			if cmdLine.Text != "" {
				taskSummary = cmdLine.Text
			} else {
				taskSummary = tt.Summary
			}

			// create task from template task tt
			task := dstask.Task{
				WritePending: true,
				Status:       dstask.STATUS_PENDING,
				Summary:      taskSummary,
				Tags:         tt.Tags,
				Project:      tt.Project,
				Priority:     tt.Priority,
				Notes:        tt.Notes,
			}

			// Modify the task with any tags/projects/antiProjects/priorities in cmdLine
			task.Modify(cmdLine)

			task = ts.LoadTask(task)
			ts.SavePendingChanges()
			dstask.MustGitCommit("Added %s", task)
			if tt.Status != dstask.STATUS_TEMPLATE {
				// Insert Text Statement to inform user of real Templates
				fmt.Print("\nYou've copied an open task!\nTo learn more about creating templates enter 'dstask help template'\n\n")
			}
		} else if cmdLine.Text != "" {
			context.PrintContextDescription()
			cmdLine.MergeContext(context)
			task := dstask.Task{
				WritePending: true,
				Status:       dstask.STATUS_PENDING,
				Summary:      cmdLine.Text,
				Tags:         cmdLine.Tags,
				Project:      cmdLine.Project,
				Priority:     cmdLine.Priority,
				Notes:        cmdLine.Note,
			}
			task = ts.LoadTask(task)
			ts.SavePendingChanges()
			dstask.MustGitCommit("Added %s", task)
		}

	case dstask.CMD_RM, dstask.CMD_REMOVE:
		if len(cmdLine.IDs) < 1 {
			dstask.ExitFail("%s", "missing argument: id")
		}
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)
		for _, id := range cmdLine.IDs {
			task, err := ts.MustGetByID(id)
			if err != nil {
				dstask.ExitFail("In CMD_RM: %v", err)
			}
			// Mark our task for deletion
			task.Deleted = true

			// MustUpdateTask validates and normalises our task object
			ts.MustUpdateTask(task)
			ts.SavePendingChanges()
			dstask.MustGitCommit("Removed: %s", task)
		}

	case dstask.CMD_TEMPLATE:
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)

		if len(cmdLine.IDs) > 0 {
			for _, id := range cmdLine.IDs {
				task, err := ts.MustGetByID(id)
				if err != nil {
					dstask.ExitFail("In CMD_TEMPLATE: %v", err)
				}
				task.Status = dstask.STATUS_TEMPLATE

				ts.MustUpdateTask(task)
				ts.SavePendingChanges()
				dstask.MustGitCommit("Changed %s to Template", task)
			}
		} else if cmdLine.Text != "" {
			context.PrintContextDescription()
			cmdLine.MergeContext(context)
			task := dstask.Task{
				WritePending: true,
				Status:       dstask.STATUS_TEMPLATE,
				Summary:      cmdLine.Text,
				Tags:         cmdLine.Tags,
				Project:      cmdLine.Project,
				Priority:     cmdLine.Priority,
				Notes:        cmdLine.Note,
			}
			task = ts.LoadTask(task)
			ts.SavePendingChanges()
			dstask.MustGitCommit("Created Template %s", task)
		}

	case dstask.CMD_LOG:
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)

		if cmdLine.Text != "" {
			context.PrintContextDescription()
			cmdLine.MergeContext(context)
			task := dstask.Task{
				WritePending: true,
				Status:       dstask.STATUS_RESOLVED,
				Summary:      cmdLine.Text,
				Tags:         cmdLine.Tags,
				Project:      cmdLine.Project,
				Priority:     cmdLine.Priority,
				Resolved:     time.Now(),
			}
			task = ts.LoadTask(task)
			ts.SavePendingChanges()
			dstask.MustGitCommit("Logged %s", task)
		}

	case dstask.CMD_START:
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)
		if len(cmdLine.IDs) > 0 {
			// start given tasks by IDs
			for _, id := range cmdLine.IDs {
				task, err := ts.MustGetByID(id)
				if err != nil {
					dstask.ExitFail("In CMD_START: %v", err)
				}
				task.Status = dstask.STATUS_ACTIVE
				if cmdLine.Text != "" {
					task.Notes += "\n" + cmdLine.Text
				}
				ts.MustUpdateTask(task)

				ts.SavePendingChanges()
				dstask.MustGitCommit("Started %s", task)

				if task.Notes != "" {
					fmt.Printf("\nNotes on task %d:\n\033[38;5;245m%s\033[0m\n\n", task.ID, task.Notes)
				}
			}
		} else if cmdLine.Text != "" {
			// create a new task that is already active (started)
			cmdLine.MergeContext(context)
			task := dstask.Task{
				WritePending: true,
				Status:       dstask.STATUS_ACTIVE,
				Summary:      cmdLine.Text,
				Tags:         cmdLine.Tags,
				Project:      cmdLine.Project,
				Priority:     cmdLine.Priority,
				Notes:        cmdLine.Note,
			}
			task = ts.LoadTask(task)
			ts.SavePendingChanges()
			dstask.MustGitCommit("Added and started %s", task)
		}

	case dstask.CMD_STOP:
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)
		for _, id := range cmdLine.IDs {
			task, err := ts.MustGetByID(id)
			if err != nil {
				dstask.ExitFail("In CMD_STOP: %v", err)
			}
			task.Status = dstask.STATUS_PAUSED
			if cmdLine.Text != "" {
				task.Notes += "\n" + cmdLine.Text
			}
			ts.MustUpdateTask(task)
			ts.SavePendingChanges()
			dstask.MustGitCommit("Stopped %s", task)
		}

	case dstask.CMD_DONE, dstask.CMD_RESOLVE:
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)
		for _, id := range cmdLine.IDs {
			task, err := ts.MustGetByID(id)
			if err != nil {
				dstask.ExitFail("In CMD_DONE/CMD_RESOLVE: %v", err)
			}

			task.Status = dstask.STATUS_RESOLVED
			if cmdLine.Text != "" {
				task.Notes += "\n" + cmdLine.Text
			}
			ts.MustUpdateTask(task)
			ts.SavePendingChanges()
			dstask.MustGitCommit("Resolved %s", task)
		}

	case dstask.CMD_CONTEXT:
		if len(os.Args) < 3 {
			fmt.Printf("Current context: %s\n", context)
		} else if os.Args[2] == "none" {
			if err := state.SetContext(dstask.CmdLine{}); err != nil {
				dstask.ExitFail(err.Error())
			}
		} else {
			if err := state.SetContext(cmdLine); err != nil {
				dstask.ExitFail(err.Error())
			}
		}
		state.Save()

	case dstask.CMD_MODIFY:

		identifiers, taskSet, nil := cmdLine.MustGetIdentifiers()

		ts := dstask.LoadTasksFromDisk(taskSet)

		if len(identifiers) == 0 {
			ts.Filter(context)
			dstask.ConfirmOrAbort("No IDs specified. Apply to all %d tasks in current context?", len(ts.Tasks()))

			for _, task := range ts.Tasks() {
				task.Modify(cmdLine)
				ts.MustUpdateTask(task)
				ts.SavePendingChanges()
				dstask.MustGitCommit("Modified %s", task)
			}
			return
		}

		for _, id := range identifiers {
			task, err := ts.MustGetTask(id)
			if err != nil {
				dstask.ExitFail("In CMD_MODIFY: %v", err)
			}

			task.Modify(cmdLine)
			ts.MustUpdateTask(task)
			ts.SavePendingChanges()
			dstask.MustGitCommit("Modified %s", task)
		}

	case dstask.CMD_EDIT:

		identifiers, taskSet, err := cmdLine.MustGetIdentifiers()
		if err != nil {
			dstask.ExitFail("In CMD_EDIT: %v", err)
		}
		editResolved := taskSet[0] == dstask.STATUS_RESOLVED

		ts := dstask.LoadTasksFromDisk(taskSet)

		for _, id := range identifiers {
			task, err := ts.MustGetTask(id)
			if err != nil {
				dstask.ExitFail("In CMD_EDIT: %v", err)
			}

			// Hide ID
			task.ID = 0

			data, err := yaml.Marshal(&task)
			if err != nil {
				dstask.ExitFail(fmt.Sprintf("Failed to marshal task %s: %v\n", task, err))
			}

			data = dstask.MustEditBytes(data, "yml")

			err = yaml.Unmarshal(data, &task)
			if err != nil {
				// TODO reattempt mechanism
				dstask.ExitFail(fmt.Sprintf("Failed to unmarshal yml: %v\n", err))
			}

			// Re-add ID
			if !editResolved {
				task.ID, _ = id.(int)
			}
			// TODO edit MustUpdateTask to erase resolved date if necessary
			ts.MustUpdateTask(task)
			ts.SavePendingChanges()
			dstask.MustGitCommit("Edited %s", task)

			if editResolved {
				// If the editing a resolved Task and status is changing, load NON_RESOLVED_STATUSES to assign ID
				// After LoadTasksFromDisk() is called, the new ID of the resolved task will be displayed
				if task.Status != dstask.STATUS_RESOLVED {
					ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)
					task, err := ts.MustGetByUUID(cmdLine.UUID)
					if err != nil {
						dstask.ExitFail("In CMD_EDIT: Failed to Assign ID to pending task: %v", err)
					}
					fmt.Printf("This task has now been assigned the ID: %v\n", task.ID)
				}
			}
		}

	case dstask.CMD_NOTE, dstask.CMD_NOTES:

		// If stdout is not a TTY, we simply write markdown notes to stdout
		openEditor := dstask.IsTTY()

		identifiers, taskSet, nil := cmdLine.MustGetIdentifiers()

		ts := dstask.LoadTasksFromDisk(taskSet)

		for _, identifier := range identifiers {
			task, err := ts.MustGetTask(identifier)
			if err != nil {
				dstask.ExitFail("In CMD_NOTE(S): Unable to retrieve task: %v", err)
			}
			err = task.AddNote(cmdLine, openEditor)
			if err != nil {
				dstask.ExitFail("In CMD_NOTE(S): %v", err)
			}

			if openEditor {
				ts.MustUpdateTask(task)
				ts.SavePendingChanges()
				dstask.MustGitCommit("Edit note %s", task)
			}
		}

	case dstask.CMD_UNDO:
		var err error
		n := 1
		if len(os.Args) == 3 {
			n, err = strconv.Atoi(os.Args[2])
			if err != nil {
				dstask.Help(dstask.CMD_UNDO)
			}
		}

		dstask.MustRunGitCmd("revert", "--no-gpg-sign", "--no-edit", "HEAD~"+strconv.Itoa(n)+"..")

	case dstask.CMD_SYNC:
		dstask.Sync()

	case dstask.CMD_GIT:
		dstask.MustRunGitCmd(os.Args[2:]...)

	case dstask.CMD_SHOW_ACTIVE:
		context.PrintContextDescription()
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)
		ts.Filter(context)
		ts.Filter(cmdLine)
		ts.FilterByStatus(dstask.STATUS_ACTIVE)
		ts.SortByPriority()
		ts.DisplayByNext(true)

	case dstask.CMD_SHOW_PAUSED:
		context.PrintContextDescription()
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)
		ts.Filter(context)
		ts.Filter(cmdLine)
		ts.FilterByStatus(dstask.STATUS_PAUSED)
		ts.SortByPriority()
		ts.DisplayByNext(true)

	case dstask.CMD_OPEN:
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)
		for _, id := range cmdLine.IDs {
			task, err := ts.MustGetByID(id)
			if err != nil {
				dstask.ExitFail("In CMD_OPEN: %v", err)
			}
			urls := xurls.Relaxed.FindAllString(task.Summary+" "+task.Notes, -1)

			if len(urls) == 0 {
				dstask.ExitFail("No URLs found in task %v", task.ID)
			}

			for _, url := range urls {
				dstask.MustOpenBrowser(url)
			}
		}

	case dstask.CMD_IMPORT_TW:
		ts := dstask.LoadTasksFromDisk(dstask.ALL_STATUSES)
		ts.ImportFromTaskwarrior()
		ts.SavePendingChanges()
		dstask.MustGitCommit("Import from taskwarrior")

	case dstask.CMD_SHOW_PROJECTS:
		context.PrintContextDescription()
		ts := dstask.LoadTasksFromDisk(dstask.ALL_STATUSES)
		cmdLine.MergeContext(context)
		ts.Filter(context)
		ts.DisplayProjects()

	case dstask.CMD_SHOW_TAGS:
		context.PrintContextDescription()
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)
		cmdLine.MergeContext(context)
		ts.Filter(context)
		for tag := range ts.GetTags() {
			fmt.Println(tag)
		}

	case dstask.CMD_SHOW_TEMPLATES:
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)
		ts.Filter(context)
		ts.Filter(cmdLine)
		ts.FilterByStatus(dstask.STATUS_TEMPLATE)
		ts.SortByPriority()
		ts.DisplayByNext(false)
		context.PrintContextDescription()

	case dstask.CMD_SHOW_RESOLVED:
		ts := dstask.LoadTasksFromDisk(dstask.ALL_STATUSES)
		ts.Filter(context)
		ts.Filter(cmdLine)
		ts.FilterByStatus(dstask.STATUS_RESOLVED)
		ts.SortByResolved()
		ts.DisplayByWeek()
		context.PrintContextDescription()

	case dstask.CMD_SHOW_UNORGANISED:
		ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)
		ts.Filter(cmdLine)
		ts.FilterUnorganised()
		ts.DisplayByNext(true)

	case dstask.CMD_HELP:
		if len(os.Args) > 2 {
			dstask.Help(os.Args[2])
		} else {
			dstask.Help("")
		}

	case dstask.CMD_VERSION:
		fmt.Printf(
			"Version: %s\nGit commit: %s\nBuild date: %s\n",
			dstask.VERSION,
			dstask.GIT_COMMIT,
			dstask.BUILD_DATE,
		)

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

		// args are dstask _completions <user command line>
		// parse command line as normal to set rules
		cmdLine := dstask.ParseCmdLine(originalArgs...)

		// no command specified, default given
		if !cmdLine.IDsExhausted || cmdLine.Cmd == dstask.CMD_HELP || cmdLine.Cmd == "" {
			for _, cmd := range dstask.ALL_CMDS {
				if !strings.HasPrefix(cmd, "_") {
					completions = append(completions, cmd)
				}
			}
		}

		//uuid
		if dstask.StrSliceContains([]string{
			dstask.CMD_EDIT,
			dstask.CMD_NOTE,
			dstask.CMD_NOTES,
			dstask.CMD_MODIFY,
		}, cmdLine.Cmd) && len(cmdLine.IDs) < 1 {
			// This will complete with a space after the colon.
			ts := dstask.LoadTasksFromDisk([]string{dstask.STATUS_RESOLVED})
			for _, task := range ts.Tasks() {
				completions = append(completions, "uuid:"+task.UUID)
			}
		}

		if dstask.StrSliceContains([]string{
			"",
			dstask.CMD_NEXT,
			dstask.CMD_ADD,
			dstask.CMD_REMOVE,
			dstask.CMD_LOG,
			dstask.CMD_START,
			dstask.CMD_STOP,
			dstask.CMD_DONE,
			dstask.CMD_RESOLVE,
			dstask.CMD_CONTEXT,
			dstask.CMD_MODIFY,
			dstask.CMD_SHOW_NEXT,
			dstask.CMD_SHOW_PROJECTS,
			dstask.CMD_SHOW_ACTIVE,
			dstask.CMD_SHOW_PAUSED,
			dstask.CMD_SHOW_OPEN,
			dstask.CMD_SHOW_RESOLVED,
			dstask.CMD_SHOW_TEMPLATES,
		}, cmdLine.Cmd) {
			ts := dstask.LoadTasksFromDisk(dstask.NON_RESOLVED_STATUSES)
			// limit completions to available context, but not if the user is
			// trying to change context, context ignore is on, or modify
			// command is being completed
			if !cmdLine.IgnoreContext &&
				cmdLine.Cmd != dstask.CMD_CONTEXT &&
				cmdLine.Cmd != dstask.CMD_MODIFY {
				ts.Filter(context)
			}

			// templates
			if cmdLine.Cmd == dstask.CMD_ADD {
				for _, task := range ts.Tasks() {
					if task.Status == dstask.STATUS_TEMPLATE {
						completions = append(completions, "template:"+strconv.Itoa(task.ID))
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
