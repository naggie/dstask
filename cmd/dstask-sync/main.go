package main

import (
	"os"

	"github.com/naggie/dstask"
	"github.com/naggie/dstask/pkg/sync"
	"github.com/naggie/dstask/pkg/sync/config"
	"github.com/naggie/dstask/pkg/sync/github"
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
	configFile := os.ExpandEnv("$HOME/.dstask-sync.toml")

	cfg, err := config.Load(configFile, repo)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	for _, cfgGithub := range cfg.Github {
		if cfgGithub.Token == "" {
			logrus.Infof("GitHub config section %s/%s: skipping because no token configured", cfgGithub.User, cfgGithub.Repo)
			continue
		}
		logrus.Infof("GitHub config section %s/%s: processing", cfgGithub.User, cfgGithub.Repo)
		var src sync.Source
		src, err := github.NewClient(cfgGithub)
		if err != nil {
			logrus.Fatal(err.Error())
		}

		for {
			tasks, err := src.Next()
			if err != nil {
				logrus.Fatal(err.Error())
			}
			if len(tasks) == 0 {
				break
			}

			for _, t := range tasks {
				err = sync.ProcessTask(repo, t)
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	}
	dstask.MustGitCommit(repo, "GitHub import")
}
