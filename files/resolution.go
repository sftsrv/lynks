package files

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/sftsrv/lynks/config"
)

type ResolutionStrategy struct {
	toMarkdownLink func(config config.Resolution, from string, to string, toAlias string) string
	// fromMarkdownLink func(config config.Resolution, from string, to string, toAlias string) string
}

const mdExtension = ".md"

var rootResolutionStrategy = ResolutionStrategy{
	toMarkdownLink: func(config config.Resolution, _ string, to string, toAlias string) string {
		hasAlias := toAlias != to
		if hasAlias {
			ext := path.Ext(toAlias)
			if config.KeepExtension {
				return toAlias
			}

			return strings.TrimSuffix(toAlias, ext)
		}

		ext := path.Ext(to)
		if config.KeepExtension {
			return to
		}

		return strings.TrimSuffix(to, ext)
	},
}

var relativeResolutionStrategy = ResolutionStrategy{
	toMarkdownLink: func(config config.Resolution, from string, to string, toAlias string) string {
		hasAlias := toAlias != to
		if hasAlias {
			ext := path.Ext(toAlias)
			if config.KeepExtension {
				return toAlias
			}

			return strings.TrimSuffix(toAlias, ext)
		}

		dir := path.Dir(from)

		absDir, dirErr := filepath.Abs(dir)
		absTo, toErr := filepath.Abs(to)

		if dirErr != nil || toErr != nil {
			panic(fmt.Errorf("Received incompatible paths. Link from %s to %s", from, to))
		}

		rel, err := filepath.Rel(absDir, absTo)
		if err != nil {
			panic(fmt.Errorf("Received incompatible paths. Link from %s to %s", from, to))
		}

		if config.KeepExtension {
			return rel
		}

		ext := path.Ext(rel)
		return strings.TrimSuffix(rel, ext)
	},
}

var resolutionStrategies = map[config.ResolutionStrategy]ResolutionStrategy{
	config.RootResolutionStrategy:     rootResolutionStrategy,
	config.RelativeResolutionStrategy: relativeResolutionStrategy,
}
