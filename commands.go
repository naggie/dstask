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

// CommandAdd adds a new task to the task database.
func CommandAdd(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}
	if query.Template > 0 {
		var taskSummary string
		tt := ts.MustGetByID(query.Template)
		query.MergeContext(ctx)

		if query.Text != "" {
			taskSummary = query.Text
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

		// Modify the task with any tags/projects/antiProjects/priorities in query
		task.Modify(query)

		task = ts.LoadTask(task)
		ts.SavePendingChanges()
		MustGitCommit(conf.Repo, "Added %s", task)
		if tt.Status != STATUS_TEMPLATE {
			// Insert Text Statement to inform user of real Templates
			fmt.Print("\nYou've copied an open task!\nTo learn more about creating templates enter 'dstask help template'\n\n")
		}
	} else if query.Text != "" {
		ctx.PrintContextDescription()
		query.MergeContext(ctx)
		task := Task{
			WritePending: true,
			Status:       STATUS_PENDING,
			Summary:      query.Text,
			Tags:         query.Tags,
			Project:      query.Project,
			Priority:     query.Priority,
			Notes:        query.Note,
		}
		task = ts.LoadTask(task)
		ts.SavePendingChanges()
		MustGitCommit(conf.Repo, "Added %s", task)

	}
	return nil
}

// CommandContext sets a global context for dstask.
func CommandContext(conf Config, state State, ctx, query Query) error {
	if len(os.Args) < 3 {
		fmt.Printf("Current context: %s\n", ctx)
	} else if os.Args[2] == "none" {
		if err := state.SetContext(Query{}); err != nil {
			return err
		}
	} else {
		if err := state.SetContext(query); err != nil {
			return err
		}
	}

	state.Save(conf.StateFile)
	return nil
}

// CommandDone marks a task as done.
func CommandDone(conf Config, ctx, query Query) error {
	if query.HasOperators() {
		return errors.New("Operators not valid in this context.")
	}

	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}

	if len(query.IDs) == 0 {
		return errors.New("No ID(s) specified")
	}

	// iterate over IDs instead of filtering; it's clearer and enables us to
	// test each ID exists.
	for _, id := range query.IDs {
		task := ts.MustGetByID(id)
		task.Status = STATUS_RESOLVED
		task.Resolved = time.Now()
		if query.Text != "" {
			task.Notes += "\n" + query.Text
		}
		ts.MustUpdateTask(task)
		ts.SavePendingChanges()
		MustGitCommit(conf.Repo, "Resolved %s", task)
	}

	return nil
}

// CommandEdit edits a task's metadata, such as status, projects, tags, etc.
func CommandEdit(conf Config, ctx, query Query) error {
	if query.HasOperators() {
		return errors.New("Operators not valid in this context.")
	}

	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}

	if len(query.IDs) == 0 {
		return errors.New("No ID(s) specified")
	}

	for _, id := range query.IDs {
		task := ts.MustGetByID(id)
		data, err := yaml.Marshal(&task)
		if err != nil {
			// TODO present error to user, specific error message is important
			return fmt.Errorf("failed to marshal task %s", task)
		}

		for {
			edited := MustEditBytes(data, "yml")
			err = yaml.Unmarshal(edited, &task)
			if err == nil {
				break
			} else {
				// edit is a special case that won't be used as part of an API,
				// so it's OK to exit
				ConfirmOrAbort("Failed to unmarshal %s\nTry again?", err)
			}
		}

		ts.MustUpdateTask(task)
		ts.SavePendingChanges()
		MustGitCommit(conf.Repo, "Edited %s", task)
	}
	return nil
}

// CommandHelp prints for a specific command or all commands.
func CommandHelp(args []string) {
	if len(os.Args) > 2 {
		Help(os.Args[2])
	} else {
		Help("")
	}
}

// CommandImportTW imports a taskwarrior database.
func CommandImportTW(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, true)
	if err != nil {
		return err
	}
	ts.ImportFromTaskwarrior()
	ts.SavePendingChanges()
	MustGitCommit(conf.Repo, "Import from taskwarrior")
	return nil
}

