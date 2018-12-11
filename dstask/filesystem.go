package dstask

// an interface to the filesystem/git based database -- loading, saving, committing

import (
	"path"
	"os"
)

func MustGetRepoDirectory(directory ...string) string {
	root := MustExpandHome(GIT_REPO)
	return path.Join(append([]string{root}, directory...)...)
}

func LoadTaskSetFromDisk(statuses []string) *TaskSet {
	gitDotGitLocation := MustGetRepoDirectory(".git")

	if _, err := os.Stat(gitDotGitLocation); os.IsNotExist(err) {
		ExitFail("Could not find git repository at "+GIT_REPO+", please clone or create")
	}

	for _, status := range ALL_STATUSES {
		dir := MustGetRepoDirectory(status)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, 0755)
			if err != nil {
				ExitFail("Failed to create directory in git repository")
			}
		}
	}

	return &TaskSet{
		knownUuids: make(map[string]bool),
	}
}


func (t *Task) Save() error {
	return nil
}
