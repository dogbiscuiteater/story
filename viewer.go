package gistviewer

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"os"
	"strconv"
)

var app *views.Application

type viewer struct {
	input *views.TextArea
	list  *list
	model *listModel

	views.Panel
}

func (m *viewer) HandleEvent(e tcell.Event) bool {
	switch ev := e.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape {
			app.Quit()
			return true
		}

		if ev.Key() == tcell.KeyDown {
			m.list.HandleEvent(e)
			app.Update()
			return true
		}

		if ev.Key() == tcell.KeyUp {
		m.list.HandleEvent(e)
			app.Update()
			return true
		}
	}
	return m.Panel.HandleEvent(e)
}

func hst() []string {
	h := make([]string, 100)
	for i := 0; i < 100; i++ {
		h[i] = "history item " + strconv.Itoa(i)
	}
	return h
}

func NewViewer() *viewer {
	v := &viewer{}

	i := views.NewTextArea()
	i.SetStyle(tcell.StyleDefault.Background(tcell.ColorNavy))

	m := &listModel{history:hst(), endx:60, endy:120}
	l := &list{
		main: views.NewCellView(),
	}
	l.SetContent(l.main)
	l.main.SetModel(m)

	v.input = i
	v.list = l
	v.SetOrientation(views.Vertical)
	v.AddWidget(i, 0.01)
	v.AddWidget(l, 0.5)
	v.model = m

	app = &views.Application{}
	app.SetRootWidget(v)

	if e := app.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
	return v
}
