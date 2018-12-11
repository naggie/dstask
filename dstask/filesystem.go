package dstask

// an interface to the filesystem/git based database -- loading, saving, committing

import (
	"os/user"
	"strings"
)

func LoadTaskSetFromDisk(statuses []string) *TaskSet {
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
