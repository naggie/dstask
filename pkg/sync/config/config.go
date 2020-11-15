package config

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Github []Github
}
type Github struct {
	Token     string
	User      string
	Repo      string
	GetClosed bool `toml:"get_closed"`
	Assignee  string
}

func Load(path string) (Config, error) {
	var config Config

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("Couldn't read config file %q: %s", path, err.Error())
	}

	_, err = toml.Decode(string(contents), &config)
	if err != nil {
		return config, fmt.Errorf("Invalid config file %q: %s", path, err.Error())
	}
	return config, nil
}
