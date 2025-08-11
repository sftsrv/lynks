package main

import (
	"os"

	"github.com/sftsrv/lynks/cli"
	"github.com/sftsrv/lynks/config"
	"github.com/sftsrv/lynks/files"
	"github.com/sftsrv/lynks/ui"
)

func main() {
	configPath := "lynks.config.json"
	config := config.Load(configPath)

	files := files.GetMarkdownFiles(config)

	if len(os.Args) < 2 {
		ui.Run(config, files)
		return
	}

	switch os.Args[1] {
	case "lint":
		cli.Lint(config, files)
	}
}
