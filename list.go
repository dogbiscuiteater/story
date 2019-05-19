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

func (l *list) HandleEvent(ev tcell.Event) bool {
	return l.view.HandleEvent(ev)
}

type listModel struct {

	history        *History
	selectedItem   *item
	allItems	   []*item
	allVisibleItems[]*item
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
		l.model.allVisibleItems = l.model.allItems
	} else {
		l.model.mode = grouped
		l.model.allVisibleItems = l.model.groupedItems
	}
	l.model.endy = len(l.model.allVisibleItems)-1
	l.model.sort()
	l.view.HandleEvent(tcell.NewEventKey(tcell.KeyHome, ' ', 0))
}

func (m *listModel) sort() {
	var sortFunc func(i, j int) bool
	if m.mode == grouped {
		sortFunc = m.sortGrouped()
	} else {
		sortFunc = m.sortInDateOrder()
	}
	sort.Slice(m.allVisibleItems, sortFunc)
}

func (m *listModel) sortGrouped() func(i, j int) bool {
	return func(i, j int) bool {
		a := m.groupedItemMap[m.groupedItems[i].cmdexpr]
		b := m.groupedItemMap[m.groupedItems[j].cmdexpr]
		if m.highlights[a[0]].matches == m.highlights[b[0]].matches {
			return len(m.groupedItemMap[m.groupedItems[i].cmdexpr]) > len(m.groupedItemMap[m.groupedItems[j].cmdexpr])
		}
		return m.highlights[a[0]].matches > m.highlights[b[0]].matches
	}
}

func (m *listModel) sortInDateOrder() func (i, j int) bool {
	return func(i, j int) bool {
		a := m.allItems[i]
		b := m.allItems[j]
		if m.highlights[a].matches == m.highlights[b].matches {
			return m.allItems[i].timestamp.After(m.allItems[j].timestamp)
		}
		return m.highlights[a].matches > m.highlights[b].matches
	}
}

func (m *listModel) sortByHighglights() func (i, j int) bool {
	return func(i, j int) bool {
		return m.highlights[m.allVisibleItems[i]].matches > m.highlights[m.allVisibleItems[j]].matches
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

func (m *listModel) createItems() {
	for i := len(m.history.lines) - 1; i >= 0; i-- {
		v := m.history.lines[i]
		if !validHistLine(v) {
			continue
		}
		i := newItem(v, m.history.fmt)
		m.allItems = append(m.allItems, i)
	}
	m.allVisibleItems = m.allItems
	m.endx = 60
	m.endy = len(m.allItems)
}

func (m *listModel) GetBounds() (int, int) {
	return m.endx, len(m.allVisibleItems)
}

func (m *listModel) MoveCursor(offx, offy int) {
	if m.y+offy >= len(m.allVisibleItems) {
		m.y = len(m.allVisibleItems) - 1
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

	if y >= len(m.allVisibleItems) {
		return ' ', style, nil, 1
	}

	var ch rune

	if y < 0 {
		y = 0
	}

	if y == m.y {
		style = style.Background(tcell.ColorOldLace).Foreground(tcell.ColorBlack).Bold(true)
		m.selectedItem = m.allVisibleItems[y]
	}

	selectedMode := m.mode
	leftMargin := 29
	text := m.allVisibleItems[y].formatted

	if selectedMode == grouped {
		text = m.allVisibleItems[y].grouped
	}


	if x < leftMargin {
		style = style.Foreground(tcell.ColorRed)
	}


	if x >= len(text) {
		ch = ' '
	} else {
		i := m.allVisibleItems[y]

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
