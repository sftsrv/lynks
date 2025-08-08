package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
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
	name string
	url  string

	resolved relativePath
	status   linkStatus
}

type Config struct {
	resolveExtension string
	aliases          []alias
}

func (l link) color() lg.Color {
	switch l.status {
	case remote:
		return theme.ColorSecondary
	case resolved:
		return theme.ColorSecondary
	case unresolved:
		return theme.ColorWarn
	}

	return theme.ColorError
}

type alias struct {
	alias, actual string
}

func (c Config) addAlias(link string) string {
	for _, alias := range c.aliases {
		if after, ok := strings.CutPrefix(link, alias.actual); ok {
			return alias.alias + after
		}
	}

	return link
}

func (c Config) removeAlias(link string) string {
	for _, alias := range c.aliases {
		if after, ok := strings.CutPrefix(link, alias.alias); ok {
			return alias.actual + after
		}
	}

	return link
}

func (l link) Title() string {
	return lg.NewStyle().Foreground(l.color()).Render(lg.NewStyle().Bold(true).Render(l.name) + " " + string(l.resolved) + " (" + l.url + ")")
}

type file struct {
	path     relativePath
	contents string
}

type relativePath string

func (s relativePath) Title() string {
	return string(s)
}

type state int

const (
	filePickerView state = iota
	linkPickerView
	linkFixerView
)

type model struct {
	config Config

	state  state
	window window

	file       file
	filepicker picker.Model[relativePath]

	link       link
	linkpicker picker.Model[link]
	linkfixer  picker.Model[relativePath]
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
		m.linkfixer = m.linkfixer.Height(msg.Height - 3)
		return m, nil

	case picker.SelectedMsg[relativePath]:
		switch m.state {
		case filePickerView:
			file, links := readFile(m.config, msg.Selected)

			m.state = linkPickerView
			m.file = file
			m.linkpicker = m.linkpicker.Items(links)

		case linkFixerView:
			m.state = linkPickerView
			updated := fixLink(m.config, m.file, m.link, msg.Selected)
			updateFile(updated)
			file, links := readFile(m.config, updated.path)

			m.file = file
			m.linkpicker = m.linkpicker.Items(links)
		}

	case picker.SelectedMsg[link]:
		m.state = linkFixerView

		parts := strings.Split(msg.Selected.url, "/")
		last := parts[len(parts)-1]

		m.linkfixer = m.linkfixer.Search(last)
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

func resolveLink(config Config, relative string, url string) (linkStatus, relativePath) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return remote, relativePath(url)
	}

	p := url + config.resolveExtension

	if strings.HasPrefix(p, "../") {
		p = filepath.Join(filepath.Dir(relative), p)
	} else {
		p = config.removeAlias(p)
	}

	stat, statErr := os.Stat(p)
	if statErr != nil {
		return unresolved, relativePath(p)
	}

	if stat.IsDir() {
		return unresolved, relativePath(p)
	}

	return resolved, relativePath(p)
}

func fixLink(config Config, file file, link link, path relativePath) file {
	oldLink := fmt.Sprintf("[%s](%s)", link.name, link.url)
	newLink := fmt.Sprintf("[%s](%s)", link.name, strings.TrimSuffix(config.addAlias(string(path)), config.resolveExtension))

	file.contents = strings.Replace(file.contents, oldLink, newLink, 1)
	return file
}

func updateFile(file file) {
	osFile, err := os.Create(string(file.path))
	if err != nil {
		panic(fmt.Errorf("Failed to open file: %v", err))
	}

	_, err = osFile.WriteString(file.contents)
	if err != nil {
		panic(fmt.Errorf("Failed to update file: %v", err))
	}

	err = osFile.Close()
	if err != nil {
		panic(fmt.Errorf("Failed to close file: %v", err))
	}
}

func readFile(config Config, path relativePath) (file, []link) {
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
			status, resolved := resolveLink(config, string(path), url)

			links = append(links,
				link{
					name,
					url,
					resolved,
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

func (m model) linkFixerView() string {
	selected := m.file.path
	header := lg.JoinVertical(
		lg.Top,
		theme.Heading.Render("Fix links for")+theme.Primary.MarginLeft(1).Render(string(selected)),
		lg.NewStyle().MarginLeft(4).Foreground(theme.ColorSecondary).Render(m.link.name),
	)

	return lg.JoinVertical(
		lg.Top,
		header,
		m.linkfixer.View(),
	)
}

func (m model) View() string {
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

func initialModel(config Config, files []relativePath) model {
	return model{
		config:     config,
		state:      filePickerView,
		filepicker: picker.New[relativePath]().Title("File to check").Accent(theme.ColorPrimary).Items(files),
		linkpicker: picker.New[link]().Title("Edit Link").Accent(theme.ColorSecondary),
		linkfixer:  picker.New[relativePath]().Title("Fix link").Accent(theme.ColorWarn).Items(files),
	}
}

func main() {
	config := Config{
		resolveExtension: ".md",
		aliases: []alias{
			{actual: "docs-md", alias: "/docs"},
			{actual: "blog-md", alias: "/blog"},
			{actual: "", alias: "/"},
		},
	}

	files := getMarkdownFiles(config)

	m := initialModel(config, files)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func getMarkdownFiles(config Config) []relativePath {
	var files []relativePath

	root := "./"
	filepath.WalkDir(root,
		func(s string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && strings.HasSuffix(s, config.resolveExtension) {
				files = append(files, relativePath(path.Join(root, s)))
			}

			return nil
		},
	)

	return files
}