// CommandLog logs a completed task immediately. Useful for tracking tasks after
// they're already completed.
func CommandLog(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)

	if err != nil {
		return err
	}

	query = query.Merge(ctx)

	if query.Text == "" {
		return errors.New("Task description required")
	}

	ctx.PrintContextDescription()
	query.MergeContext(ctx)
	task := Task{
		WritePending: true,
		Status:       STATUS_RESOLVED,
		Summary:      query.Text,
		Tags:         query.Tags,
		Project:      query.Project,
		Priority:     query.Priority,
		Resolved:     time.Now(),
	}
	task = ts.LoadTask(task)
	ts.SavePendingChanges()
	MustGitCommit(conf.Repo, "Logged %s", task)

	return nil
}

// CommandModify applies a change to tasks specified by ID, or all tasks in
// current context
func CommandModify(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}

	if !query.HasOperators() {
		return errors.New("No operations specified")
	}

	if len(query.IDs) == 0 {
		ts.Filter(ctx)

		if StdoutIsTTY() {
			ConfirmOrAbort("No IDs specified. Apply to all %d tasks in current ctx?", len(ts.Tasks()))
		}

		for _, task := range ts.Tasks() {
			task.Modify(query)
			ts.MustUpdateTask(task)
			ts.SavePendingChanges()
			MustGitCommit(conf.Repo, "Modified %s", task)
		}
	} else {
		for _, id := range query.IDs {
			task := ts.MustGetByID(id)
			task.Modify(query)
			ts.MustUpdateTask(task)
			ts.SavePendingChanges()
			MustGitCommit(conf.Repo, "Modified %s", task)
		}
	}

	return nil
}

// CommandNext prints the unresolved tasks associated with the current context.
// This is the default command.
func CommandNext(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}

	query = query.Merge(ctx)
	ts.Filter(query)
	ts.DisplayByNext(ctx, true)

	return nil
}

// CommandNote edits or prints the markdown note associated with the task.
func CommandNote(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}

	for _, task := range ts.Tasks() {
		// If stdout is a TTY, we open the editor
		if StdoutIsTTY() {
			if query.Text == "" {
				task.Notes = string(MustEditBytes([]byte(task.Notes), "md"))
			} else {
				if task.Notes == "" {
					task.Notes = query.Text
				} else {
					task.Notes += "\n" + query.Text
				}
			}
			ts.MustUpdateTask(task)
			ts.SavePendingChanges()
			MustGitCommit(conf.Repo, "Edit note %s", task)
		} else {
			// If stdout is not a TTY, we simply write markdown notes to stdout
			if err := WriteStdout([]byte(task.Notes)); err != nil {
				ExitFail("Could not write to stdout: %v", err)
			}
		}
	}
	return nil
}

// CommandOpen opens a task URL in the browser, if the task has a URL.
func CommandOpen(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}
	for _, task := range ts.Tasks() {
		urls := xurls.Relaxed.FindAllString(task.Summary+" "+task.Notes, -1)
		if len(urls) == 0 {
			return fmt.Errorf("no URLs found in task %v", task.ID)
		}
		for _, url := range urls {
			MustOpenBrowser(url)
		}
	}

	return nil
}

// CommandRemove removes a task by ID from the database.
func CommandRemove(conf Config, ctx, query Query) error {
	if len(query.IDs) < 1 {
		return errors.New("missing argument: id")
	}
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}

	for _, task := range ts.Tasks() {
		fmt.Println(task)
	}

	if StdoutIsTTY() {
		ConfirmOrAbort("\nThe above %d task(s) will be deleted without checking subtasks. Continue?", len(ts.Tasks()))
	}

	for _, task := range ts.Tasks() {
		// Mark our task for deletion
		task.Deleted = true

		// MustUpdateTask validates and normalises our task object
		ts.MustUpdateTask(task)
		ts.SavePendingChanges()

		if query.Text != "" {
			// commit comment, put in body
			MustGitCommit(conf.Repo, "Removed: %s\n\n%s", task, query.Text)
		} else {
			MustGitCommit(conf.Repo, "Removed: %s", task)
		}
	}
	return nil
}

// CommandShowActive prints a list of active tasks.
func CommandShowActive(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}
	ts.DisplayByNext(ctx, true)

	return nil
}

// CommandShowProjects prints a list of projects associated with all tasks.
func CommandShowProjects(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}
	ts.DisplayProjects()
	return nil
}

// CommandShowOpen prints a list of open tasks.
func CommandShowOpen(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}
	ts.DisplayByNext(ctx, false)
	return nil
}

// CommandShowPaused prints a list of paused tasks.
func CommandShowPaused(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}
	ts.DisplayByNext(ctx, true)
	return nil
}

