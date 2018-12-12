package dstask

import (
	"fmt"
	"os"
	"path"
	"strings"
	"os/user"
	"github.com/satori/go.uuid"
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
	return uuid.NewV4().String()
}

func IsValidUuid4String(str string) bool {
	_, err := uuid.FromString(str)
	return err == nil
}
