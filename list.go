package gistviewer

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"sort"
)


type mode int

const (
	dateOrder mode = iota
	grouped

)

type list struct {
	view  *views.CellView
	model *listModel
	views.Panel
}

func (m *list) HandleEvent(ev tcell.Event) bool {
	return m.view.HandleEvent(ev)
}

// Row number, followed by start and end locations in the item to identify the highlighted span of runes
type highlightedSpan [3]int

type listModel struct {
	history        *History
	selectedItem   *item
	highlights     map[*item]highlights
	groupedItemMap map[string][]*item
	groupedItems	[]*item
	mode mode

	x    int
	y    int
	endx int
	endy int
}

func (l *list)switchMode() {
	m := l.model.mode
	if m == grouped {
		l.model.mode = dateOrder
	} else {
		l.model.mode = grouped
		groupedItems := l.model.groupedItems
		l.model.history.allVisibleItems = groupedItems
		sort.Slice(groupedItems,
			func(i, j int) bool {
				return len(l.model.groupedItemMap[groupedItems[i].cmdexpr]) > len(l.model.groupedItemMap[groupedItems[j].cmdexpr])
			},
		)
	}
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

	selectedMode := m.mode
	leftMargin := 29
	text := m.history.allVisibleItems[y].formatted

	if selectedMode == grouped {
		text = m.history.allVisibleItems[y].grouped
	}

	if x < leftMargin {
		style = style.Foreground(tcell.ColorRed)
	}


	if x >= len(text) {
		ch = ' '
	} else {
		i := m.history.allVisibleItems[y]

		ch = rune(text[x])
		if m.highlights[i].spans != nil {
			for _, h := range m.highlights[i].spans{
				if x >= h[0] && x < h[1]  {
					style = style.Background(tcell.ColorOrange).Foreground(tcell.ColorBlack).Bold(true)
					break
				}
			}
		}
	}
	return ch, style, nil, 1
}
