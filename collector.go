package gistviewer

import "strconv"

// Collector gathers together identical lines
func (l *list) collect(){
	cmdExprs := make(map[string]int, len(l.model.history.allItems))
	for _, i := range l.model.history.allItems {
		cmdExprs[i.cmdexpr] ++
		g := l.model.groupedItemMap[i.cmdexpr]
		g = append(g, i)
		i.formatted += strconv.Itoa(len(i.cmdexpr))
	}
}