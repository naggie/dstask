package dstask

import (
	"fmt"
	"github.com/gofrs/uuid"
	"os"
	"io/ioutil"
	"os/exec"
	"os/user"
	"path"
	"strconv"
	"strings"
	"encoding/gob"
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

	err := cmd.Run()

	if err != nil {
		ExitFail("%s cmd failed", name)
	}
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
	for _, item := range(haystack) {
		if item == needle {
			return true
		}
	}

	return false
}

func MustWriteGob(filePath string,object interface{}) {
	file, err := os.Create(filePath)
	defer file.Close()

	if err != nil {
		ExitFail("Failed to open %s for writing: ", filePath)
	}

	encoder := gob.NewEncoder(file)
	encoder.Encode(object)
}

func MustReadGob(filePath string,object interface{}) {
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

