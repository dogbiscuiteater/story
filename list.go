package gistviewer

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"strconv"
)

type list struct {
	view  *views.CellView
	model *listModel

	views.Panel
}

func (m *list) HandleEvent(ev tcell.Event) bool {
	return m.view.HandleEvent(ev)
}

type listModel struct {
	history []string

	x    int
	y    int
	endx int
	endy int
}

func (m *listModel) GetBounds() (int, int) {
	return m.endx, m.endy
}

func (m *listModel) MoveCursor(offx, offy int) {
	fmt.Sprintln("moving " + strconv.Itoa(offy))
	if m.y+offy >= len(m.history) {
		m.y = len(m.history) - 1
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
	var ch rune

	if y < 0 {
		y = 0
	}

	if y >= len(m.history) {
		y = len(m.history) - 1
	}

	if y == m.y {
		style = style.Background(tcell.ColorLightGreen)
	}

	if x < 29 {
		style = style.Foreground(tcell.ColorRed)
	}

	if x >= len(m.history[y]) {
		ch = ' '
	} else {
		ch = rune(m.history[y][x])
	}
	return ch, style, nil, 1
}
