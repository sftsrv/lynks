package files

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/sftsrv/lynks/config"
	"github.com/sftsrv/lynks/theme"
)

type linkStatus int

const (
	resolved linkStatus = iota
	unresolved
	remote
)

type RelativePath string

func (s RelativePath) Title() string {
	return string(s)
}

type Link struct {
	Name string
	Url  string

	Resolved RelativePath
	Status   linkStatus
}

type File struct {
	Path               RelativePath
	Contents           string
	HasLinks           bool
	HasUnresolvedLinks bool
}

const resolveExtension = ".md"

func (l Link) color() lg.Color {
	switch l.Status {
	case remote:
		return theme.ColorSecondary
	case resolved:
		return theme.ColorSecondary
	case unresolved:
		return theme.ColorWarn
	}

	return theme.ColorError
}

func (l Link) Title() string {
	return lg.NewStyle().Foreground(l.color()).Render(lg.NewStyle().Bold(true).Render(l.Name) + " " + l.Url + "->" + string(l.Resolved))
}

func (l Link) FileName() string {
	parts := strings.Split(l.Url, "/")
	return parts[len(parts)-1]
}

func (l Link) IsUnresolved() bool {
	return l.Status == unresolved
}

func ResolveLink(config config.Config, relative string, url string) (linkStatus, RelativePath) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return remote, RelativePath(url)
	}

	p := url + resolveExtension

	if strings.HasPrefix(p, "../") {
		p = filepath.Join(filepath.Dir(relative), p)
	} else {
		p = config.RemoveAlias(p)
	}

	stat, statErr := os.Stat(p)
	if statErr != nil {
		return unresolved, RelativePath(p)
	}

	if stat.IsDir() {
		return unresolved, RelativePath(p)
	}

	return resolved, RelativePath(p)
}

func FixLink(config config.Config, file File, link Link, path RelativePath) File {
	oldLink := fmt.Sprintf("[%s](%s)", link.Name, link.Url)
	newLink := fmt.Sprintf("[%s](%s)", link.Name, strings.TrimSuffix(config.AddAlias(string(path)), resolveExtension))

	file.Contents = strings.Replace(file.Contents, oldLink, newLink, 1)
	return file
}

func GetMarkdownFiles(config config.Config) []RelativePath {
	var files []RelativePath

	root := config.Root
	filepath.WalkDir(root,
		func(s string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && strings.HasSuffix(s, resolveExtension) {
				files = append(files, RelativePath(s))
			}

			return nil
		},
	)

	return files
}

func UpdateFile(file File) {
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

func ReadFile(config config.Config, path RelativePath) (File, []Link) {
	buf, err := os.ReadFile(string(path))
	if err != nil {
		panic(err)
	}

	contents := string(buf)
	linkRe := regexp.MustCompile(`(\s|^)\[.+?\]\(.+?\)`)
	nameRe := regexp.MustCompile(`\[.+?\]`)
	urlRe := regexp.MustCompile(`\(.+?\)`)

	matches := linkRe.FindAllString(contents, -1)
	links := []Link{}

	hasLinks := len(links) > 0
	hasUnresolvedLinks := false

	for _, match := range matches {
		namePart := nameRe.FindString(match)
		urlPart := urlRe.FindString(match)

		if namePart != "" && urlPart != "" {
			name := namePart[1 : len(namePart)-1]
			url := urlPart[1 : len(urlPart)-1]
			status, resolved := ResolveLink(config, string(path), url)

			link := Link{Name: name, Url: url, Resolved: resolved, Status: status}

			links = append(links, link)
			if link.IsUnresolved() {
				hasUnresolvedLinks = true
			}
		}

	}

	return File{Path: path, Contents: contents, HasLinks: hasLinks, HasUnresolvedLinks: hasUnresolvedLinks}, links
}
