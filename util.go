package dstask

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/gofrs/uuid"
	"golang.org/x/sys/unix"
)

func ExitFail(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "\033[31m"+format+"\033[0m\n", a...)
	os.Exit(1)
}

func ConfirmOrAbort(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+" [y/n] ", a...)

	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	if input == "y\n" {
		return
	} else {
		ExitFail("Aborted.")
	}
}

func MustGetUUID4String() string {
	// does not match docs...
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return u.String()
}

func IsValidUUID4String(str string) bool {
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

func IsValidStatus(status string) bool {
	return StrSliceContains(ALL_STATUSES, status)
}

func SumInts(vals ...int) int {
	var total int

	for _, v := range vals {
		total += v
	}

	return total
}

func RunCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func MustEditBytes(data []byte, ext string) []byte {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		editor = "vim"
	}

	tmpfile, err := ioutil.TempFile("", "dstask.*."+ext)
	if err != nil {
		ExitFail("Could not create temporary file to edit")
	}
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(data)
	tmpfile.Close()

	if err != nil {
		ExitFail("Could not write to temporary file to edit")
	}

	err = RunCmd(editor, tmpfile.Name())
	if err != nil {
		ExitFail("Failed to run $EDITOR")
	}

	data, err = ioutil.ReadFile(tmpfile.Name())

	if err != nil {
		ExitFail("Could not read back temporary edited file")
	}

	return data
}

func StrSliceContains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}

func IsValidStateTransition(from string, to string) bool {
	for _, transition := range VALID_STATUS_TRANSITIONS {
		if from == transition[0] && to == transition[1] {
			return true
		}
	}

	return false
}

func MustOpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		ExitFail("unsupported platform")
	}

	if err != nil {
		ExitFail("Failed to open browser")
	}
}

func DeduplicateStrings(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func MustGetTermSize() (int, int) {
	if FAKE_PTY {
		return 80, 24
	}

	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		ExitFail("Not a TTY")
	}

	return int(ws.Col), int(ws.Row)
}

func IsTTY() bool {
	_, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	return err == nil || FAKE_PTY
}

func WriteStdout(data []byte) error {
	if _, err := os.Stdout.Write(data); err != nil {
		return err
	}
	return nil
}
