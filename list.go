package gistviewer

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"strconv"
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

func (m *list) filter(searchTerm string) {
	v := make([]*Item,0)

	for _, i := range m.model.history.allItems{
		if strings.Contains(i.formatted, searchTerm){
			v = append(v, i)
		}
	}
	m.model.history.allVisibleItems = v
}

type listModel struct {
	history *History
	items []*Item
	filteredHistory map[bool]string
	selectedItem *Item

	x    int
	y    int
	endx int
	endy int
}

func (m *listModel) filterHistory(searchTerm string){

}

func (m *listModel) loadHistory() *listModel{
	done := make(chan bool, 1)
	go func(chan bool) {
		h := NewHistory()
		done <-true
		m.history = h
	}(done)

	<- done
	return m
}

func (m *listModel) GetBounds() (int, int) {
	return m.endx, m.endy
}

func (m *listModel) MoveCursor(offx, offy int) {
	fmt.Sprintln("moving " + strconv.Itoa(offy))
	if m.y+offy >= len(m.history.allVisibleItems) {
		m.y = len(m.history.allVisibleItems) - 1
	} else if m.y+offy < 0{
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
		style = style.Background(tcell.ColorLightGreen)
		m.selectedItem = m.history.allVisibleItems[y]
	}

	if x < 29 {
		style = style.Foreground(tcell.ColorRed)
	}

	if  x >= len(m.history.allVisibleItems[y].formatted) {
		ch = ' '
	} else {
		ch = rune(m.history.allVisibleItems[y].formatted[x])
	}
	return ch, style, nil, 1
}
