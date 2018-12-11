package dstask

// an interface to the filesystem/git based database -- loading, saving, committing

import (
	"os/user"
	"strings"
	"path"
	"os"
)

func LoadTaskSetFromDisk(statuses []string) *TaskSet {
	GitRepoLocation := MustExpandHome(GIT_REPO)

	for _, status := range ALL_STATUSES {
		dir := path.Join(GitRepoLocation, status)
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

func MustExpandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		return usr.HomeDir + path[2:len(path)]
	} else {
		return path
	}
}

func (t *Task) Save() error {
	return nil
}
