package dstask

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type State struct {
	context CmdLine
	// git ref before the last consequential command
	checkpoint string
	// git ref after the last consequential command (if does not match HEAD.
	// undo should fail) -- this can happen as a consequence of sync.
	lastKnown string
	// last command -- joined args. Used to confirm an undo
	// TODO confirm undo
	lastChangeCmd string
	// when did change occur? more than a day?
	lastChangeTime time.Time
}

// TODO separate validate context fn then move to context cmd
func (state State) Save() {
	fp := MustExpandHome(STATE_FILE)
	os.MkdirAll(filepath.Dir(fp), os.ModePerm)
	MustWriteGob(fp, &state)
}

func LoadState() State {
	fp := MustExpandHome(STATE_FILE)
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		return State{}
	}

	state := State{}
	MustReadGob(fp, &state)
	return state
}

func (state State) GetContext() CmdLine {
	return state.context
}

func (state *State) SetContext(context CmdLine) {
	if len(context.IDs) != 0 {
		ExitFail("Context cannot contain IDs")
	}

	if context.Text != "" {
		ExitFail("Context cannot contain text")
	}

	state.context = context
}

func (state *State) ClearContext() {
	state.SetContext(CmdLine{})
}

func (state *State) SetCheckpoint() {
	state.checkpoint = MustGetGitRef()
}

func (state *State) SetLastCmd() {
	state.lastChangeCmd = strings.Join(os.Args[1:], " ")
	state.lastKnown = MustGetGitRef()
}