// CommandShowResolved prints a list of resolved tasks.
func CommandShowResolved(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, true)
	if err != nil {
		return err
	}
	ts.DisplayByWeek()
	return nil
}

// CommandShowTags prints a list of all tags associated with non-resolved tasks.
func CommandShowTags(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)

	if err != nil {
		return err
	}
	for tag := range ts.GetTags() {
		fmt.Println(tag)
	}
	return nil
}

// CommandShowTemplates show a list of task templates.
func CommandShowTemplates(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)

	if err != nil {
		return err
	}

	ts.DisplayByNext(ctx, false)
	return nil
}

// CommandShowUnorganised prints a list of tasks without tags or projects.
func CommandShowUnorganised(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}
	ts.DisplayByNext(ctx, true)
	return nil
}

// CommandStart marks a task as started.
func CommandStart(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}
	if len(query.IDs) > 0 {
		// start given tasks by IDs
		for _, id := range query.IDs {
			task := ts.MustGetByID(id)
			task.Status = STATUS_ACTIVE
			if query.Text != "" {
				task.Notes += "\n" + query.Text
			}
			ts.MustUpdateTask(task)

			ts.SavePendingChanges()
			MustGitCommit(conf.Repo, "Started %s", task)

			if task.Notes != "" {
				fmt.Printf("\nNotes on task %d:\n\033[38;5;245m%s\033[0m\n\n", task.ID, task.Notes)
			}
		}
	} else if query.Text != "" {
		// create a new task that is already active (started)
		query.MergeContext(ctx)
		task := Task{
			WritePending: true,
			Status:       STATUS_ACTIVE,
			Summary:      query.Text,
			Tags:         query.Tags,
			Project:      query.Project,
			Priority:     query.Priority,
			Notes:        query.Note,
		}
		task = ts.LoadTask(task)
		ts.SavePendingChanges()
		MustGitCommit(conf.Repo, "Added and started %s", task)
	}
	return nil

}

// CommandStop marks a task as stopped.
func CommandStop(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}
	for _, task := range ts.Tasks() {
		task.Status = STATUS_PAUSED
		if query.Text != "" {
			task.Notes += "\n" + query.Text
		}
		ts.MustUpdateTask(task)
		ts.SavePendingChanges()
		MustGitCommit(conf.Repo, "Stopped %s", task)
	}
	return nil
}

// CommandSync pushes and pulls task database changes from the remote repository.
func CommandSync(repoPath string) error {
	// TODO(dontlaugh) return error
	Sync(repoPath)
	return nil
}

// CommandTemplate creates a new task template.
func CommandTemplate(conf Config, ctx, query Query) error {
	ts, err := LoadTaskSet(conf.Repo, conf.IDsFile, false)
	if err != nil {
		return err
	}

	if len(query.IDs) > 0 {
		for _, id := range query.IDs {
			task := ts.MustGetByID(id)
			task.Status = STATUS_TEMPLATE

			ts.MustUpdateTask(task)
			ts.SavePendingChanges()
			MustGitCommit(conf.Repo, "Changed %s to Template", task)
		}
	} else if query.Text != "" {
		query.MergeContext(ctx)
		task := Task{
			WritePending: true,
			Status:       STATUS_TEMPLATE,
			Summary:      query.Text,
			Tags:         query.Tags,
			Project:      query.Project,
			Priority:     query.Priority,
			Notes:        query.Note,
		}
		task = ts.LoadTask(task)
		ts.SavePendingChanges()
		MustGitCommit(conf.Repo, "Created Template %s", task)
	}
	return nil

}

// CommandUndo performs undo with git revert.
func CommandUndo(conf Config, args []string, ctx, query Query) error {
	var err error
	n := 1
	if len(args) == 3 {
		n, err = strconv.Atoi(args[2])
		if err != nil {
			Help(CMD_UNDO)
			return err
		}
	}

	MustRunGitCmd(conf.Repo, "revert", "--no-gpg-sign", "--no-edit", "HEAD~"+strconv.Itoa(n)+"..")

	return nil
}

// CommandVersion prints version information for the dstask binary.
func CommandVersion() {
	fmt.Printf(
		"Version: %s\nGit commit: %s\nBuild date: %s\n",
		VERSION,
		GIT_COMMIT,
		BUILD_DATE,
	)
}
