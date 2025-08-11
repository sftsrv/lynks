package main

import (
	"fmt"
	"os"

	"github.com/sftsrv/lynks/cli"
	"github.com/sftsrv/lynks/config"
	"github.com/sftsrv/lynks/files"
	"github.com/sftsrv/lynks/ui"
)

func main() {
	configPath := "lynks.config.json"
	config, configErr := config.Load(configPath)

	fmt.Printf("config %v", config)

	if configErr != nil {
		panic(configErr)
	}

	files := files.GetMarkdownFiles(config)

	if len(os.Args) < 2 {
		ui.Run(config, files)
	}
	
	switch os.Args[1] {
	case "lint":
		cli.Lint(config, files)
	}
}
