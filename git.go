package dstask

import (
	"fmt"
	"os"
	"path"
)

// RunGitCmd shells out to git in the context of the dstask repo.
func RunGitCmd(args ...string) error {
	args = append([]string{"-C", GIT_REPO}, args...)
	return RunCmd("git", args...)
}

// MustRunGitCmd delegates to RunGitCmd and exits the program on any error.
func MustRunGitCmd(args ...string) {
	err := RunGitCmd(args...)
	if err != nil {
		ExitFail("Failed to run git cmd.")
	}
}

// MustGitCommit stages changes in the dstask repository and commits them. If
// any error is encountered, the program exits.
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

// MustGetRepoPath returns the full path to a file within the dstask git repo.
// Pass file as an empty string to return the git repo directory itself.
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

// EnsureRepoExists checks for the existence of a dstask repository, or exits the program.
func EnsureRepoExists(repoPath string) {
	// TODO make sure git is installed
	gitDotGitLocation := path.Join(repoPath, ".git")

	if _, err := os.Stat(gitDotGitLocation); os.IsNotExist(err) {
		ExitFail("Could not find git repository at " + repoPath + ", please clone or create. Try `dstask help` for more information.")
	}
}

// Sync performs a git pull, and then a git push. If any conflicts are encountered,
// the user will need to resolve them.
func Sync() {
	MustRunGitCmd("pull", "--no-rebase", "--no-edit", "--commit", "origin", "master")
	MustRunGitCmd("push", "origin", "master")
}
