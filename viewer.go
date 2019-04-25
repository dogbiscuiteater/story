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

type viewerModel struct {
}

func (m *viewer) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape {
			app.Quit()
			return true
		}
	}
	return m.Panel.HandleEvent(ev)
}

func hst() []string {
	h := make([]string, 10)
	for i := 0; i < 10; i++ {
		h[i] = "history item " + strconv.Itoa(i)
	}
	return h
}

func NewViewer() *viewer {
	v := &viewer{}

	i := views.NewTextArea()
	i.SetStyle(tcell.StyleDefault.Background(tcell.ColorNavy))

	m := &listModel{history:hst(), endx:60, endy:20}
	l := &list{}
	l.Init()
	l.SetModel(m)
	l.SetStyle(tcell.StyleDefault.Background(tcell.ColorOrange))

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
