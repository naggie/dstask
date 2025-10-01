package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/naggie/dstask"
	"github.com/naggie/dstask/pkg/imp/config"
	"github.com/naggie/dstask/pkg/imp/github"
	"github.com/naggie/dstask/pkg/imp/tw"
	"github.com/sirupsen/logrus"
)

// getEnv returns an env var's value, or a default.
func getEnv(key string, _default string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return _default
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: dstask-import github|tw|--help|help")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "       dstask-import help or --help       # this menu")
	fmt.Fprintln(
		os.Stderr,
		"       dstask-import github               # import from GitHub as specified in configuration",
	)
	fmt.Fprintln(
		os.Stderr,
		"       cat export.json | dstask-import tw # import from a taskwarrior json dump which can be obtained with the taskwarrior command 'task export'",
	)
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "--help", "help":
		usage()
	case "tw":
		conf := dstask.NewConfig()
		if err := tw.Do(conf); err != nil {
			dstask.ExitFail(err.Error())
		}
	case "github":
		// Plattformsichere Standardpfade bestimmen
		home, err := os.UserHomeDir()
		if err != nil {
			home = os.Getenv("HOME")
		}
		repo := getEnv("DSTASK_GIT_REPO", filepath.Join(home, ".dstask"))
		configFile := filepath.Join(home, ".dstask-import.toml")

		cfg, err := config.Load(configFile, repo)
		if err != nil {
			logrus.Fatal(err.Error())
		}

		err = github.Do(repo, cfg)
		if err != nil {
			logrus.Fatal(err.Error())
		}
	default:
		usage()
		os.Exit(2)
	}
}
