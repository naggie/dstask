package main

import (
	"os"
	"path"

	"github.com/naggie/dstask"
)

func main() {

	repoPath := getEnv("DSTASK_GIT_REPO", os.ExpandEnv("$HOME/.dstask"))
	dstask.EnsureRepoExists(repoPath)
	stateFilePath := path.Join(repoPath, ".git", "dstask", "state.bin")
	idsFilePath := path.Join(repoPath, ".git", "dstask", "ids.bin")
	// Load state for getting and setting ctx
	state := dstask.LoadState(stateFilePath)
	ctx := state.Context
	cmdLine := dstask.ParseCmdLine(os.Args[1:]...)

	if cmdLine.IgnoreContext {
		ctx = dstask.CmdLine{}
	}

	switch cmdLine.Cmd {
	// The default command
	case "", dstask.CMD_NEXT:
		if err := dstask.CommandNext(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_OPEN:
		if err := dstask.CommandShowOpen(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_ADD:
		if err := dstask.CommandAdd(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_RM, dstask.CMD_REMOVE:
		if err := dstask.CommandRemove(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_TEMPLATE:
		if err := dstask.CommandTemplate(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_LOG:
		if err := dstask.CommandLog(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_START:
		if err := dstask.CommandStart(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_STOP:
		if err := dstask.CommandStop(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_DONE, dstask.CMD_RESOLVE:
		if err := dstask.CommandDone(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_CONTEXT:
		if err := dstask.CommandContext(stateFilePath, state, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_MODIFY:
		if err := dstask.CommandModify(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_EDIT:
		if err := dstask.CommandEdit(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_NOTE, dstask.CMD_NOTES:
		if err := dstask.CommandNote(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_UNDO:
		if err := dstask.CommandUndo(repoPath, idsFilePath, stateFilePath, os.Args, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SYNC:
		if err := dstask.CommandSync(repoPath); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_GIT:
		dstask.MustRunGitCmd(repoPath, os.Args[2:]...)

	case dstask.CMD_SHOW_ACTIVE:
		if err := dstask.CommandShowActive(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_PAUSED:
		if err := dstask.CommandShowPaused(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_OPEN:
		if err := dstask.CommandOpen(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_IMPORT_TW:
		if err := dstask.CommandImportTW(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_PROJECTS:
		if err := dstask.CommandShowProjects(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_TAGS:
		if err := dstask.CommandShowTags(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_TEMPLATES:
		if err := dstask.CommandShowTemplates(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_RESOLVED:
		if err := dstask.CommandShowResolved(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_UNORGANISED:
		if err := dstask.CommandShowUnorganised(repoPath, idsFilePath, stateFilePath, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_HELP:
		dstask.CommandHelp(os.Args)

	case dstask.CMD_VERSION:
		dstask.CommandVersion()

	case dstask.CMD_COMPLETIONS:
		dstask.Completions(repoPath, idsFilePath, stateFilePath, os.Args, ctx)

	default:
		panic("this should never happen?")
	}
}

// getEnv returns an env var's value, or a default.
func getEnv(key string, _default string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return _default
}
