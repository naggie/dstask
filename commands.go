package dstask

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mvdan/xurls"
	yaml "gopkg.in/yaml.v2"
)

// CommandAdd ...
func CommandAdd(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}
	if cmdLine.Template > 0 {
		var taskSummary string
		tt := ts.MustGetByID(cmdLine.Template)
		ctx.PrintContextDescription()
		cmdLine.MergeContext(ctx)

		if cmdLine.Text != "" {
			taskSummary = cmdLine.Text
		} else {
			taskSummary = tt.Summary
		}

		// create task from template task tt
		task := Task{
			WritePending: true,
			Status:       STATUS_PENDING,
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
		MustGitCommit(repoPath, "Added %s", task)
		if tt.Status != STATUS_TEMPLATE {
			// Insert Text Statement to inform user of real Templates
			fmt.Print("\nYou've copied an open task!\nTo learn more about creating templates enter 'dstask help template'\n\n")
		}
	} else if cmdLine.Text != "" {
		ctx.PrintContextDescription()
		cmdLine.MergeContext(ctx)
		task := Task{
			WritePending: true,
			Status:       STATUS_PENDING,
			Summary:      cmdLine.Text,
			Tags:         cmdLine.Tags,
			Project:      cmdLine.Project,
			Priority:     cmdLine.Priority,
			Notes:        cmdLine.Note,
		}
		task = ts.LoadTask(task)
		ts.SavePendingChanges()
		MustGitCommit(repoPath, "Added %s", task)

	}
	return nil
}

// CommandContext ...
func CommandContext(repoPath, idsFilePath, stateFilePath string, state State, ctx, cmdLine CmdLine) error {
	if len(os.Args) < 3 {
		fmt.Printf("Current context: %s\n", ctx)
	} else if os.Args[2] == "none" {
		if err := state.SetContext(CmdLine{}); err != nil {
			ExitFail(err.Error())
		}
	} else {
		if err := state.SetContext(cmdLine); err != nil {
			ExitFail(err.Error())
		}
	}
	state.Save(stateFilePath)
	return nil
}

// CommandDone ...
func CommandDone(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}
	for _, id := range cmdLine.IDs {
		task := ts.MustGetByID(id)
		task.Status = STATUS_RESOLVED
		if cmdLine.Text != "" {
			task.Notes += "\n" + cmdLine.Text
		}
		ts.MustUpdateTask(task)
		ts.SavePendingChanges()
		MustGitCommit(repoPath, "Resolved %s", task)
	}
	return nil
}

// CommandEdit ...
func CommandEdit(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}
	for _, id := range cmdLine.IDs {
		task := ts.MustGetByID(id)

		// hide ID
		task.ID = 0

		data, err := yaml.Marshal(&task)
		if err != nil {
			// TODO present error to user, specific error message is important
			ExitFail("Failed to marshal task %s", task)
		}

		data = MustEditBytes(data, "yml")

		err = yaml.Unmarshal(data, &task)
		if err != nil {
			// TODO present error to user, specific error message is important
			// TODO reattempt mechanism
			ExitFail("Failed to unmarshal yml")
		}

		// re-add ID
		task.ID = id

		ts.MustUpdateTask(task)
		ts.SavePendingChanges()
		MustGitCommit(repoPath, "Edited %s", task)
	}
	return nil
}

// CommandHelp ...
func CommandHelp(args []string) {
	if len(os.Args) > 2 {
		Help(os.Args[2])
	} else {
		Help("")
	}
}

// CommandImportTW ...
func CommandImportTW(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(ALL_STATUSES...),
	)
	if err != nil {
		return err
	}
	ts.ImportFromTaskwarrior()
	ts.SavePendingChanges()
	MustGitCommit(repoPath, "Import from taskwarrior")
	return nil
}

// CommandLog ...
func CommandLog(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}

	if cmdLine.Text != "" {
		ctx.PrintContextDescription()
		cmdLine.MergeContext(ctx)
		task := Task{
			WritePending: true,
			Status:       STATUS_RESOLVED,
			Summary:      cmdLine.Text,
			Tags:         cmdLine.Tags,
			Project:      cmdLine.Project,
			Priority:     cmdLine.Priority,
			Resolved:     time.Now(),
		}
		task = ts.LoadTask(task)
		ts.SavePendingChanges()
		MustGitCommit(repoPath, "Logged %s", task)
	}

	return nil
}
func CommandModify(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}

	if len(cmdLine.IDs) == 0 {
		ts.Filter(ctx)
		ConfirmOrAbort("No IDs specified. Apply to all %d tasks in current ctx?", len(ts.Tasks()))

		for _, task := range ts.Tasks() {
			task.Modify(cmdLine)
			ts.MustUpdateTask(task)
			ts.SavePendingChanges()
			MustGitCommit(repoPath, "Modified %s", task)
		}
		return nil
	}

	for _, id := range cmdLine.IDs {
		task := ts.MustGetByID(id)
		task.Modify(cmdLine)
		ts.MustUpdateTask(task)
		ts.SavePendingChanges()
		MustGitCommit(repoPath, "Modified %s", task)
	}

	return nil
}

