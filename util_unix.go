//go:build !windows

package dstask

import (
	"os"

	"golang.org/x/sys/unix"
)

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
