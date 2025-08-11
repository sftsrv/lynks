package cli

import (
	"fmt"
	"os"

	"github.com/sftsrv/lynks/config"
	"github.com/sftsrv/lynks/files"
	"github.com/sftsrv/lynks/theme"

	lg "github.com/charmbracelet/lipgloss"
)

func Lint(config config.Config, paths []files.RelativePath) {
	fileCount := len(paths)
	linkCount := 0
	unresolvedCount := 0

	for _, path := range paths {
		file, links := files.ReadFile(config, path)
		linkCount += len(links)

		if !file.HasUnresolvedLinks {
			continue
		}

		fmt.Println(theme.Heading.Render(string(file.Path)))
		fmt.Println(theme.Faded.Render("Unresolved links:"))
		for _, link := range links {
			if link.IsUnresolved() {
				unresolvedCount++
				fmt.Println(theme.Warn.PaddingLeft(2).Render(link.Title()))
			}
		}
	}

	result := theme.Heading.Render("No unresolved links!")
	if unresolvedCount > 0 {
		result = theme.Alert.Render("Found unresolved links")
	}

	fmt.Println(
		lg.NewStyle().Padding(1, 2).Border(lg.NormalBorder()).Render(
			lg.JoinVertical(lg.Top,
				theme.Heading.Render("Summary"),
				theme.Primary.Render(fmt.Sprintf("%d files checked", fileCount)),
				theme.Primary.Render(fmt.Sprintf("%d links checked", linkCount)),
				theme.Primary.Render(fmt.Sprintf("%d unresolved links found", unresolvedCount)),
				lg.NewStyle().MarginTop(1).Render(result),
			),
		))

	if unresolvedCount > 0 {
		os.Exit(1)
	}

	os.Exit(0)
}
