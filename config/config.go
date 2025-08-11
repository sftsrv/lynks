package config

import (
	"encoding/json"
	"os"
	"path"
	"strings"
)

type aliases = map[string]string

type ResolutionStrategy string

const (
	RootResolutionStrategy     ResolutionStrategy = "root"
	RelativeResolutionStrategy ResolutionStrategy = "relative"
)

type Resolution struct {
	Strategy      ResolutionStrategy `json:"strategy"`
	KeepExtension bool               `json:"keepExtension"`
}

type Config struct {
	Root       string     `json:"root"`
	Resolution Resolution `json:"resolution"`
	Ignore     []string   `json:"ignore"`
	Aliases    aliases    `json:"aliases"`
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

func defaultConfig() Config {
	return Config{
		Root: "./",
		Resolution: Resolution{
			Strategy:      RelativeResolutionStrategy,
			KeepExtension: true,
		},
	}
}

func Load(path string) Config {
	config := defaultConfig()

	file, fileErr := os.ReadFile(path)
	if fileErr != nil {
		return config
	}

	jsonErr := json.Unmarshal(file, &config)
	if jsonErr != nil {
		return config
	}

	return config
}
