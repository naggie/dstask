package dstask

import (
	"os"
	"path/filepath"
)

type State struct {
	Context CmdLine
	// git ref before the last consequential command
	LastOperationFrom string
	// git ref after the last consequential command (if does not match HEAD.
	// undo should fail) -- this can happen as a consequence of sync.
	LastOperationTo string
}

// TODO separate validate context fn then move to context cmd
func SaveState(state State) {
	if len(state.Context.IDs) != 0 {
		ExitFail("Context cannot contain IDs")
	}

	if state.Context.Text != "" {
		ExitFail("Context cannot contain text")
	}

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

func (*State) SetGitRef() {

}
