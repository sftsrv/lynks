package files

import (
	"path"
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

var relativeResolutionStrategy = ResolutionStrategy{}

var resolutionStrategies = map[config.ResolutionStrategy]ResolutionStrategy{
	config.RootResolutionStrategy:     rootResolutionStrategy,
	config.RelativeResolutionStrategy: relativeResolutionStrategy,
}
