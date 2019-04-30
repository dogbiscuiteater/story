package gistviewer

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type input struct {
	view *views.TextArea
	model *inputModel
}

type inputModel struct {
	line string
	cursor int

	views.CellModel
}

func (i *input) appendRune(r rune){
	m := i.model
	m.line += string(r)
	m.cursor = len(m.line) - 1
}

func (i *input) deleteRune(){
	m := i.model
	if m.cursor == 0 || len(m.line) == 0 { return }
	m.line = m.line[:len(m.line)-1]
}

func (m *inputModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
    //
	style := tcell.StyleDefault.Background(tcell.ColorOrange)
    var r rune
    if m.line =="" || x >= len(m.line)  {
    	r = ' '
    } else {
    	r = rune(m.line[x])
	}

	return r, style, nil, 1
}

func (m *inputModel) GetBounds() (int, int) {
    return 180, 0
}

func (m *inputModel) SetCursor(x, y int) {
	 m.cursor = x
}

func (m *inputModel) GetCursor() (int, int, bool, bool) {
    return m.cursor, 0, true, false
}

func (m *inputModel) MoveCursor(offx, offy int) {
	if m.line == "" { return }
	if m.cursor + offx >= len(m.line) { m.cursor = len(m.line) -1 } else {
		m.cursor = m.cursor + offx
	}
}
