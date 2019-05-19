package gistviewer

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"os"
)

var app *views.Application

type viewer struct {
	input *input
	list  *list
	status *views.SimpleStyledTextBar
	model *listModel
	Selection string

	views.Panel
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
			v.list.view.SetCursor(0, 0)
			v.addRuneToSearch(ev.Rune())
			app.Update()
			return true
		}

		if ev.Key() == tcell.KeyBackspace2 || ev.Key() == tcell.KeyBackspace {
			v.deleteRuneFromSearch()
			app.Update()
			return true
		}

		if (ev.Key() == tcell.KeyCtrlG) {
			v.list.switchMode()
			v.list.model.selectedItem = v.list.model.allVisibleItems[0]
			//v.list.s//

			// TODO Switch between grouped (with count) and ungrouped (with date) for formatted item string
			//groupedItemMap := v.list.model.groupedItemMap
			////for _, e := range groupedItems {
			//	e.formatted += strconv.Itoa(len(groupedItemMap[e.cmdexpr]))
			//}

			app.Update()
			return true
		}
	}
	return v.Panel.HandleEvent(e)
}

func (v *viewer) addRuneToSearch(r rune) {
	v.input.model.appendRune(r)
	v.list.filter(v.input.terms())
}

func (v *viewer) deleteRuneFromSearch() {
	v.input.model.deleteRune()
	v.list.filter(v.input.terms())
}

func NewViewer() *viewer {

	v := &viewer{}

	inputModel := &inputModel{line: ""}
	i := &input{
		view:  views.NewTextArea(),
		model: inputModel,
	}
	i.view.SetModel(inputModel)
	i.view.SetStyle(tcell.StyleDefault.Background(tcell.ColorNavy))

	history := NewHistory()
	listModel := &listModel{
		history: history,
		groupedItemMap: make(map[string][]*item, 0),
	}
	listModel.createItems()

	l := &list{
		view: views.NewCellView(),
	}
	l.SetContent(l.view)
	l.model = listModel
	l.view.SetModel(listModel)
	l.collect()

	v.input = i
	v.list = l
	v.SetOrientation(views.Vertical)
	v.AddWidget(i.view, 0.01)
	v.AddWidget(l, 0.5)
	v.model = listModel

	v.status = views.NewSimpleStyledTextBar()
	v.SetStatus(v.status)

	app = &views.Application{}
	app.SetRootWidget(v)

	if e := app.Run(); e != nil {
		println(e.Error())
		os.Exit(1)
	}

	return v
}
