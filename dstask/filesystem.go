package dstask

// an interface to the filesystem/git based database -- loading, saving, committing

import (
	"os/user"
	"strings"
)

func LoadTaskSetFromDisk(statuses []string) *TaskSet {
	return &TaskSet{
		knownUuids: make(map[string]bool),
	}
}

func ExpandHome(path string) (string, err) {
	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err != nil {
			return err
		}
		return usr.HomeDir + path[2:len(path)]
	} else {
		return path
	}
}

func (t *Task) Save() error {
	return
}
