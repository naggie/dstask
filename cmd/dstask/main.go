package main

import (
	"os"
	"strings"

	"github.com/naggie/dstask"
)

func main() {
	// special case: allow users to run the help command without requiring
	// initialisation. Other commands are handled in the below switch statement
	// after initialisation.
	if len(os.Args) > 1 && os.Args[1] == dstask.CMD_HELP {
		dstask.CommandHelp(os.Args)
	}

	conf := dstask.NewConfig()
	dstask.EnsureRepoExists(conf.Repo)

	// Load state for getting and setting ctx
	state := dstask.LoadState(conf.StateFile)
	ctx := state.Context
	query := dstask.ParseQuery(os.Args[1:]...)

	// Check if we have a context override.
	if conf.CtxFromEnvVar != "" {
		if query.Cmd == dstask.CMD_CONTEXT && len(os.Args) >= 3 {
			dstask.ExitFail("setting context not allowed while DSTASK_CONTEXT is set")
		}
		splitted := strings.Fields(conf.CtxFromEnvVar)
		ctx = dstask.ParseQuery(splitted...)
	}

	// Check if we ignore context with the "--" token
	if query.IgnoreContext {
		ctx = dstask.Query{}
	}

	switch query.Cmd {
	// The default command
	case "", dstask.CMD_NEXT, dstask.CMD_SHOW_NEXT:
		if err := dstask.CommandNext(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_OPEN:
		if err := dstask.CommandShowOpen(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_ADD:
		if err := dstask.CommandAdd(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_RM, dstask.CMD_REMOVE:
		if err := dstask.CommandRemove(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_TEMPLATE:
		if err := dstask.CommandTemplate(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_LOG:
		if err := dstask.CommandLog(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_START:
		if err := dstask.CommandStart(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_STOP:
		if err := dstask.CommandStop(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_DONE, dstask.CMD_RESOLVE:
		if err := dstask.CommandDone(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_CONTEXT:
		if err := dstask.CommandContext(conf, state, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_MODIFY:
		if err := dstask.CommandModify(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_EDIT:
		if err := dstask.CommandEdit(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_NOTE, dstask.CMD_NOTES:
		if err := dstask.CommandNote(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_UNDO:
		if err := dstask.CommandUndo(conf, os.Args, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SYNC:
		if err := dstask.CommandSync(conf.Repo); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_GIT:
		dstask.MustRunGitCmd(conf.Repo, os.Args[2:]...)

	case dstask.CMD_SHOW_ACTIVE:
		if err := dstask.CommandShowActive(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_PAUSED:
		if err := dstask.CommandShowPaused(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_OPEN:
		if err := dstask.CommandOpen(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_PROJECTS:
		if err := dstask.CommandShowProjects(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_TAGS:
		if err := dstask.CommandShowTags(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_TEMPLATES:
		if err := dstask.CommandShowTemplates(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_RESOLVED:
		if err := dstask.CommandShowResolved(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_SHOW_UNORGANISED:
		if err := dstask.CommandShowUnorganised(conf, ctx, query); err != nil {
			dstask.ExitFail(err.Error())
		}

	case dstask.CMD_VERSION:
		dstask.CommandVersion()

	case dstask.CMD_COMPLETIONS:
		dstask.Completions(conf, os.Args, ctx)

	default:
		panic("this should never happen?")
	}
}
