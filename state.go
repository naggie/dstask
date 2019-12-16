package dstask

import (
	"encoding/gob"
	"log"
	"os"
	"os/exec"
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

func MustWriteGob(filePath string, object interface{}) {
	file, err := os.Create(filePath)
	defer file.Close()

	if err != nil {
		ExitFail("Failed to open %s for writing: ", filePath)
	}

	encoder := gob.NewEncoder(file)
	encoder.Encode(object)
}

func MustReadGob(filePath string, object interface{}) {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		ExitFail("Failed to open %s for reading: ", filePath)
	}

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(object)

	if err != nil {
		ExitFail("Failed to parse gob: %s", filePath)
	}
}

func MustGetGitRef() string {
	root := MustExpandHome(GIT_REPO)
	out, err := exec.Command("git", "-C", root, "rev-parse", "HEAD").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(out))
}

// revert commits after last checkpoint if current commit is known
func (state *State) Undo() {
	ExitFail("Not possible to undo, please edit history manually in git repository")
}
