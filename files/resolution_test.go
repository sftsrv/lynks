package files

import (
	"testing"

	"github.com/sftsrv/lynks/config"
)

type Case struct {
	to       string
	toAlias  string
	expected string
}

func TestRootResolution(t *testing.T) {
	// rootResolutionStrategy
}

func TestRelativeResolutionStrategyWithoutExtension(t *testing.T) {
	config := config.Resolution{
		Strategy:      config.RelativeResolutionStrategy,
		KeepExtension: false,
	}

	toMd := relativeResolutionStrategy.toMarkdownLink

	from := "my-example/folder/file.md"
	cases := []Case{
		{"my-example/uncle/cousin.md", "my-example/uncle/cousin.md", "../uncle/cousin"},
		{"my-example/parent.md", "my-example/parent.md", "../parent"},
		{"my-example/folder/sibling.md", "my-example/folder/sibling.md", "sibling"},

		// alias
		{"my-example/folder/sibling.md", "my-alias/sibling.md", "my-alias/sibling"},
	}

	for _, c := range cases {
		result := toMd(config, from, c.to, c.toAlias)
		if result != c.expected {
			t.Errorf("\ngiven %v\ngot %v\nexpected %v", c, result, c.expected)
		}
	}
}

func TestRelativeResolutionStrategyWitExtension(t *testing.T) {
	config := config.Resolution{
		Strategy:      config.RelativeResolutionStrategy,
		KeepExtension: true,
	}

	toMd := relativeResolutionStrategy.toMarkdownLink

	from := "my-example/folder/file.md"
	cases := []Case{
		{"my-example/uncle/cousin.md", "my-example/uncle/cousin.md", "../uncle/cousin.md"},
		{"my-example/parent.md", "my-example/parent.md", "../parent.md"},
		{"my-example/folder/sibling.md", "my-example/folder/sibling.md", "sibling.md"},

		// alias
		{"my-example/folder/sibling.md", "my-alias/sibling.md", "my-alias/sibling.md"},
	}

	for _, c := range cases {
		result := toMd(config, from, c.to, c.toAlias)
		if result != c.expected {
			t.Errorf("\ngiven %v\ngot %v\nexpected %v", c, result, c.expected)
		}
	}
}

func TestRootResolutionStrategyWithoutExtension(t *testing.T) {
	config := config.Resolution{
		Strategy:      config.RootResolutionStrategy,
		KeepExtension: false,
	}

	toMd := rootResolutionStrategy.toMarkdownLink

	from := "my-example/folder/file.md"
	cases := []Case{
		{"my-example/folder/sibling.md", "my-example/folder/sibling.md", "my-example/folder/sibling"},

		// alias
		{"my-example/folder/sibling.md", "my-alias/sibling.md", "my-alias/sibling"},
	}

	for _, c := range cases {
		result := toMd(config, from, c.to, c.toAlias)
		if result != c.expected {
			t.Errorf("\ngiven %v\ngot %v\nexpected %v", c, result, c.expected)
		}
	}
}

func TestRootResolutionStrategyWitExtension(t *testing.T) {
	config := config.Resolution{
		Strategy:      config.RootResolutionStrategy,
		KeepExtension: true,
	}

	toMd := rootResolutionStrategy.toMarkdownLink

	from := "my-example/folder/file.md"
	cases := []Case{
		{"my-example/folder/sibling.md", "my-example/folder/sibling.md", "my-example/folder/sibling.md"},

		// alias
		{"my-example/folder/sibling.md", "my-alias/sibling.md", "my-alias/sibling.md"},
	}

	for _, c := range cases {
		result := toMd(config, from, c.to, c.toAlias)
		if result != c.expected {
			t.Errorf("\ngiven %v\ngot %v\nexpected %v", c, result, c.expected)
		}
	}
}
