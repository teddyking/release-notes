package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func MustLoad(path string) Config {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	c := Config{}
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		panic(err)
	}

	return c
}

type Config struct {
	ServerURL string    `yaml:"server_url"`
	Owner     string    `yaml:"owner"`
	Repo      string    `yaml:"repo"`
	Title     string    `yaml:"title"`
	Includes  []Include `yaml:"includes"`
}

type Include struct {
	Type           string `yaml:"type"`
	Name           string `yaml:"name"`
	Title          string `yaml:"title"`
	IncludeCommits bool   `yaml:"include_commits"`
}
