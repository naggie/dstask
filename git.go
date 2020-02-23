package dstask

import (
	"fmt"
	"os"
	"path"
)

func MustRunGitCmd(args ...string) {
	root := MustExpandHome(GIT_REPO)
	args = append([]string{"-C", root}, args...)
	err := RunCmd("git", args...)
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

// TODO check git exists within path
func InitialiseRepo() {
	gitDotGitLocation := MustExpandHome(path.Join(GIT_REPO, ".git"))

	if _, err := os.Stat(gitDotGitLocation); os.IsNotExist(err) {
		ExitFail("Could not find git repository at " + GIT_REPO + ", please clone or create. Try `dstask help` for more information.")
	}
}

func Sync() {
	dstask.MustRunGitCmd("pull", "--no-edit", "--commit", "origin", "master")
	dstask.MustRunGitCmd("push", "origin", "master")
}