// CommandNext prints the unresolved tasks associated with the current context.
// This is the default command.
func CommandNext(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithoutStatuses(STATUS_TEMPLATE),
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}
	ts.Filter(ctx)
	ts.Filter(cmdLine)
	ts.SortByPriority()
	ctx.PrintContextDescription()
	ts.DisplayByNext(true)
	ts.DisplayCriticalTaskWarning()

	return nil
}

// CommandNote ...
func CommandNote(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}

	// If stdout is not a TTY, we simply write markdown notes to stdout
	openEditor := IsTTY()

	for _, id := range cmdLine.IDs {
		task := ts.MustGetByID(id)
		if openEditor {
			if cmdLine.Text == "" {
				task.Notes = string(MustEditBytes([]byte(task.Notes), "md"))
			} else {
				if task.Notes == "" {
					task.Notes = cmdLine.Text
				} else {
					task.Notes += "\n" + cmdLine.Text
				}
			}
			ts.MustUpdateTask(task)
			ts.SavePendingChanges()
			MustGitCommit(repoPath, "Edit note %s", task)
		} else {
			if err := WriteStdout([]byte(task.Notes)); err != nil {
				ExitFail("Could not write to stdout: %v", err)
			}
		}
	}
	return nil
}

// CommandOpen ...
func CommandOpen(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}
	for _, id := range cmdLine.IDs {
		task := ts.MustGetByID(id)
		urls := xurls.Relaxed.FindAllString(task.Summary+" "+task.Notes, -1)

		if len(urls) == 0 {
			return fmt.Errorf("No URLs found in task %v", task.ID)
		}

		for _, url := range urls {
			MustOpenBrowser(url)
		}
	}

	return nil
}

// CommandRemove ...
func CommandRemove(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	if len(cmdLine.IDs) < 1 {
		return errors.New("missing argument: id")
	}
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}
	for _, id := range cmdLine.IDs {
		task := ts.MustGetByID(id)

		// Mark our task for deletion
		task.Deleted = true

		// MustUpdateTask validates and normalises our task object
		ts.MustUpdateTask(task)
		ts.SavePendingChanges()
		MustGitCommit(repoPath, "Removed: %s", task)
	}
	return nil
}

// CommandShowActive ...
func CommandShowActive(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ctx.PrintContextDescription()
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}
	ts.Filter(ctx)
	ts.Filter(cmdLine)
	ts.FilterByStatus(STATUS_ACTIVE)
	ts.SortByPriority()
	ts.DisplayByNext(true)

	return nil
}

// CommandShowProjects ...
func CommandShowProjects(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ctx.PrintContextDescription()
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(ALL_STATUSES...),
	)
	if err != nil {
		return err
	}
	cmdLine.MergeContext(ctx)
	ts.Filter(ctx)
	ts.DisplayProjects()
	return nil
}

// CommandShowOpen ...
func CommandShowOpen(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithoutStatuses(STATUS_TEMPLATE),
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}
	ts.Filter(ctx)
	ts.Filter(cmdLine)
	ts.SortByPriority()
	ctx.PrintContextDescription()
	ts.DisplayByNext(false)
	ts.DisplayCriticalTaskWarning()
	return nil
}
func CommandShowPaused(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ctx.PrintContextDescription()
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}
	ts.Filter(ctx)
	ts.Filter(cmdLine)
	ts.FilterByStatus(STATUS_PAUSED)
	ts.SortByPriority()
	ts.DisplayByNext(true)
	return nil
}

// CommandShowResolved ...
func CommandShowResolved(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(ALL_STATUSES...),
	)
	if err != nil {
		return err
	}
	ts.Filter(ctx)
	ts.Filter(cmdLine)
	ts.FilterByStatus(STATUS_RESOLVED)
	ts.SortByResolved()
	ts.DisplayByWeek()
	ctx.PrintContextDescription()
	return nil
}

