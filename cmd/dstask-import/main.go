package main

import (
	"os"

	"github.com/naggie/dstask/pkg/imp/config"
	"github.com/naggie/dstask/pkg/imp/github"
	"github.com/sirupsen/logrus"
)

// getEnv returns an env var's value, or a default.
func getEnv(key string, _default string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return _default
}

func main() {

	repo := getEnv("DSTASK_GIT_REPO", os.ExpandEnv("$HOME/.dstask"))
	configFile := os.ExpandEnv("$HOME/.dstask-import.toml")

	cfg, err := config.Load(configFile, repo)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	err = github.Do(repo, cfg)
	if err != nil {
		logrus.Fatal(err.Error())
	}
}
