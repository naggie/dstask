package dstask

import (
	"fmt"
	"os"
	"path"
	"strings"
	"os/user"
)

func ExitFail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func MustExpandHome(filepath string) string {
	if strings.HasPrefix(filepath, "~/") {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		return path.Join(usr.HomeDir, filepath[2:len(filepath)])
	} else {
		return filepath
	}
}
