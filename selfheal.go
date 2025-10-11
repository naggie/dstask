package dstask

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	shortcutDir = ".local/bin"
)

var shortcuts = []string{"p0", "p1", "p2", "p3", "ds"}

func SelfHeal(conf Config) {
	CreateShortcuts(conf)
}

func CreateShortcuts(conf Config) {
	if !ensureShortcuts(conf) {
		createShortcuts(conf)
	}
}

func validateShortcuts(conf Config, shortcutPath string) bool {
	currentBinary, err := os.Executable()
	if err != nil {
		return false
	}

	currentBinary, err = filepath.EvalSymlinks(currentBinary)
	if err != nil {
		return false
	}

	for _, shortcut := range shortcuts {
		symlinkPath := filepath.Join(shortcutPath, shortcut)

		target, err := os.Readlink(symlinkPath)
		if err != nil {
			return false
		}

		resolvedTarget, err := filepath.EvalSymlinks(target)
		if err != nil {
			return false
		}

		if resolvedTarget != currentBinary {
			return false
		}
	}

	return true
}

func ensureShortcuts(conf Config) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	shortcutPath := filepath.Join(home, shortcutDir)

	if _, err := os.Stat(shortcutPath); os.IsNotExist(err) {
		if err := os.MkdirAll(shortcutPath, 0755); err != nil {
			return false
		}
	}

	return validateShortcuts(conf, shortcutPath)
}

func createShortcuts(conf Config) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	shortcutPath := filepath.Join(home, shortcutDir)

	if err := os.MkdirAll(shortcutPath, 0755); err != nil {
		return
	}

	currentBinary, err := os.Executable()
	if err != nil {
		return
	}

	for _, shortcut := range shortcuts {
		symlinkPath := filepath.Join(shortcutPath, shortcut)

		os.Remove(symlinkPath)

		if err := os.Symlink(currentBinary, symlinkPath); err != nil {
			continue
		}
	}
}

func pathVerification() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	shortcutPath := filepath.Join(home, shortcutDir)
	pathEnv := os.Getenv("PATH")

	for _, dir := range strings.Split(pathEnv, ":") {
		cleanDir := filepath.Clean(dir)
		if cleanDir == shortcutPath {
			return true
		}

		expandedHome := strings.Replace(dir, "~", home, 1)
		if filepath.Clean(expandedHome) == shortcutPath {
			return true
		}
	}

	return false
}
