package gistviewer


import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"os"
)


var app *views.Application

type viewer struct {
	input *input
	list  *list
	model *listModel
	views.Panel

	Selection string
}

func (v *viewer) HandleEvent(e tcell.Event) bool {
	switch ev := e.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape {
			app.Quit()
			return true
		}

		if ev.Key() == tcell.KeyEnter {
			v.Selection = v.list.model.selectedItem.cmdexpr
			app.Quit()
			return true
		}

		if ev.Key() == tcell.KeyDown {
			v.list.HandleEvent(e)
			app.Update()
			return true
		}

		if ev.Key() == tcell.KeyUp {
			v.list.HandleEvent(e)
			app.Update()
			return true
		}

		if ev.Key() == tcell.KeyRune {
			v.addRuneToSearch(ev.Rune())
			app.Update()
			return true
		}

		if ev.Key() == tcell.KeyBackspace2|| ev.Key() == tcell.KeyBackspace {
			v.deleteRuneFromSearch()
			app.Update()
			return true
		}
	}
	return v.Panel.HandleEvent(e)
}

func (v *viewer) addRuneToSearch(r rune){
	v.input.model.appendRune(r)
	v.list.filter(v.input.line())
}

func (v *viewer) deleteRuneFromSearch(){
	v.input.model.deleteRune()
	v.list.filter(v.input.line())
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

	history := NewHistory()
	listModel := &listModel{history: history, endx:60, endy:len(history.allVisibleItems)}
	l := &list{
		view: views.NewCellView(),
	}

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
