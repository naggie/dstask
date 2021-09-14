package dstask

import (
	"os"
	"path"
)

// Config models the dstask application's required configuration. All paths
// are absolute.
type Config struct {
	// Branch to use when pushing and pulling to git remote
	Branch string
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

	conf.Branch = getEnv("DSTASK_BRANCH", "master")
	conf.CtxFromEnvVar = getEnv("DSTASK_CONTEXT", "")
	conf.Repo = getEnv("DSTASK_GIT_REPO", os.ExpandEnv("$HOME/.dstask"))
	conf.StateFile = path.Join(conf.Repo, ".git", "dstask", "state.bin")
	conf.IDsFile = path.Join(conf.Repo, ".git", "dstask", "ids.bin")

	return conf
}

// getEnv returns an env var's value, or a default.
func getEnv(key string, _default string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return _default
}
