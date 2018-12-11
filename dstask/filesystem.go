package dstask

// an interface to the filesystem/git based database -- loading, saving, committing

import (
	"path"
	"os"
)

func LoadTaskSetFromDisk(statuses []string) *TaskSet {
	gitRepoLocation := MustExpandHome(GIT_REPO)
	gitDotGitLocation := path.Join(gitRepoLocation, ".git")

	if _, err := os.Stat(gitDotGitLocation); os.IsNotExist(err) {
		ExitFail("Could not find git repository at "+gitRepoLocation+", please clone or create")
	}

	for _, status := range ALL_STATUSES {
		dir := path.Join(gitRepoLocation, status)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, 0755)
			if err != nil {
				ExitFail("Failed to create directory in git repository")
			}
		}
	}

	return &TaskSet{
		knownUuids:      make(map[string]bool),
		GitRepoLocation: MustExpandHome(GIT_REPO),
	}
}


func (t *Task) Save() error {
	return nil
}
