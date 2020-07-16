package dstask

import (
	"fmt"
	"os"
	"path"
)

func RunGitCmd(args ...string) error {
	args = append([]string{"-C", GIT_REPO}, args...)
	return RunCmd("git", args...)
}

func MustRunGitCmd(args ...string) {
	err := RunGitCmd(args...)
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

	// check for changes -- returns exit status 1 on change
	if RunGitCmd("diff-index", "--quiet", "HEAD", "--") == nil {
		fmt.Println("No changes detected")
		return
	}

	MustRunGitCmd("commit", "--no-gpg-sign", "-m", msg)
	fmt.Printf("\033[0m")
}

// leave file as an empty string to return directory
func MustGetRepoPath(directory, file string) string {
	dir := path.Join(GIT_REPO, directory)

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
	gitDotGitLocation := path.Join(GIT_REPO, ".git")

	if _, err := os.Stat(gitDotGitLocation); os.IsNotExist(err) {
		ExitFail("Could not find git repository at " + GIT_REPO + ", please clone or create. Try `dstask help` for more information.")
	}
}

func Sync() {
	MustRunGitCmd("pull", "--no-rebase", "--no-edit", "--commit", "origin", "master")
	MustRunGitCmd("push", "origin", "master")
}
