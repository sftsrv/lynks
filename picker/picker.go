package picker

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/sftsrv/lynks/theme"
)

type Model struct {
	title     string
	search    string
	searching bool
	cursor    int
	count     int
	items     []string
	filtered  []string
}

func New() Model {
	return Model{
		count: 5,
	}
}

func (m Model) Items(items []string) Model {
	debugItems := []string{}

	for i, item := range items {
		debugItems = append(debugItems, fmt.Sprintf("%d %s", i+1, item))
	}

	m.items = debugItems
	m.filtered = debugItems
	return m

	// m.items = items
	// m.filtered = items
	// return m
}

func (m Model) Title(title string) Model {
	m.title = title
	return m
}

func (_ Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	header := theme.
		Heading.
		Render(fmt.Sprintf("%s (%d/%d)", m.title, m.cursor+1, len(m.filtered))) +
		theme.Faded.Render(" / to search")

	if m.searching {
		header = theme.Heading.Render("Search") + " " + m.search + "_"
	}

	cursor, items := m.cursorWindow()
	content := []string{}

	for i, item := range items {
		if i == cursor {
			content = append(content, lg.NewStyle().Foreground(theme.Primary).Render(item))
		} else {
			content = append(content, item)
		}
	}

	return lg.JoinVertical(
		lg.Top,
		header,
		lg.JoinVertical(lg.Top, content...),
	)
}

// Gets the cursor position in a relative window with one item padding if possible.
// Prefers to keep cursor at the top
func (m Model) cursorWindow() (int, []string) {
	itemCount := len(m.items)

	if m.cursor < 2 {
		return m.cursor, m.items[0:min(m.count, itemCount)]
	}

	if m.cursor > itemCount-1 {
		items := m.items[max(0, itemCount-m.count):itemCount]
		lastItem := len(items) - 1
		return lastItem, items
	}

	first := m.cursor - 1
	last := min(m.cursor+m.count, itemCount)
	return 1, m.filtered[first:last]

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
				// filter items based on search here
				m.filtered = m.items
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
