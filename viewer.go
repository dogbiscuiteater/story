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
	views.BoxLayout
	input *input
	list  *list
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
	return m.BoxLayout.HandleEvent(ev)
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

	i := &input{}
	i.view = views.NewTextArea()
	i.model = &inputModel{}
	i.view.SetStyle(tcell.StyleDefault.Background(tcell.ColorNavy))

	l := &list{}
	l.view = views.NewCellView()
	l.view.SetModel(&listModel{
		history: hst(),
	})
	l.view.SetStyle(tcell.StyleDefault.Background(tcell.ColorOrange))

	v.input = i
	v.list = l
	v.SetOrientation(views.Vertical)
	v.AddWidget(i.view, 0.01)
	v.AddWidget(l.view, 0.5)

	app = &views.Application{}
	app.SetRootWidget(v)

	if e := app.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}

	return v
}
