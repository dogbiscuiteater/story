package gistviewer

import (
	"github.com/gdamore/tcell"
	"regexp"
	"strings"
)

type highlights struct {
	spans [][]int
	matches int
}

func (l *list) filter(searchTerms []string) {

	// Create a mapping of Item to text span
	spans := make(map[*item][][]int, 0)
	h := make(map[*item]highlights, 0)

	l.model.highlights = h

	// Nudge the view port to reveal the top line. TODO: find out why the top line gets hidden.
	l.view.HandleEvent(tcell.NewEventKey(tcell.KeyUp, ' ', 0))

	if l.model.mode == grouped {
		l.model.allVisibleItems = l.model.groupedItems
	} else {
		l.model.allVisibleItems = l.model.allItems
	}

	// If the search string is blank then show all of the items.
	if len(searchTerms) == 0 {
		return
	}

	// Prepare to reset the visible items list
	v := make([]*item, 0)
	matches := make(map[*item]int)
	var item *item

	// For each item, check to see if its command expression contains any search term and highlight the first 10
	// occurrences.
	for _, item = range l.model.allVisibleItems {
		for _, searchTerm := range searchTerms {
			if strings.Contains(item.formatted, searchTerm) {
				matches[item]++
				re := regexp.MustCompile(regexp.QuoteMeta(searchTerm))

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

	l.model.allVisibleItems = v

	if len(l.model.allVisibleItems) > 0 {
		l.model.selectedItem = l.model.allVisibleItems[0]
	}

	l.model.endy = len(l.model.allVisibleItems) - 1
	l.model.y = 0
}

