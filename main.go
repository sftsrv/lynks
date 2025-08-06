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

type linkStatus int

const (
	resolved linkStatus = iota
	unresolved
	remote
)

type link struct {
	name, url, absolute string
	status              linkStatus
}

func (l link) icon() string {
	switch l.status {
	case remote:
		return "www"
	case resolved:
		return "loc"
	case unresolved:
		return "xxx"
	}

	return "???"
}

func (l link) Title() string {
	return l.icon() + " " + lg.NewStyle().Bold(true).Render(l.name) + " " + l.absolute + " (" + l.url + ")"
}

type file struct {
	path
	contents string
}

type path string

func (s path) Title() string {
	return string(s)
}

type state int

const (
	filePicker state = iota
	linkPickerView
)

// TODO: we need to have some kind of state of selectfile/viewlinks/fixlinks/savelinks
type model struct {
	state state
	base  string
	window
	file
	filepicker picker.Model[path]
	linkpicker picker.Model[link]
}

func (m model) Init() tea.Cmd {
	return nil
}

func (w *window) updateWindowSize(width int, height int) {
	w.width = width
	w.height = height
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.window.updateWindowSize(msg.Width, msg.Height)
		m.filepicker = m.filepicker.Height(msg.Height)
		m.linkpicker = m.linkpicker.Height(msg.Height - 1)
		return m, nil

	case picker.SelectedMsg[path]:
		file, links := readFile(m.base, msg.Selected)

		m.state = linkPickerView
		m.file = file
		m.linkpicker = m.linkpicker.Items(links)

	case tea.KeyMsg:
		str := msg.String()

		if str == "ctrl+c" {
			return m, tea.Quit
		}

		// if no file selected, delegate messag handling to picker
		switch m.state {
		case filePicker:
			var cmd tea.Cmd
			m.filepicker, cmd = m.filepicker.Update(msg)
			return m, cmd

		case linkPickerView:
			if str == "esc" {
				m.state = filePicker
			}

			var cmd tea.Cmd
			m.linkpicker, cmd = m.linkpicker.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// TODO: fix this function, does not work correctly for relative links
func resolveLink(base string, relative string, url string) (linkStatus, string) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return remote, url
	}

	if strings.HasPrefix(url, "/") {
		url = base + url
	}

	dir, dirErr := filepath.Rel(string(relative), "..")
	if dirErr != nil {
		return unresolved, url
	}

	absPath, absErr := filepath.Abs(filepath.Join(dir, url))
	if absErr != nil {
		return unresolved, relative
	}

	stat, statErr := os.Stat(absPath)
	if statErr != nil {
		return unresolved, absPath
	}

	if stat.IsDir() {
		return unresolved, absPath
	}

	return resolved, absPath
}

func readFile(base string, path path) (file, []link) {
	buf, err := os.ReadFile(string(path))
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

		namePart := nameRe.FindString(match)
		urlPart := urlRe.FindString(match)

		if namePart != "" && urlPart != "" {
			name := namePart[1 : len(namePart)-1]
			url := urlPart[1 : len(urlPart)-1]
			status, absolute := resolveLink(base, string(path), url)

			links = append(links,
				link{
					name,
					url,
					absolute,
					status,
				},
			)
		}

	}

	// TODO: what do we do with this once we have it?
	return file{
		path,
		contents,
	}, links
}

func (m model) filePickerView() string {
	return m.filepicker.View()
}

func (m model) linkPickerView() string {
	selected := m.file.path
	header := theme.Heading.Render("Links for") + theme.Primary.MarginLeft(1).Render(string(selected))

	return lg.JoinVertical(
		lg.Top,
		header,
		m.linkpicker.View(),
	)
}

func (m model) View() string {
	switch m.state {
	case filePicker:
		return m.filePickerView()

	case linkPickerView:
		return m.linkPickerView()
	}

	return "unexpected state"
}

func initialModel(files []path) model {
	return model{
		base:       ".",
		state:      filePicker,
		filepicker: picker.New[path]().Title("File to check").Accent(theme.ColorPrimary).Items(files),
		linkpicker: picker.New[link]().Title("Edit Link").Accent(theme.ColorSecondary),
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

func getMarkdownFiles() []path {
	var files []path

	filepath.WalkDir(".",
		func(s string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && strings.HasSuffix(s, ".md") {
				files = append(files, path(s))
			}

			return nil
		},
	)

	return files
}
