package config

import (
	"fmt"
	"io/ioutil"
	"path"

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
	Template     string
	TemplateTask dstask.Task `toml:"-"`
}

func Load(configFile, repo string) (Config, error) {
	var config Config

	contents, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, fmt.Errorf("Couldn't read config file %q: %s", configFile, err.Error())
	}

	_, err = toml.Decode(string(contents), &config)
	if err != nil {
		return config, fmt.Errorf("Invalid config file %q: %s", configFile, err.Error())
	}

	for i, gh := range config.Github {

		tplFile := path.Join(repo, "templates-github", gh.Template+".yml")

		data, err := ioutil.ReadFile(tplFile)
		if err != nil {
			return config, fmt.Errorf("Failed to read %s: %s", tplFile, err.Error())
		}
		err = yaml.Unmarshal(data, &config.Github[i].TemplateTask)
		if err != nil {
			return config, fmt.Errorf("Failed to unmarshal %s: %s", tplFile, err.Error())
		}
	}

	return config, nil
}
