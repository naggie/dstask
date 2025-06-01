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
	Context Query
}

// Persistent DB of UUID -> ID to ensure that tasks have a persistent ID
// local to this machine for their lifetime. This is important to ensure
// the correct task is targeted between operations. Historically, each task
// stored its preferred ID but this resulted in merge conflicts when 2
// machines were using dstask concurrently on the same repository.
type IdsMap map[string]int

// Save serialises State to disk as gob binary data.
func (state State) Save(stateFilePath string) {
	if err := os.MkdirAll(filepath.Dir(stateFilePath), os.ModePerm); err != nil {
		ExitFail("Failed to create directories for %s: %s", stateFilePath, err)
	}

	mustWriteGob(stateFilePath, &state)
}

// LoadState reads the state file, if it exists. Otherwise a default State is returned.
func LoadState(stateFilePath string) State {
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		return State{}
	}

	state := State{}
	mustReadGob(stateFilePath, &state)

	return state
}

// SetContext sets a context on State, with some validation.
func (state *State) SetContext(context Query) error {
	if len(context.IDs) != 0 {
		return errors.New("context cannot contain IDs")
	}

	if context.Text != "" {
		return errors.New("context cannot contain text")
	}

	state.Context = context

	return nil
}

func mustWriteGob(filePath string, object any) {
	file, err := os.Create(filePath)
	if err != nil {
		ExitFail("Failed to open %s for writing: ", filePath)
	}

	defer func() {
		if err := file.Close(); err != nil {
			ExitFail("Failed to close file: %v", err)
		}
	}()

	encoder := gob.NewEncoder(file)

	err = encoder.Encode(object)
	if err != nil {
		ExitFail("Failed to encode state gob: %s, %s", filePath, err)
	}
}

func mustReadGob(filePath string, object any) {
	file, err := os.Open(filePath)
	if err != nil {
		ExitFail("Failed to open %s for reading: ", filePath)
	}

	defer func() {
		if err := file.Close(); err != nil {
			ExitFail("Failed to close file: %v", err)
		}
	}()

	decoder := gob.NewDecoder(file)

	err = decoder.Decode(object)
	if err != nil {
		ExitFail("Failed to parse state gob: %s, %s", filePath, err)
	}
}

func (ids *IdsMap) Save(idsFilePath string) {
	if err := os.MkdirAll(filepath.Dir(idsFilePath), os.ModePerm); err != nil {
		ExitFail("Failed to create directories for %s: %s", idsFilePath, err)
	}

	mustWriteGob(idsFilePath, &ids)
}

func LoadIds(idsFilePath string) IdsMap {
	if _, err := os.Stat(idsFilePath); os.IsNotExist(err) {
		return make(IdsMap)
	}

	ids := make(IdsMap, 1000)
	mustReadGob(idsFilePath, &ids)

	return ids
}
