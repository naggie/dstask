package dstask

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func MustRunGitCmd(args ...string) {
	root := MustExpandHome(GIT_REPO)
	args = append([]string{"-C", root}, args...)
	err := MustRunCmd("git", args...)
	if err != nil {
		ExitFail("Failed to run git cmd.")
	}
}

func MustGitCommit(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)

	// git add all changed/created files
	// could optimise this to be given an explicit list of
	// added/modified/deleted files -- only if slow.
	fmt.Printf("\n%s\n", msg)
	fmt.Printf("\033[38;5;245m")
	MustRunGitCmd("add", ".")
	MustRunGitCmd("commit", "--no-gpg-sign", "-m", msg)
	fmt.Printf("\033[0m")
}

// leave file as an empty string to return directory
func MustGetRepoPath(directory, file string) string {
	root := MustExpandHome(GIT_REPO)
	dir := path.Join(root, directory)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0700)
		if err != nil {
			ExitFail("Failed to create directory in git repository")
		}
	}

	return path.Join(dir, file)
}
