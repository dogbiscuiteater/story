package gistviewer

import (
	"github.com/gdamore/tcell"
	"regexp"
	"sort"
	"strings"
)

type highlights struct {
	spans [][]int
	matches int
}

func (m *list) filter(searchTerms []string) {

	// Create a mapping of Item to text span (
	spans := make(map[*item][][]int, 0)
	h := make(map[*item]highlights, 0)

	m.model.highlights = h
	m.view.HandleEvent(tcell.NewEventKey(tcell.KeyUp, ' ', 0))

	if len(searchTerms) == 0 {
		m.model.history.allVisibleItems = m.model.history.allItems
		return
	}
	
	v := make([]*item, 0)
	matches := make(map[*item]int)
	var item *item

	for _, item = range m.model.history.allItems {
		for _, searchTerm := range searchTerms {
			if strings.Contains(item.formatted, searchTerm) {
				matches[item]++
				re := regexp.MustCompile(searchTerm)

				for _, indices := range re.FindAllStringIndex(item.formatted, 10) {
					spans[item] = append(spans[item], indices)
				}
			}
		}

		if matches[item] == 0 { continue }

		v = append(v, item)
		h[item] = highlights{
			spans:   spans[item],
			matches: matches[item],
		}
	}

	sort.Slice(v, func(i,j int)bool {
		return h[v[i]].matches > h[v[j]].matches
	})

	//sort.Slice(v, func(i,j int) bool { v[i]. }())
	m.model.history.allVisibleItems = v

	if len(m.model.history.allVisibleItems) > 0 {
		m.model.selectedItem = m.model.history.allVisibleItems[0]
	}

	m.model.endy = len(m.model.history.allVisibleItems) - 1
	m.model.y = 0
}

