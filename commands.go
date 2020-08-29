package dstask

import (
	"errors"
	"fmt"
)

func CommandAdd(repoPath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath,
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
		MustGitCommit("Added %s", task)
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
		MustGitCommit("Added %s", task)

	}
	return nil
}

// CommandNext prints the unresolved tasks associated with the current context.
// This is the default command.
func CommandNext(repoPath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath,
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

// CommandRemove ...
func CommandRemove(repoPath string, ctx, cmdLine CmdLine) error {
	if len(cmdLine.IDs) < 1 {
		return errors.New("missing argument: id")
	}
	ts, err := NewTaskSet(
		repoPath,
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
		MustGitCommit("Removed: %s", task)
	}
	return nil
}

// CommandShowOpen ...
func CommandShowOpen(repoPath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(
		repoPath,
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
