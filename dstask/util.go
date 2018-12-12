package dstask

import (
	"fmt"
	"github.com/gofrs/uuid"
	"os"
	"os/user"
	"path"
	"strings"
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
		PRIORITY_HIGH: true,
		PRIORITY_NORMAL: true,
		PRIORITY_LOW: true,
	}[priority]
}

func Help() {
	fmt.Fprintf(os.Stderr, `
Usage: task add <filter> <summary>

	Add a task with a summary and context. Current global context is
	added.


Usage: task
Usage: task <filter>
Usage: task next
Usage: task next <filter>

	List available tasks.


Usage: task context <context>
Usage: task context none

	Set a global context for all queries and inserts.


Usage: taskwarrior export | task import-from-taskwarrior


Usage: task help

	Show this help dialog


Usage: task modify <id> <attributes...>

Usage: task edit <id>

Usage: task describe <id>
`)
}
