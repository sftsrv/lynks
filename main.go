package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/sftsrv/lynks/picker"
	"github.com/sftsrv/lynks/theme"
)

type window struct {
	width  int
	height int
}

type model struct {
	window
	filepicker picker.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (w *window) updateWindowSize(width int, height int) {
	w.width = width
	w.height = height
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.window.updateWindowSize(msg.Width, msg.Height)
		m.filepicker = m.filepicker.Height(msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	}

	m.filepicker, cmd = m.filepicker.Update(msg)

	return m, cmd
}

func (m model) pickerView() string {
	return m.filepicker.View()
}

func (m model) View() string {
	selected := m.filepicker.GetSelected()
	if selected == "" {
		return m.pickerView()
	}

	return lg.JoinVertical(
		lg.Top,
		theme.Heading.Render("Selected file"),
		theme.Primary.MarginLeft(1).Render(selected),
	)

}

func initialModel(files []string) model {
	return model{
		filepicker: picker.New().Title("Select a file").Items(files),
	}
}

func main() {
	files := getMarkdownFiles()

	m := initialModel(files)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func getMarkdownFiles() []string {
	var files []string

	filepath.WalkDir(".",
		func(s string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && strings.HasSuffix(s, ".md") {
				files = append(files, ""+s)
			}

			return nil
		},
	)

	return files
}
