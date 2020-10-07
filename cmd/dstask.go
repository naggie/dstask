package main

import (
	"os"

	"github.com/naggie/dstask"
)

func main() {

	conf := dstask.NewConfig()

	dstask.EnsureRepoExists(conf.Repo)
	// Load state for getting and setting ctx
	state := dstask.LoadState(conf.StateFile)
	ctx := state.Context
	cmdLine := dstask.ParseCmdLine(os.Args[1:]...)

	if len(cmdLine.IDs) > 0 &&
		(lenNotZero(cmdLine.Tags, cmdLine.AntiTags, cmdLine.AntiProjects) || cmdLine.Project != "") {
		dstask.ExitFail("IDs cannot be combined with other attributes")
	}

	if cmdLine.IgnoreContext {
		ctx = dstask.CmdLine{}
	}

	switch cmdLine.Cmd {
	// The default command
	case "", dstask.CMD_NEXT:
		if err := dstask.CommandNext(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_OPEN:
		if err := dstask.CommandShowOpen(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_ADD:
		if err := dstask.CommandAdd(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_RM, dstask.CMD_REMOVE:
		if err := dstask.CommandRemove(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_TEMPLATE:
		if err := dstask.CommandTemplate(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_LOG:
		if err := dstask.CommandLog(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_START:
		if err := dstask.CommandStart(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_STOP:
		if err := dstask.CommandStop(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_DONE, dstask.CMD_RESOLVE:
		if err := dstask.CommandDone(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_CONTEXT:
		if err := dstask.CommandContext(conf, state, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_MODIFY:
		if err := dstask.CommandModify(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_EDIT:
		if err := dstask.CommandEdit(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_NOTE, dstask.CMD_NOTES:
		if err := dstask.CommandNote(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_UNDO:
		if err := dstask.CommandUndo(conf, os.Args, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SYNC:
		if err := dstask.CommandSync(conf.Repo); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_GIT:
		dstask.MustRunGitCmd(conf.Repo, os.Args[2:]...)

	case dstask.CMD_SHOW_ACTIVE:
		if err := dstask.CommandShowActive(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_PAUSED:
		if err := dstask.CommandShowPaused(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_OPEN:
		if err := dstask.CommandOpen(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_IMPORT_TW:
		if err := dstask.CommandImportTW(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_PROJECTS:
		if err := dstask.CommandShowProjects(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_TAGS:
		if err := dstask.CommandShowTags(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_TEMPLATES:
		if err := dstask.CommandShowTemplates(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_RESOLVED:
		if err := dstask.CommandShowResolved(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_UNORGANISED:
		if err := dstask.CommandShowUnorganised(conf, ctx, cmdLine); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_HELP:
		dstask.CommandHelp(os.Args)

	case dstask.CMD_VERSION:
		dstask.CommandVersion()

	case dstask.CMD_COMPLETIONS:
		dstask.Completions(conf, os.Args, ctx)

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

func lenNotZero(arrays ...[]string) bool {
	for _, arr := range arrays {
		if len(arr) > 0 {
			return true
		}
	}
	return false
}
