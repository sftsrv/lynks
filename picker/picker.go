package picker

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/sftsrv/lynks/theme"
)

type Model struct {
	title     string
	search    string
	items     []string
	filtered  []string
	searching bool
	cursor    int
}

func New() Model {
	return Model{}
}

func (m Model) Items(items []string) Model {
	m.items = items
	return m
}

func (m Model) Title(title string) Model {
	m.title = title
	return m
}

func (_ Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	header := theme.Heading.Render(m.title) + theme.Faded.Render(" / to search")
	if m.searching {
		header = theme.Heading.Render("Search") + " " + m.search + "_"
	}

	return lg.JoinVertical(
		lg.Top,
		header,
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	maxIndex := max(len(m.filtered)-1, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		str := msg.String()
		switch str {
		case "left", "h":
			m.cursor = 0
		case "right", "l":
			m.cursor = maxIndex

		case "up", "k":
			m.cursor = clamp(m.cursor-1, 0, maxIndex)

		case "down", "j":
			m.cursor = clamp(m.cursor+1, 0, maxIndex)

		case "esc":
			m.searching = false

		case "/":
			m.searching = true

		case "backspace":
			if m.search != "" {
				m.search = m.search[0 : len(m.search)-1]
			}

		default:
			if len(str) == 1 {
				m.search += str
			}
		}
	}

	return m, nil
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
