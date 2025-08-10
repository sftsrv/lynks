package config

import (
	"encoding/json"
	"os"
	"path"
	"strings"
)

type aliases = map[string]string

type Config struct {
	Root    string  `json:"root"`
	Aliases aliases `json:"aliases"`
}

func (c Config) AddAlias(link string) string {
	for alias, actual := range c.Aliases {
		fullPrefix := path.Join(c.Root, actual)

		if after, ok := strings.CutPrefix(link, fullPrefix); ok {
			return alias + after
		}
	}

	return link
}

func (c Config) RemoveAlias(link string) string {
	for alias, actual := range c.Aliases {
		if after, ok := strings.CutPrefix(link, alias); ok {
			return path.Join(c.Root, actual, after)
		}
	}

	return link
}

func Load(path string) (Config, error) {
	config := Config{}
	file, fileErr := os.ReadFile(path)

	if fileErr != nil {
		return config, fileErr
	}

	jsonErr := json.Unmarshal(file, &config)
	if jsonErr != nil {
		return config, jsonErr
	}

	return config, nil
}
