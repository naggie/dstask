package dstask

import (
	"os"
	"path/filepath"
)

// Config models the dstask application's required configuration. All paths
// are absolute.
type Config struct {
	// Path to the git repository
	Repo string
	// Path to the dstask local state file. State will differ between machines
	StateFile string
	// Path to the ids file
	IDsFile string
	// An unparsed context string, provided via DSTASK_CONTEXT
	CtxFromEnvVar string
}

// NewConfig generates a new Config struct from the environment.
func NewConfig() Config {
	var conf Config

	conf.CtxFromEnvVar = getEnv("DSTASK_CONTEXT", "")
	// Bestimme Home-Verzeichnis plattformunabhängig
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback: benutze $HOME falls vorhanden
		home = os.Getenv("HOME")
	}
	defaultRepo := filepath.Join(home, ".dstask")
	conf.Repo = getEnv("DSTASK_GIT_REPO", defaultRepo)
	conf.StateFile = filepath.Join(conf.Repo, ".git", "dstask", "state.bin")
	conf.IDsFile = filepath.Join(conf.Repo, ".git", "dstask", "ids.bin")

	return conf
}

// getEnv returns an env var's value, or a default.
func getEnv(key string, _default string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return _default
}
