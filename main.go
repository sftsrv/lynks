package main

import (
	"fmt"
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

	ui.Run(config, files)
}
