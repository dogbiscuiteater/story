package gistviewer

import "github.com/gdamore/tcell/views"

type input struct {
	views.TextArea
	model *inputModel
}

type inputModel struct {
	line string
}
