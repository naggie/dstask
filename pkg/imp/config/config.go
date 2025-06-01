package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/naggie/dstask"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Github []Github
}

type Github struct {
	Token        string
	Repos        []string
	GetClosed    bool `toml:"get_closed"`
	Assignee     string
	Milestone    string
	Labels       []string
	TemplateStr  string      `toml:"template_str"`
	TemplateTask dstask.Task `toml:"-"`
}

func Load(configFile, repo string) (Config, error) {
	var config Config

	contents, err := os.ReadFile(configFile)
	if err != nil {
		return config, fmt.Errorf("couldn't read config file %q: %s", configFile, err.Error())
	}

	_, err = toml.Decode(string(contents), &config)
	if err != nil {
		return config, fmt.Errorf("invalid config file %q: %s", configFile, err.Error())
	}

	for i, gh := range config.Github {
		err = yaml.Unmarshal([]byte(gh.TemplateStr), &config.Github[i].TemplateTask)
		if err != nil {
			return config, fmt.Errorf("failed to unmarshal template: %s", err.Error())
		}
	}

	return config, nil
}
