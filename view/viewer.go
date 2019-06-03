package story

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"os"
	story "story/history"
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

		if ev.Key() == tcell.KeyLeft {
			c, _,_,_ := v.input.model.GetCursor()
			if c == 0 { return true }
			c--
			v.input.model.SetCursor(c, 0)
			return true
		}

		if ev.Key() == tcell.KeyRight {
			c, _,_,_ := v.input.model.GetCursor()
			if c == len(v.input.line()) { return true }
			c++
			v.input.model.SetCursor(c, 0)
			return true
		}

		if ev.Key() == tcell.KeyRune {
			v.list.cellView.SetCursor(0, 0)
			v.addRuneToSearch(ev.Rune())
			app.Update()
			return true
		}

		if ev.Key() == tcell.KeyBackspace2 || ev.Key() == tcell.KeyBackspace {
			v.deleteRuneFromSearch()
			app.Update()
			return true
		}

		if ev.Key() == tcell.KeyCtrlG {
			v.list.switchMode()
			v.list.model.selectedItem = v.list.model.allVisibleItems[0]
			v.status.SetLeft("Order by: " + string(v.list.view()))
			app.Update()
			return true
		}
	}
	return v.Panel.HandleEvent(e)
}

func (v *viewer) keybarText() string {
	s := "%B[CTL-M]%N Change Mode %B[ESC]%N Exit"
	return s
}

func (v *viewer) statusText() string{
	return "Order by: " + string(v.list.view()) + "Showing "
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

	history := story.NewHistory()
	listModel := &listModel{
		history:        history,
		groupedItemMap: make(map[string][]*item, 0),
		mode:           date,
	}
	listModel.createItems()

	l := &list{
		cellView: views.NewCellView(),
	}
	l.model = listModel
	l.cellView.SetModel(listModel)
	l.collect()

	v.input = i
	v.list = l
	v.SetOrientation(views.Vertical)
	v.model = listModel

	title := views.NewSimpleStyledTextBar()
	title.SetCenter("story : history viewer")
	title.SetRight("v.0.1")

	menu := views.NewSimpleStyledText()
	menu.SetMarkup(v.keybarText())

	app = &views.Application{}
	app.SetRootWidget(v)

	status := views.NewSimpleStyledTextBar()
	status.SetLeft(v.statusText())
	v.status = status

	p := views.NewBoxLayout(views.Vertical)
	p.AddWidget(i.view, 0.01)
	p.AddWidget(l.cellView, 0.5)
	v.SetTitle(title)
	v.SetMenu(menu)
	v.SetContent(p)
	v.SetStatus(status)

	if e := app.Run(); e != nil {
		println(e.Error())
		os.Exit(1)
	}

	return v
}
