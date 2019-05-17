package gistviewer

import (
	"strconv"
	"strings"
)

// Collector gathers together identical lines
func (l *list) collect(){
	cmdExprs := make(map[string]bool, len(l.model.allItems))
	for _, i := range l.model.allItems {
		if cmdExprs[i.cmdexpr] {
			l.model.groupedItemMap[i.cmdexpr] = append(l.model.groupedItemMap[i.cmdexpr], i)
		} else {
			l.model.groupedItemMap[i.cmdexpr] = []*item{i}
			l.model.groupedItems = append(l.model.groupedItems, i)
			cmdExprs[i.cmdexpr] = true
		}
	}

	for _, i := range l.model.allItems {
		count := "(" + strconv.Itoa(len(l.model.groupedItemMap[i.cmdexpr])) + ")"
		padding := strings.Repeat(" ", 29-len(count))
		i.grouped =  count + padding + " : " + i.cmdexpr
	}
}
