package picker

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
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
	m.items = items
	m.filtered = items
	return m
}

// The height of the picker is header + count == 1 + count
func (m Model) GetHeight() int {
	return 1 + m.count
}

// The count depends on how much space we have
func (m Model) Height(height int) Model {
	m.count = height - 1
	return m.applyFilter()
}

func (m Model) Title(title string) Model {
	m.title = title
	return m
}

func (_ Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	count := fmt.Sprintf("(%d/%d)", m.cursor+1, len(m.filtered))

	fallback := " / to search"
	if m.search != "" {
		fallback = m.search
	}

	header := theme.
		Heading.
		Render(m.title+" "+count) +
		theme.Faded.Render(fallback)

	if m.searching {
		header = theme.Heading.Render("Search "+count) + " " + m.search + "_"
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
	itemCount := len(m.filtered)

	if m.cursor < 2 {
		return m.cursor, m.filtered[0:min(m.count+1, itemCount)]
	}

	if m.cursor > itemCount-1 {
		items := m.filtered[max(0, itemCount-m.count):itemCount]
		lastItem := len(items) - 1
		return lastItem, items
	}

	first := m.cursor - 1
	last := min(m.cursor+m.count, itemCount)
	return 1, m.filtered[first:last]

}

func (m Model) applyFilter() Model {
	// Must reset the cursor since we're modifying the underlying list
	m.cursor = 0

	if m.search == "" {
		m.filtered = m.items
		return m
	}

	matches := fuzzy.Find(m.search, m.items)

	m.filtered = []string{}
	for _, match := range matches {
		m.filtered = append(m.filtered, m.items[match.Index])
	}

	return m
}

func (m Model) cursorUp() Model {
	maxIndex := max(len(m.filtered)-1, 0)
	m.cursor = clamp(m.cursor-1, 0, maxIndex)

	return m
}

func (m Model) cursorDown() Model {
	maxIndex := max(len(m.filtered)-1, 0)
	m.cursor = clamp(m.cursor+1, 0, maxIndex)

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	maxIndex := max(len(m.filtered)-1, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		str := msg.String()
		if m.searching {
			switch str {

			case "up":
				m = m.cursorUp()

			case "down":
				m = m.cursorDown()

			case "esc", "enter":
				m.searching = false

			case "backspace":
				if m.search != "" {
					m.search = m.search[0 : len(m.search)-1]
					m = m.applyFilter()
				}

			default:
				if len(str) == 1 {
					m.search += str
					m = m.applyFilter()
				}
			}
		} else {
			switch str {
			case "left", "h":
				m.cursor = 0
			case "right", "l":
				m.cursor = maxIndex

			case "up", "k":
				m = m.cursorUp()

			case "down", "j":
				m = m.cursorDown()

			case "/":
				m.searching = true
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
