package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type filepicker struct {
	search string
	cursor int
	files  []string
}

type selected struct {
	path string
}

type model struct {
	filepicker
	selected
}

func (m model) Init() tea.Cmd {
	return nil
}

func clamp(i int, min int, max int) int {
	if i > max {
		return max
	}

	if i < min {
		return min
	}

	return i
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	max := len(m.files) - 1

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left", "h":
			m.cursor = 0
		case "right", "l":
			m.cursor = max

		case "up", "k":
			m.cursor--
			m.cursor = clamp(m.cursor, 0, max)
		case "down", "j":
			m.cursor++
			m.cursor = clamp(m.cursor, 0, max)
		}

	}

	return m, nil
}

func (m model) View() string {

	var view strings.Builder
	view.WriteString("Select a file")

	for i, file := range m.files {
		view.WriteString("\n")
		if i == m.cursor {
			view.WriteString("> ")
		} else {
			view.WriteString("  ")
		}
		view.WriteString(file)
	}

	view.WriteString("\n")

	return view.String()
}

func main() {
	var files []string

	filepath.WalkDir(".",
		func(s string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && strings.HasSuffix(s, ".md") {
				files = append(files, s)
			}

			return nil
		},
	)

	m := model{
		filepicker: filepicker{
			search: "",
			files:  files,
		},
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
