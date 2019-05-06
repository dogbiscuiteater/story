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
	terms []string
	cursor int

	views.CellModel
}

func (m *inputModel) appendRune(r rune){
	m.line += string(r)
	m.cursor = len(m.line) - 1
}

func (m *inputModel) deleteRune(){
	if m.cursor == 0 || len(m.line) == 0 { return }
	m.line = m.line[:len(m.line)-1]
	m.cursor = len(m.line)
}

func (i *input) terms() []string {
	return i.model.terms
}

func (i *input) line() string {
	return i.model.line
}

func (m *inputModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
    //
	style := tcell.StyleDefault.Bold(true)
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
    return m.cursor, 0, true, true
}

func (m *inputModel) MoveCursor(offx, offy int) {
	if m.line == "" { return }
	if m.cursor + offx >= len(m.line) { m.cursor = len(m.line) -1 } else {
		m.cursor = m.cursor + offx
	}
}
