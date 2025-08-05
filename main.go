package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
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

type link struct {
	name, url string
}

type file struct {
	path, contents string
	links          []link
}

// TODO: we need to have some kind of state of selectfile/viewlinks/fixlinks/savelinks
type model struct {
	window
	file
	filepicker picker.Model
	linkpicker picker.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (w *window) updateWindowSize(width int, height int) {
	w.width = width
	w.height = height
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	selected := m.filepicker.GetSelected()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.window.updateWindowSize(msg.Width, msg.Height)
		m.filepicker = m.filepicker.Height(msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			if selected != "" {
				m.filepicker = m.filepicker.ClearSelected()
				return m, nil
			}
		}

		if selected == "" {
			var cmd tea.Cmd
			m.filepicker, cmd = m.filepicker.Update(msg)
			return m, cmd
		}

		// for some reason this seems to be one update late, need
		// to figure out why that's happening
		m.file = readFile(selected)

	}

	return m, nil
}

func readFile(path string) file {
	buf, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	contents := string(buf)
	linkRe := regexp.MustCompile(`\[.+?\]\(.+?\)`)
	nameRe := regexp.MustCompile(`\[.+?\]`)
	urlRe := regexp.MustCompile(`\(.+?\)`)

	matches := linkRe.FindAllString(contents, -1)
	links := []link{}

	for _, match := range matches {

		name := nameRe.FindString(match)
		url := urlRe.FindString(match)

		if name != "" && url != "" {
			links = append(links,
				link{
					name: name[1 : len(name)-1],
					url:  url[1 : len(url)-1],
				},
			)
		}

	}

	// TODO: what do we do with this once we have it?
	return file{
		path,
		contents,
		links,
	}
}

func (m model) pickerView() string {
	return m.filepicker.View()
}

func (m model) fixLinksView() string {
	selected := m.filepicker.GetSelected()
	header := theme.Heading.Render("Selected file") + theme.Primary.MarginLeft(1).Render(selected)

	links := fmt.Sprintf("links: %v", m.links)

	return lg.JoinVertical(
		lg.Top,
		header,
		links,
	)
}

func (m model) View() string {
	selected := m.filepicker.GetSelected()
	if selected == "" {
		return m.pickerView()
	}

	return m.fixLinksView()
}

func initialModel(files []string) model {
	return model{
		filepicker: picker.New().Title("File to check").Items(files),
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
