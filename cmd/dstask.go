package main

import (
	"os"

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
	// The default command
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

	case dstask.CMD_UNDO:
		if err := dstask.CommandUndo(repoPath, os.Args, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SYNC:
		if err := dstask.CommandSync(repoPath); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_GIT:
		dstask.MustRunGitCmd(os.Args[2:]...)

	case dstask.CMD_SHOW_ACTIVE:
		if err := dstask.CommandShowActive(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_PAUSED:
		if err := dstask.CommandShowPaused(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_OPEN:
		if err := dstask.CommandOpen(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_IMPORT_TW:
		if err := dstask.CommandImportTW(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_PROJECTS:
		if err := dstask.CommandShowProjects(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_TAGS:
		if err := dstask.CommandShowTags(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_TEMPLATES:
		if err := dstask.CommandShowTemplates(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_RESOLVED:
		if err := dstask.CommandShowResolved(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_UNORGANISED:
		if err := dstask.CommandShowUnorganised(repoPath, context, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_HELP:
		dstask.CommandHelp(os.Args)

	case dstask.CMD_VERSION:
		dstask.CommandVersion()

	case dstask.CMD_COMPLETIONS:
		dstask.Completions(os.Args, context)

	default:
		panic("this should never happen?")
	}

}
