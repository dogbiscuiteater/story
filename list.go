package gistviewer

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"regexp"
	"strings"
)

type list struct {
	view  *views.CellView
	model *listModel
	views.Panel
}

func (m *list) HandleEvent(ev tcell.Event) bool {
	return m.view.HandleEvent(ev)
}

func (m *list) filter(searchTerms []string) {
	h := make(map[*Item][][]int, 0)
	m.model.highlights = h
	m.view.HandleEvent(tcell.NewEventKey(tcell.KeyUp, ' ', 0))
	if len(searchTerms) ==0 { return }

	v := make([]*Item, 0)
	for _, searchTerm := range searchTerms {
		if strings.TrimSpace(searchTerm) == "" {
			m.model.history.allVisibleItems = m.model.history.allItems
			return
		}


		for _, item := range m.model.history.allItems {
			if strings.Contains(item.formatted, searchTerm) {
				re := regexp.MustCompile(searchTerm)
				v = append(v, item)
				for _, indices := range re.FindAllStringIndex(item.formatted, 10){
					h[item]=append(h[item],indices)
				}
			}
		}
	}

	//sort.Slice(v, func(i,j int) bool { v[i]. }())
	m.model.history.allVisibleItems = v

	if len(m.model.history.allVisibleItems) > 0 {
		m.model.selectedItem = m.model.history.allVisibleItems[0]
	}

	m.model.endy = len(m.model.history.allVisibleItems) - 1
	m.model.y = 0
}

// Row number, followed by start and end locations in the item to identify the highlighted span of runes
type highlightedSpan [3]int

type listModel struct {
	history      *History
	selectedItem *Item
	highlights map[*Item][][]int

	x    int
	y    int
	endx int
	endy int
}

func (m *listModel) loadHistory() *listModel {
	done := make(chan bool, 1)
	go func(chan bool) {
		h := NewHistory()
		done <- true
		m.history = h
	}(done)

	<-done
	return m
}

func (m *listModel) GetBounds() (int, int) {
	return m.endx, len(m.history.allVisibleItems)
}

func (m *listModel) MoveCursor(offx, offy int) {
	if m.y+offy >= len(m.history.allVisibleItems) {
		m.y = len(m.history.allVisibleItems) - 1
	} else if m.y+offy < 0 {
		m.y = 0
	} else {
		m.y += offy
	}
}

func (m *listModel) GetCursor() (int, int, bool, bool) {
	return m.x, m.y, true, false
}

func (m *listModel) SetCursor(x int, y int) {
	m.y = y
}

func (m *listModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	style := tcell.StyleDefault

	if y >= len(m.history.allVisibleItems) {
		return ' ', style, nil, 1
	}

	var ch rune

	if y < 0 {
		y = 0
	}

	if y == m.y {
		style = style.Background(tcell.ColorOldLace).Foreground(tcell.ColorBlack).Bold(true)
		m.selectedItem = m.history.allVisibleItems[y]
	}

	if x < 29 {
		style = style.Foreground(tcell.ColorRed)
	}

	if x >= len(m.history.allVisibleItems[y].formatted) {
		ch = ' '
	} else {
		i := m.history.allVisibleItems[y]

		ch = rune(i.formatted[x])
		if m.highlights[i] != nil {
			for _, h := range m.highlights[i]{
				if x >= h[0] && x < h[1]  {
					style = style.Background(tcell.ColorOrange).Foreground(tcell.ColorBlack).Bold(true)
					break
				}
			}
		}
	}
	return ch, style, nil, 1
}
