package dstask

import (
	"os/user"
	"strings"
)

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

func (t *Task) InitDirectories() error {

}

func (t *Task) Save() error {
	return
}
