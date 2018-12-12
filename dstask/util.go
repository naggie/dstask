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

func Help() {
	fmt.Fprintf(os.Stderr, `Usage: task add <context> <summary>

			Add a task with a summary and context. Current global context is
			added.


		Usage: task
		Usage: task next

		    List available tasks.


		Usage: task context <context>
		Usage: task context none

			Set a global context for all queries and inserts.


		Usage: taskwarrior export | task import-from-taskwarrior


		Usage: task help

			Show this help dialog
`)
}
