package picker

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Model struct {
	title     string
	search    string
	items     []string
	filtered  []string
	searching bool
	cursor    int
}

func New(items []string) Model {
	return Model{
		items: items,
	}
}

func (p Model) Title(title string) Model {
	p.title = title

	return p
}

func (_ Model) Init() tea.Cmd {
	return nil
}

func (p Model) View() string {
	header := p.title
	if p.searching {
		header = "Search + " + p.search
	}

	return lg.JoinVertical(
		lg.Top,
		header,
	)
}

func (p Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	m := max(len(p.filtered)-1, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		str := msg.String()
		switch str {
		case "left", "h":
			p.cursor = 0
		case "right", "l":
			p.cursor = m

		case "up", "k":
			p.cursor = clamp(p.cursor-1, 0, m)

		case "down", "j":
			p.cursor = clamp(p.cursor+1, 0, m)

		case "esc":
			p.searching = false

		case "/":
			p.searching = true

		case "backspace":
			if p.search != "" {
				p.search = p.search[0 : len(p.search)-1]
			}

		default:
			if len(str) == 1 {
				p.search += str
			}
		}
	}

	return p, nil
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
