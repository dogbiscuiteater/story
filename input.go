package gistviewer

import "github.com/gdamore/tcell/views"

type input struct {
	view  *views.TextArea
	model *inputModel
}

type inputModel struct {
	line string
}
