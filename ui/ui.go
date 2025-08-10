package ui

import (
	"fmt"
	"os"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/sftsrv/lynks/config"
	"github.com/sftsrv/lynks/files"
	"github.com/sftsrv/lynks/picker"
	"github.com/sftsrv/lynks/theme"
)

type window struct {
	width  int
	height int
}

type state int

const (
	filePickerView state = iota
	linkPickerView
	linkFixerView
)

type Model struct {
	config config.Config

	state  state
	window window

	file       files.File
	filepicker picker.Model[files.RelativePath]

	link       files.Link
	linkpicker picker.Model[files.Link]
	linkfixer  picker.Model[files.RelativePath]
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (w *window) updateWindowSize(width int, height int) {
	w.width = width
	w.height = height
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.window.updateWindowSize(msg.Width, msg.Height)
		m.filepicker = m.filepicker.Height(msg.Height)
		m.linkpicker = m.linkpicker.Height(msg.Height - 1)
		m.linkfixer = m.linkfixer.Height(msg.Height - 3)
		return m, nil

	case picker.SelectedMsg[files.RelativePath]:
		switch m.state {
		case filePickerView:
			file, links := readFile(m.config, msg.Selected)

			m.state = linkPickerView
			m.file = file
			m.linkpicker = m.linkpicker.Items(links)

		case linkFixerView:
			m.state = linkPickerView
			updated := files.FixLink(m.config, m.file, m.link, msg.Selected)
			updateFile(updated)
			file, links := readFile(m.config, updated.Path)

			m.file = file
			m.linkpicker = m.linkpicker.Items(links)
		}

	case picker.SelectedMsg[files.Link]:
		m.state = linkFixerView

		m.linkfixer = m.linkfixer.Search(msg.Selected.FileName())
		m.link = msg.Selected

	case tea.KeyMsg:
		str := msg.String()

		if str == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.state {
		case filePickerView:
			var cmd tea.Cmd
			m.filepicker, cmd = m.filepicker.Update(msg)
			return m, cmd

		case linkPickerView:
			if str == "esc" {
				m.state = filePickerView
			}

			var cmd tea.Cmd
			m.linkpicker, cmd = m.linkpicker.Update(msg)
			return m, cmd

		case linkFixerView:
			if str == "esc" {
				m.state = linkPickerView
			}

			var cmd tea.Cmd
			m.linkfixer, cmd = m.linkfixer.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func updateFile(file files.File) {
	osFile, err := os.Create(string(file.Path))
	if err != nil {
		panic(fmt.Errorf("Failed to open file: %v", err))
	}

	_, err = osFile.WriteString(file.Contents)
	if err != nil {
		panic(fmt.Errorf("Failed to update file: %v", err))
	}

	err = osFile.Close()
	if err != nil {
		panic(fmt.Errorf("Failed to close file: %v", err))
	}
}

func readFile(config config.Config, path files.RelativePath) (files.File, []files.Link) {
	buf, err := os.ReadFile(string(path))
	if err != nil {
		panic(err)
	}

	contents := string(buf)
	linkRe := regexp.MustCompile(`\[.+?\]\(.+?\)`)
	nameRe := regexp.MustCompile(`\[.+?\]`)
	urlRe := regexp.MustCompile(`\(.+?\)`)

	matches := linkRe.FindAllString(contents, -1)
	links := []files.Link{}

	for _, match := range matches {
		namePart := nameRe.FindString(match)
		urlPart := urlRe.FindString(match)

		if namePart != "" && urlPart != "" {
			name := namePart[1 : len(namePart)-1]
			url := urlPart[1 : len(urlPart)-1]
			status, resolved := files.ResolveLink(config, string(path), url)

			links = append(links,
				files.Link{Name: name, Url: url, Resolved: resolved, Status: status},
			)
		}

	}

	hasLinks := len(links) > 0

	return files.File{Path: path, Contents: contents, HasLinks: hasLinks}, links
}

func (m Model) filePickerView() string {
	return m.filepicker.View()
}

func (m Model) linkPickerView() string {

	selected := m.file.Path
	header := theme.Heading.Render("Links for") + theme.Primary.MarginLeft(1).Render(string(selected))

	noLinksMessage := theme.Faded.Render("No links found in file")
	exitMessage := theme.Faded.Render("<esc> to go back to files")

	if !m.file.HasLinks {
		return lg.JoinVertical(lg.Top, header, noLinksMessage, exitMessage)
	}

	return lg.JoinVertical(
		lg.Top,
		header,
		m.linkpicker.View(),
	)
}

func (m Model) linkFixerView() string {
	selected := m.file.Path
	header := lg.JoinVertical(
		lg.Top,
		theme.Heading.Render("Fix links for")+theme.Primary.MarginLeft(1).Render(string(selected)),
		lg.NewStyle().MarginLeft(4).Foreground(theme.ColorSecondary).Render(m.link.Title()),
	)

	return lg.JoinVertical(
		lg.Top,
		header,
		m.linkfixer.View(),
	)
}

func (m Model) View() string {
	switch m.state {
	case filePickerView:
		return m.filePickerView()

	case linkPickerView:
		return m.linkPickerView()

	case linkFixerView:
		return m.linkFixerView()
	}

	return "unexpected state"
}

func initialModel(config config.Config, f []files.RelativePath) Model {
	return Model{
		config:     config,
		state:      filePickerView,
		filepicker: picker.New[files.RelativePath]().Title("File to check").Accent(theme.ColorPrimary).Items(f),
		linkpicker: picker.New[files.Link]().Title("Edit Link").Accent(theme.ColorSecondary),
		linkfixer:  picker.New[files.RelativePath]().Title("Fix link").Accent(theme.ColorWarn).Items(f),
	}
}

func Run(config config.Config, f []files.RelativePath) {
	m := initialModel(config, f)

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
