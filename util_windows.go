//go:build windows

package dstask

import (
	"os"

	"github.com/mattn/go-isatty"
	"golang.org/x/sys/windows"
)

func MustGetTermSize() (int, int) {
	if FAKE_PTY {
		return 80, 24
	}

	fd := os.Stdout.Fd()

	// Fallback: if not a TTY, fail as before
	if !(isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)) {
		ExitFail("Not a TTY")
	}
	
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(windows.Handle(fd), &info); err != nil {
		return 80, 24
	}

	return int(info.Window.Right - info.Window.Left + 1), int(info.Window.Bottom - info.Window.Top + 1)
}
