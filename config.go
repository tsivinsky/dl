package main

import (
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type App struct {
	Name        string   `yaml:"name"`
	URL         string   `yaml:"url"`
	Build       []string `yaml:"build"`
	Destination string   `yaml:"dest"`
}

type Conf struct {
	DL []App `yaml:"dl"`
}

func setupConfig() (string, string, error) {
	userConfDir, err := os.UserConfigDir()
	if err != nil {
		return "", "", err
	}

	confDir := path.Join(userConfDir, "dl")
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		err = os.Mkdir(confDir, 0755)
		if err != nil {
			return "", "", err
		}
	}

	confFile := path.Join(confDir, "config.yml")
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		_, err = os.OpenFile(confFile, os.O_CREATE, 0644)
		if err != nil {
			return "", "", err
		}
	}

	return confDir, confFile, nil
}

func parseConfig(configPath string) (Conf, error) {
	var conf Conf

	f, err := os.OpenFile(configPath, os.O_RDONLY, 0644)
	if err != nil {
		return conf, err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&conf); err != nil {
		return conf, err
	}

	return conf, nil
}
