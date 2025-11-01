//go:build windows

package dstask

import (
	"os"

	"github.com/mattn/go-isatty"
)

func MustGetTermSize() (int, int) {
	if FAKE_PTY {
		return 80, 24
	}

	// Fallback: if not a TTY, fail as vorher
	if !(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())) {
		ExitFail("Not a TTY")
	}

	// On Windows, golang.org/x/sys/unix is unavailable; use conservative default
	// Many Windows terminals handle wrapping; we pick a reasonable width
	return 80, 24
}
