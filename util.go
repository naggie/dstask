package dstask

import (
	"encoding/gob"
	"fmt"
	"github.com/gofrs/uuid"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strconv"
	"strings"
)

func ExitFail(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "\033[31m"+format+"\033[0m\n", a...)
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

func SumInts(vals ...int) int {
	var total int

	for _, v := range vals {
		total += v
	}

	return total
}

func FixStr(text string, width int) string {
	// remove after newline
	text = strings.Split(text, "\n")[0]
	if len(text) <= width {
		return fmt.Sprintf("%-"+strconv.Itoa(width)+"v", text)
	} else {
		return text[:width]
	}
}

func MustRunCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func MustRunGitCmd(args ...string) {
	root := MustExpandHome(GIT_REPO)
	args = append([]string{"-C", root}, args...)
	MustRunCmd("git", args...)
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

	MustRunCmd(editor, tmpfile.Name())
	data, err = ioutil.ReadFile(tmpfile.Name())

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

func MustWriteGob(filePath string, object interface{}) {
	file, err := os.Create(filePath)
	defer file.Close()

	if err != nil {
		ExitFail("Failed to open %s for writing: ", filePath)
	}

	encoder := gob.NewEncoder(file)
	encoder.Encode(object)
}

func MustReadGob(filePath string, object interface{}) {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		ExitFail("Failed to open %s for reading: ", filePath)
	}

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(object)

	if err != nil {
		ExitFail("Failed to parse gob: %s", filePath)
	}
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

func MustGetTermSize() (int,int) {
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
