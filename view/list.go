package lowdown

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	lowdown "lowdown/history"
	"sort"
	"strconv"
	"strings"
	"time"
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

//HandleEvent delegates to the Panel view
func (l *list) HandleEvent(ev tcell.Event) bool {
	return l.view.HandleEvent(ev)
}

type listModel struct {
	history         *lowdown.History
	selectedItem    *item
	allItems        []*item
	allVisibleItems []*item

	groupedItemMap  map[string][]*item
	groupedItems    []*item
	mode            mode

	x    int
	y    int
	endx int
	endy int
}

func (l *list) switchMode() {
	l.model.update()
	l.view.HandleEvent(tcell.NewEventKey(tcell.KeyHome, ' ', 0))
}

func (m *listModel) update() {
	if m.mode == grouped {
		m.mode = dateOrder
		m.allVisibleItems = m.allItems
	} else {
		m.mode = grouped
		m.allVisibleItems = m.groupedItems
	}

	m.endy = len(m.allVisibleItems) - 1
	m.sort()
}

func (m *listModel) sort() {
	var sortFunc func(i, j int) bool
	if m.mode == grouped {
		sortFunc = m.sortGrouped()
	} else {
		sortFunc = m.sortInDateOrder()
	}
	sort.Slice(m.allVisibleItems, sortFunc)
	if len(m.allVisibleItems) == 0 || len(m.allVisibleItems[0].highlights) == 0 { return }

	for i, item := range m.allVisibleItems {
		if len(item.highlights) ==  0 {
			m.allVisibleItems = m.allVisibleItems[0:i]
			break;
		}
	}
}

func (m *listModel) sortGrouped() func(i, j int) bool {
	return func(i, j int) bool {
		a := m.allVisibleItems[i]
		b := m.allVisibleItems[j]
		if  a.rank() == b.rank() {
			return len(m.groupedItemMap[a.cmdexpr]) > len(m.groupedItemMap[b.cmdexpr])
		}
		return a.rank() > b.rank()
	}
}

func (m *listModel) sortInDateOrder() func(i, j int) bool {
	return func(i, j int) bool {
		a := m.allVisibleItems[i]
		b := m.allVisibleItems[j]
		if  a.rank() == b.rank() {
			return a.timestamp.After(b.timestamp)
		}
		return a.rank() > b.rank()
	}
}

func (m *listModel) loadHistory() *listModel {
	done := make(chan bool, 1)
	go func(chan bool) {
		h := lowdown.NewHistory()
		done <- true
		m.history = h
	}(done)

	<-done
	return m
}

func (m *listModel) createItems() {
	for i := len(m.history.Lines()) - 1; i >= 0; i-- {
		v := m.history.Lines()[i]
		item := newItem(v, m.history.Fmt)
		m.allItems = append(m.allItems, item)
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

	// Get hold of the visible text in the item
	text := m.allVisibleItems[y].formatted
	if selectedMode == grouped {
		text = m.allVisibleItems[y].grouped
	}

	// The rune is inside the left margin (1st 29 chars?!) so it is colored differently
	if x < leftMargin {
		style = style.Foreground(tcell.ColorRed)
	}

	// The rune may be r.padding
	if x >= len(text) {
		ch = ' '
	} else {
		item := m.allVisibleItems[y]
		ch = rune(text[x])

		// If this character is part of a highlighted term, then use the highlighting style for it.
		if item.highlights != nil {
			for _, h := range item.highlights {
				for i := 0; i < len(h.indexes); i += 2 {
					if x >= h.indexes[i] && x < h.indexes[i+1] {
						style = style.Background(tcell.ColorOrange).Foreground(tcell.ColorBlack).Bold(true)
						break
					}
				}
			}
		}
	}
	return ch, style, nil, 1
}

// collect gathers together identical lines
func (l *list) collect() {
	cmdExprs := make(map[string]bool, len(l.model.allItems))
	for _, i := range l.model.allItems {
		if cmdExprs[i.cmdexpr] {
			l.model.groupedItemMap[i.cmdexpr] = append(l.model.groupedItemMap[i.cmdexpr], i)
		} else {
			l.model.groupedItemMap[i.cmdexpr] = []*item{i}
			l.model.groupedItems = append(l.model.groupedItems, i)
			cmdExprs[i.cmdexpr] = true
		}
	}

	for _, i := range l.model.allItems {
		count := "(" + strconv.Itoa(len(l.model.groupedItemMap[i.cmdexpr])) + ")"
		padding := strings.Repeat(" ", 29-len(count))
		i.grouped = count + padding + " : " + i.cmdexpr
	}
}

// Item is an entry in a shell history. It contains the timestamp, command expression, search terms and highlighted terms
type item struct {

	timestamp time.Time
	fmt       *lowdown.HistoryFormat
	entry     string
	formatted string
	grouped	  string
	cmdexpr   string
	cmd       string
	cmdArgs   string
	words     []string
	highlights []highlights
}

func newItem(entry string, fmt *lowdown.HistoryFormat) *item {
	h := &item{
		entry: entry,
		fmt:   fmt,
	}
	h.split()
	return h
}

func (i *item) split() {
	elements := strings.Split(i.entry, ";")

	// Get the timestamp element
	s := strings.Split(elements[0], ":")
	if len(s) > 1 {
		t, _ := strconv.ParseInt(strings.TrimSpace(s[1]), 10, 64)
		i.timestamp = time.Unix(t, 0)
	}

	// Get the command element
	if len(elements) == 1 {
		i.cmdexpr = ""
	} else {
		i.cmdexpr = elements[1]
		i.cmd = strings.TrimSpace(strings.Split(i.cmdexpr, " ")[0])
		i.cmdArgs = strings.TrimSpace(strings.TrimPrefix(i.cmdexpr, i.cmd))
	}
	i.formatted = i.timestamp.String() + " : " + i.cmdexpr
	i.words = strings.Split(i.cmdArgs, " ")
}

func (i *item) rank() int {
	rank := len(i.highlights) * 10
	for _, h := range i.highlights {
		l := len(h.indexes) / 2
		if l > 9 { l = 9}
		rank += l
	}
	return rank
}