package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mvdan/xurls"
	"github.com/naggie/dstask"
)

func main() {
	// Sets globals: GIT_REPO, STATE_FILE, IDS_FILE
	dstask.ParseConfig()
	dstask.EnsureRepoExists(dstask.GIT_REPO)
	repoPath := dstask.GIT_REPO
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
		if err := dstask.CommandNext(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_OPEN:
		if err := dstask.CommandShowOpen(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_ADD:
		if err := dstask.CommandAdd(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_RM, dstask.CMD_REMOVE:
		if err := dstask.CommandRemove(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_TEMPLATE:
		if err := dstask.CommandTemplate(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_LOG:
		if err := dstask.CommandLog(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_START:
		if err := dstask.CommandStart(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_STOP:
		if err := dstask.CommandStop(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_DONE, dstask.CMD_RESOLVE:
		if err := dstask.CommandDone(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_CONTEXT:
		if err := dstask.CommandContext(repoPath, state, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_MODIFY:
		if err := dstask.CommandModify(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_EDIT:
		if err := dstask.CommandEdit(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_NOTE, dstask.CMD_NOTES:
		if err := dstask.CommandNote(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
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
			task := ts.MustGetByID(id)
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
