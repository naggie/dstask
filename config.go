package dstask

import (
	"os"
	"path"
)

// Config models the dstask application's required configuration. All paths
// are absolute.
type Config struct {
	Repo      string
	StateFile string
	IDsFile   string
}

// NewConfig generates a new Config struct from the environment.
func NewConfig() Config {

	var conf Config

	repoPath := getEnv("DSTASK_GIT_REPO", os.ExpandEnv("$HOME/.dstask"))
	stateFilePath := path.Join(repoPath, ".git", "dstask", "state.bin")
	idsFilePath := path.Join(repoPath, ".git", "dstask", "ids.bin")

	conf.Repo = repoPath
	conf.StateFile = stateFilePath
	conf.IDsFile = idsFilePath

	return conf
}

// getEnv returns an env var's value, or a default.
func getEnv(key string, _default string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return _default
}
