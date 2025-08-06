package picker

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
	"github.com/sftsrv/lynks/theme"
)

type Model[I Item] struct {
	accent    lg.Color
	title     string
	search    string
	searching bool
	cursor    int
	count     int
	items     []I
	filtered  []I
}

func New[I Item]() Model[I] {
	return Model[I]{
		accent: theme.ColorPrimary,
		count:  5,
	}
}

func (m Model[I]) Items(items []I) Model[I] {
	m.items = items
	m.filtered = items
	return m
}

// The height of the picker is header + count == 1 + count
func (m Model[I]) GetHeight() int {
	return 1 + m.count
}

// The count depends on how much space we have
func (m Model[I]) Height(height int) Model[I] {
	m.count = height - 1
	return m.applyFilter()
}

func (m Model[I]) Accent(accent lg.Color) Model[I] {
	m.accent = accent
	return m
}

func (m Model[I]) Title(title string) Model[I] {
	m.title = title
	return m
}

func (_ Model[I]) Init() tea.Cmd {
	return nil
}

func indicator(accent lg.Color, selected bool, title string) string {
	if !selected {
		return lg.NewStyle().PaddingRight(2).Render("") + title
	}

	return lg.NewStyle().PaddingRight(1).Foreground(accent).Render("→") +
		lg.NewStyle().Foreground(accent).Bold(true).Render(title)
}

func (m Model[I]) View() string {
	count := fmt.Sprintf("(%d/%d)", m.cursor+1, len(m.filtered))

	fallback := "/ to search"
	if m.search != "" {
		fallback = m.search
	}

	header := theme.
		Heading.
		Background(m.accent).
		Render(m.title+" "+count) +
		theme.Faded.MarginLeft(1).Render(fallback)

	if m.searching {
		header = theme.Heading.Background(m.accent).Render("Search "+count) + " " + m.search + "_"
	}

	cursor, items := m.cursorWindow()
	content := []string{}

	for i, item := range items {
		content = append(content, indicator(m.accent, i == cursor, item))
	}

	if len(items) < m.count {
		content = append(content, theme.Faded.Render("no more items"))
	}

	return lg.JoinVertical(
		lg.Top,
		header,
		lg.NewStyle().BorderLeft(true).BorderForeground(m.accent).BorderStyle(lg.NormalBorder()).Render(lg.JoinVertical(lg.Top, content...)),
	)
}

// Gets the cursor position in a relative window with one item padding if possible.
// Prefers to keep cursor at the top
func (m Model[I]) cursorWindow() (int, []string) {
	itemCount := len(m.filtered)

	if m.cursor < 2 {
		items := m.filtered[0:min(m.count, itemCount)]
		return m.cursor, getTitles(items)
	}

	if m.cursor > itemCount-1 {
		items := m.filtered[max(0, itemCount-m.count-1):itemCount]
		lastItem := len(items) - 1
		return lastItem, getTitles(items)
	}

	first := m.cursor - 1
	last := min(m.cursor+m.count-1, itemCount)
	items := m.filtered[first:last]

	return 1, getTitles(items)

}

func (m Model[I]) applyFilter() Model[I] {
	// Must reset the cursor since we're modifying the underlying list
	m.cursor = 0

	if m.search == "" {
		m.filtered = m.items
		return m
	}

	itemSource := ItemSource[I]{m.items}

	matches := fuzzy.FindFrom(m.search, itemSource)

	m.filtered = []I{}
	for _, match := range matches {
		m.filtered = append(m.filtered, m.items[match.Index])
	}

	return m
}

func (m Model[I]) cursorUp() Model[I] {
	maxIndex := max(len(m.filtered)-1, 0)
	m.cursor = clamp(m.cursor-1, 0, maxIndex)

	return m
}

func (m Model[I]) cursorDown() Model[I] {
	maxIndex := max(len(m.filtered)-1, 0)
	m.cursor = clamp(m.cursor+1, 0, maxIndex)

	return m
}

func (m Model[I]) clearSearch() Model[I] {
	m.searching = false
	m.search = ""
	return m
}

type SelectedMsg[I Item] struct {
	Selected I
}

func (m Model[I]) selectedMsg() tea.Cmd {
	return func() tea.Msg {
		return SelectedMsg[I]{
			m.filtered[m.cursor],
		}
	}
}

func (m Model[I]) Update(msg tea.Msg) (Model[I], tea.Cmd) {
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

			case "esc":
				m = m.clearSearch().applyFilter()

			case "enter":
				return m, m.selectedMsg()

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

			case " ", "enter":
				return m, m.selectedMsg()

			case "esc":
				m = m.clearSearch().applyFilter()

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
