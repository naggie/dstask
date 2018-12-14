package dstask

import (
	"fmt"
	"github.com/gofrs/uuid"
	"os"
	"os/user"
	"path"
	"encoding/gob"
	"strings"
	"strconv"
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

func MustWriteGob(filePath string,object interface{}) {
	file, err := os.Create(filePath)
	defer file.Close()

	if err != nil {
		ExitFail("Failed to open file for writing: "+filePath)
	}

	encoder := gob.NewEncoder(file)
	encoder.Encode(object)
}

func MustReadGob(filePath string,object interface{}) {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		ExitFail("Failed to open file for reading: "+filePath)
	}

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(object)

	if err != nil {
		ExitFail("Failed to parse gob: "+filePath)
	}
}

func Help() {
	fmt.Fprintf(os.Stderr, `
Usage: task add <filter> <summary>
Example: task add +work Fix CI building P2

	Add a task with a summary and context. Current global context is
	added.


Usage: task <id>

	Show detailed information about a task


Usage: task
Usage: task <filter>
Example: task P1

	List available tasks.


Usage: task context <context>
Usage: task context none
Example: task context project:dstask
Example: task context +work +bug

	Set (or clear) a global context for all queries and inserts.



Usage: taskwarrior export | task import-from-taskwarrior
	Import tasks from taskwarrior. Note that existing tasks will not be
	updated. This is to avoid dealing with conflicts.


Usage: task help

	Show this help dialog


Usage: task modify <id> <attributes...>

Usage: task edit <id>

Usage: task describe <id>


Usage: task week

	Show tasks completed in the last week, rolling


Usage: task day

	Show tasks completed since midnight


Usage: task projects

	List project status (percentage done, estimated completion time)
`)
}

func SumInts(vals ...int) int {
	var total int

	for _, v := range(vals) {
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
