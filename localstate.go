package dstask

// this file represents the interface to the state specific to the PC dstask is
// on is stored. This is very minimal at the moment -- just the current
// context. It will probably remain that way.

import (
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"
)

// State models our local context for serialisation and deserialisation from
// our state file.
type State struct {
	// Context is an implicit command line that changes the behavior or display
	// of some commands.
	Context CmdLine
}

// Persistent DB of UUID -> ID to ensure that tasks have a persistent ID
// local to this machine for their lifetime. This is important to ensure
// the correct task is targeted between operations. Historically, each task
// stored its preferred ID but this resulted in merge conflicts when 2
// machines were using dstask concurrently on the same repository.
type IdsMap map[string]int

// Save serialises State to disk as gob binary data.
func (state State) Save() {
	os.MkdirAll(filepath.Dir(STATE_FILE), os.ModePerm)
	mustWriteGob(STATE_FILE, &state)
}

// LoadState reads the state file, if it exists. Otherwise a default State is returned.
func LoadState() State {
	if _, err := os.Stat(STATE_FILE); os.IsNotExist(err) {
		return State{}
	}

	state := State{}
	mustReadGob(STATE_FILE, &state)
	return state
}

// SetContext sets a context on State, with some validation.
func (state *State) SetContext(context CmdLine) error {
	if len(context.IDs) != 0 {
		return errors.New("Context cannot contain IDs")
	}

	if context.Text != "" {
		return errors.New("Context cannot contain text")
	}

	state.Context = context
	return nil
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

func (ids *IdsMap) Save() {
	os.MkdirAll(filepath.Dir(IDS_FILE), os.ModePerm)
	mustWriteGob(IDS_FILE, &ids)
}

func LoadIds() IdsMap {
	if _, err := os.Stat(IDS_FILE); os.IsNotExist(err) {
		return make(IdsMap)
	}

	ids := make(IdsMap, 1000)
	mustReadGob(IDS_FILE, &ids)
	return ids
}
