package lowdown

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"strings"
)

type input struct {
	view  *views.TextArea
	model *inputModel
}

type inputModel struct {
	line   string
	terms  []string
	cursor int

	views.CellModel
}

func (m *inputModel) appendRune(r rune) {
	if m.cursor == len(m.line){
		m.line += string(r)
	} else {
		m.line = m.line[0:m.cursor] + string(r) + m.line[m.cursor:len(m.line)]
	}
	m.updateTerms()
	m.cursor++
}

func (m *inputModel) deleteRune() {
	if m.cursor == 0 || len(m.line) == 0 {
		return
	}
	m.cursor--
	m.line = m.line[0:m.cursor] + m.line[m.cursor+1:len(m.line)]
	m.updateTerms()
}

func (m *inputModel) updateTerms() {
	// If the line is blank, then clear the model's terms.
	if len(m.line) == 0 {
		m.terms = make([]string, 0)
		return
	}

	// If the last character in the line is a space, then ignore it.
	if m.line[len(m.line)-1] == ' ' { return }

	// Otherwise, re-compute the model terms.
	m.terms = make([]string, 0)
	for _, t := range strings.Split(m.line, " ") {
		 m.terms = append(m.terms, strings.TrimSpace(t))
	}
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
	if m.line == "" || x >= len(m.line) {
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
	if m.line == "" {
		return
	}
	if m.cursor+offx >= len(m.line) {
		m.cursor = len(m.line) - 1
	} else {
		m.cursor = m.cursor + offx
	}
}
