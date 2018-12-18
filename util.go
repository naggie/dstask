package dstask

import (
	"fmt"
	"github.com/gofrs/uuid"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strconv"
	"strings"
)

func ExitFail(msg string) {
	fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", msg)
	os.Exit(1)
}

func MustExpandHome(filepath string) string {
	if strings.HasPrefix(filepath, "~/") {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		return path.Join(usr.HomeDir, filepath[2:])
	} else {
		return filepath
	}
}

func MustGetUuid4String() string {
	// does not match docs...
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return u.String()
}

func IsValidUuid4String(str string) bool {
	_, err := uuid.FromString(str)
	return err == nil
}

func IsValidPriority(priority string) bool {
	return map[string]bool{
		PRIORITY_CRITICAL: true,
		PRIORITY_HIGH:     true,
		PRIORITY_NORMAL:   true,
		PRIORITY_LOW:      true,
	}[priority]
}

func SumInts(vals ...int) int {
	var total int

	for _, v := range vals {
		total += v
	}

	return total
}

func FixStr(text string, width int) string {
	if len(text) <= width {
		return fmt.Sprintf("%-"+strconv.Itoa(width)+"v", text)
	} else {
		return text[:width]
	}
}

func MustRunGitCmd(args ...string) {
	root := MustExpandHome(GIT_REPO)
	args = append([]string{"-C", root}, args...)
	out, err := exec.Command("git", args...).CombinedOutput()

	fmt.Printf(string(out))
	if err != nil {
		ExitFail("Git command failed")
	}

}
