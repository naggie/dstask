package dstask

import (
	"fmt"
	"os"
	"path"
	"os/exec"
)

// RunGitCmd shells out to git in the context of the dstask repo.
func RunGitCmd(repoPath string, args ...string) error {
	args = append([]string{"-C", repoPath}, args...)
	return RunCmd("git", args...)
}

// MustRunGitCmd delegates to RunGitCmd and exits the program on any error.
func MustRunGitCmd(repoPath string, args ...string) {
	err := RunGitCmd(repoPath, args...)
	if err != nil {
		ExitFail("Failed to run git cmd.")
	}
}

// MustGitCommit stages changes in the dstask repository and commits them. If
// any error is encountered, the program exits.
func MustGitCommit(repoPath, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)

	// git add all changed/created files
	// could optimise this to be given an explicit list of
	// added/modified/deleted files -- only if slow.
	fmt.Printf("\n%s\n", msg)
	fmt.Printf("\033[38;5;245m")
	MustRunGitCmd(repoPath, "add", ".")

	// check for changes -- returns exit status 1 on change
	if RunGitCmd(repoPath, "diff-index", "--quiet", "HEAD", "--") == nil {
		fmt.Println("No changes detected")
		return
	}

	MustRunGitCmd(repoPath, "commit", "--no-gpg-sign", "-m", msg)
	fmt.Printf("\033[0m")
}

// MustGetRepoPath returns the full path to a file within the dstask git repo.
// Pass file as an empty string to return the git repo directory itself.
func MustGetRepoPath(repoPath, directory, file string) string {
	dir := path.Join(repoPath, directory)

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
	_, err := exec.LookPath("git")
	if err != nil {
		ExitFail("git required, please install")
	}

	gitDotGitLocation := path.Join(repoPath, ".git")

	if _, err := os.Stat(gitDotGitLocation); os.IsNotExist(err) {
		ConfirmOrAbort("Could not find dstask repository at ~/.dstask -- create?")
		MustRunGitCmd(repoPath, "init")
		fmt.Println("\nAdd a remote repository with:\n\n\tdstask git remote add origin <repo>\n")
	}
}

// Sync performs a git pull, and then a git push. If any conflicts are encountered,
// the user will need to resolve them.
func Sync(repoPath string) {
	MustRunGitCmd(repoPath, "pull", "--no-rebase", "--no-edit", "--commit", "origin", "master")
	MustRunGitCmd(repoPath, "push", "origin", "master")
}
