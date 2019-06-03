package story

import (
	"github.com/gdamore/tcell"
	"regexp"
	"strings"
)

type highlights struct {
	term string
	indexes []int
}

func (l *list) filter(searchTerms []string) {

	if l.model.mode == frequency {
		l.model.allVisibleItems = l.model.groupedItems
	} else {
		l.model.allVisibleItems = l.model.allItems
	}

	// If the search string is blank then show all of the items.
	if len(searchTerms) == 0 {
		for _, item := range l.model.allItems {
			item.highlights = make([]highlights, 0)
		}
		return
	}

	// Prepare to reset the visible items list
	v := make([]*item, 0)
	matches := make(map[*item]int)
	var item *item

	// For each item, check to see if its command expression contains any search term and highlight the first 10
	// occurrences.
	for _, item = range l.model.allVisibleItems {
		item.highlights = make([]highlights, 0)
		for _, searchTerm := range searchTerms {
			searchTerm = strings.TrimSpace(searchTerm)
			if strings.Contains(item.formatted, searchTerm) {
				matches[item]++
				re := regexp.MustCompile(regexp.QuoteMeta(searchTerm))

				for _, indices := range re.FindAllStringIndex(item.formatted, 10) {
					item.highlights = append(item.highlights, highlights{ searchTerm, indices })
				}
			}
		}

		if matches[item] == 0 { continue }
		v = append(v, item)
	}

	// Reduce the visible items to just the highlighted items if there are any.
	l.model.allVisibleItems = v
	l.model.endy = len(l.model.allVisibleItems) - 1
	l.model.sort()
	l.cellView.HandleEvent(tcell.NewEventKey(tcell.KeyHome, ' ', 0))
}

