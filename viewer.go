package gistviewer

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"log"
	"os"
	"os/exec"
)

var app *views.Application

type viewer struct {
	input *input
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

		if ev.Key() == tcell.KeyEnter {
			app.Quit()
			c := exec.Command("")
			fmt.Printf("\n\ncmd %v\n\n", m.list.model.selectedItem.cmd)
			if err := c.Run(); err != nil {
				log.Fatal(err)
			}
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

		if ev.Key() == tcell.KeyRune {
			m.input.appendRune(ev.Rune())
			m.list.filter(m.input.model.line)
			app.Update()
			return true
		}

		if ev.Key() == tcell.KeyBackspace2|| ev.Key() == tcell.KeyBackspace {
			m.input.deleteRune()
			app.Update()
			return true
		}
	}
	return m.Panel.HandleEvent(e)
}

func NewViewer() *viewer {

	v := &viewer{}

	inputModel := &inputModel{line:""}
	i := &input {
		view: views.NewTextArea(),
		model: inputModel,
	}
	i.view.SetModel(inputModel)
	i.view.SetStyle(tcell.StyleDefault.Background(tcell.ColorNavy))


	listModel := &listModel{history: NewHistory(), endx:60, endy:120}
	l := &list{
		view: views.NewCellView(),
	}
	listModel.loadHistory()

	l.SetContent(l.view)
	l.model = listModel
	l.view.SetModel(listModel)

	v.input = i
	v.list = l
	v.SetOrientation(views.Vertical)
	v.AddWidget(i.view, 0.01)
	v.AddWidget(l, 0.5)
	v.model = listModel

	app = &views.Application{}
	app.SetRootWidget(v)

	if e := app.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
	return v
}
