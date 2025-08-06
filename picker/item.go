package picker

type Item interface {
	Title() string
}

type ItemSource[I Item] struct {
	items []I
}

func (s ItemSource[I]) Len() int {
	return len(s.items)
}

func (s ItemSource[I]) String(i int) string {
	return s.items[i].Title()
}

func getTitles[I Item](items []I) []string {
	strs := []string{}

	for _, i := range items {
		strs = append(strs, i.Title())
	}

	return strs

}
