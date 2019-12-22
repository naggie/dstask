package dstask

import (
	"encoding/gob"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type State struct {
	// context to automatically apply to all queries and new tasks
	context CmdLine
	// last command -- joined args. Used to confirm an undo
	cmd string
	// git ref before the last consequential command
	preCmdRef string
	// git ref after the last consequential command (if does not match HEAD.
	// undo should fail) -- this can happen as a consequence of sync.
	postCmdRef string
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

func (state *State) SetPreCmdRef() {
	state.preCmdRef = MustGetGitRef()
}

func (state *State) SetPostCmdRef() {
	state.postCmdRef = MustGetGitRef()
	state.cmd = strings.Join(os.Args[1:], " ")
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
	if state.preCmdRef == "" or state.postCmdRef == "" or state.cmd = "":
		ExitFail("Last command not recorded")


	ConfirmOrAbort("This will undo the last command on this computer which was:\n    %s\nContinue?")

	// https://stackoverflow.com/questions/4991594/revert-a-range-of-commits-in-git
	// revert all without committing, then make a single commit
	MustRunGitCmd("revert", "-n", state.preCmdRef+"^.."+state.postCmdRef)
	MustRunGitCmd("commit", "--no-gpg-sign", "-m", "Undo: "+state.cmd)

	state.cmd = ""
	state.preCmdRef = ""
	state.postCmdRef = ""
	state.Save()
}