// CommandShowTags ...
func CommandShowTags(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ctx.PrintContextDescription()
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}
	cmdLine.MergeContext(ctx)
	ts.Filter(ctx)
	for tag := range ts.GetTags() {
		fmt.Println(tag)
	}
	return nil
}

// CommandShowTemplates ...
func CommandShowTemplates(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {

	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}
	ts.Filter(ctx)
	ts.Filter(cmdLine)
	ts.FilterByStatus(STATUS_TEMPLATE)
	ts.SortByPriority()
	ts.DisplayByNext(false)
	ctx.PrintContextDescription()
	return nil
}

func CommandShowUnorganised(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}
	ts.Filter(cmdLine)
	ts.FilterUnorganised()
	ts.DisplayByNext(true)
	return nil
}

// CommandStart ...
func CommandStart(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}
	if len(cmdLine.IDs) > 0 {
		// start given tasks by IDs
		for _, id := range cmdLine.IDs {
			fmt.Println("trying to get ID", id)
			task := ts.MustGetByID(id)
			task.Status = STATUS_ACTIVE
			if cmdLine.Text != "" {
				task.Notes += "\n" + cmdLine.Text
			}
			ts.MustUpdateTask(task)

			ts.SavePendingChanges()
			MustGitCommit(repoPath, "Started %s", task)

			if task.Notes != "" {
				fmt.Printf("\nNotes on task %d:\n\033[38;5;245m%s\033[0m\n\n", task.ID, task.Notes)
			}
		}
	} else if cmdLine.Text != "" {
		// create a new task that is already active (started)
		cmdLine.MergeContext(ctx)
		task := Task{
			WritePending: true,
			Status:       STATUS_ACTIVE,
			Summary:      cmdLine.Text,
			Tags:         cmdLine.Tags,
			Project:      cmdLine.Project,
			Priority:     cmdLine.Priority,
			Notes:        cmdLine.Note,
		}
		task = ts.LoadTask(task)
		ts.SavePendingChanges()
		MustGitCommit(repoPath, "Added and started %s", task)
	}
	return nil

}

// CommandStop...
func CommandStop(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}
	for _, id := range cmdLine.IDs {
		task := ts.MustGetByID(id)
		task.Status = STATUS_PAUSED
		if cmdLine.Text != "" {
			task.Notes += "\n" + cmdLine.Text
		}
		ts.MustUpdateTask(task)
		ts.SavePendingChanges()
		MustGitCommit(repoPath, "Stopped %s", task)
	}
	return nil
}

func CommandSync(repoPath, idsFilePath, stateFilePath string) error {
	// TODO update repo w/ passed in path
	Sync(repoPath)
	return nil
}

// CommandTemplate...
func CommandTemplate(repoPath, idsFilePath, stateFilePath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath, idsFilePath, stateFilePath,
		WithStatuses(NON_RESOLVED_STATUSES...),
	)
	if err != nil {
		return err
	}

	if len(cmdLine.IDs) > 0 {
		for _, id := range cmdLine.IDs {
			task := ts.MustGetByID(id)
			task.Status = STATUS_TEMPLATE

			ts.MustUpdateTask(task)
			ts.SavePendingChanges()
			MustGitCommit(repoPath, "Changed %s to Template", task)
		}
	} else if cmdLine.Text != "" {
		ctx.PrintContextDescription()
		cmdLine.MergeContext(ctx)
		task := Task{
			WritePending: true,
			Status:       STATUS_TEMPLATE,
			Summary:      cmdLine.Text,
			Tags:         cmdLine.Tags,
			Project:      cmdLine.Project,
			Priority:     cmdLine.Priority,
			Notes:        cmdLine.Note,
		}
		task = ts.LoadTask(task)
		ts.SavePendingChanges()
		MustGitCommit(repoPath, "Created Template %s", task)
	}
	return nil

}

// CommandUndo...
func CommandUndo(repoPath, idsFilePath, stateFilePath string, args []string, ctx, cmdLine CmdLine) error {
	var err error
	n := 1
	if len(args) == 3 {
		n, err = strconv.Atoi(args[2])
		if err != nil {
			Help(CMD_UNDO)
			return err
		}
	}

	MustRunGitCmd(repoPath, "revert", "--no-gpg-sign", "--no-edit", "HEAD~"+strconv.Itoa(n)+"..")

	return nil
}

func CommandVersion() {
	fmt.Printf(
		"Version: %s\nGit commit: %s\nBuild date: %s\n",
		VERSION,
		GIT_COMMIT,
		BUILD_DATE,
	)
}
