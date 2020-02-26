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
}

func (state State) Save() {
	os.MkdirAll(filepath.Dir(STATE_FILE), os.ModePerm)
	mustWriteGob(STATE_FILE, &state)
}

func LoadState() State {
	if _, err := os.Stat(STATE_FILE); os.IsNotExist(err) {
		return State{}
	}

	state := State{}
	mustReadGob(STATE_FILE, &state)
	return state
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
