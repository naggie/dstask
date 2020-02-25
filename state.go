package dstask

// this file represents the interface to the state specific to the PC dstask is
// on is stored. This is very minimal at the moment -- just the current
// context. It will probably remain that way.

import (
	"encoding/gob"
	"os"
	"path/filepath"
)

// note that fields must be exported for gob marshalling to work.
type State struct {
	// context to automatically apply to all queries and new tasks
	Context CmdLine
	// DB of UUID -> ID to ensure that tasks have a persistent ID local to this
	// machine for their lifetime. This is important to ensure the correct task
	// is targeted between operations. Historically, each task stored its
	// preferred ID but this resulted in merge conflicts when 2 machines were
	// using dstask concurrently on the same repository.
	IDmap map[int]string
}

func (state State) Save() {
	fp := MustExpandHome(STATE_FILE)
	os.MkdirAll(filepath.Dir(fp), os.ModePerm)
	mustWriteGob(fp, &state)
}

func LoadState() State {
	fp := MustExpandHome(STATE_FILE)
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		return State{}
	}

	state := State{}
	mustReadGob(fp, &state)
	return state
}

func (state State) GetContext() CmdLine {
	return state.Context
}

func (state *State) SetContext(context CmdLine) {
	if len(context.IDs) != 0 {
		ExitFail("Context cannot contain IDs")
	}

	if context.Text != "" {
		ExitFail("Context cannot contain text")
	}

	state.Context = context
}

func (state State) GetIDCache() CmdLine {
	return state.IDmap
}

func (state *State) SetIDCache(idCache map[int]string) {
	state.IDmap = idCache;
}

func (state *State) ClearContext() {
	state.SetContext(CmdLine{})
}

func mustWriteGob(filePath string, object interface{}) {
	file, err := os.Create(filePath)
	defer file.Close()

	if err != nil {
		ExitFail("Failed to open %s for writing: ", filePath)
	}

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(object)

	if err != nil {
		ExitFail("Failed to encode state gob: %s, %s", filePath, err)
	}
}

func mustReadGob(filePath string, object interface{}) {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		ExitFail("Failed to open %s for reading: ", filePath)
	}

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(object)

	if err != nil {
		ExitFail("Failed to parse state gob: %s, %s", filePath, err)
	}
}
