package gistviewer

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"regexp"
	"strings"
)

type ordering int

const (
	DateDesc ordering = iota
	DateAsc
	FrequecyDesc
)

type History struct {
	allItems        []*Item
	wordsToItems    map[string][]*Item
	allVisibleItems []*Item
	lines           []string
	ordering        ordering
	fmt             *HistoryFormat

	aliases aliases
}

type aliases map[string]string

type filteredItemCache struct {
	filteredItems     map[filterSettings][]*Item
	filteredItemRunes map[filterSettings][][]rune
}

type HistoryFormat struct {
	timestamp string
	entry     string
	cmd       string
}

type filterSettings struct {
	filters  [5]bool
	ordering ordering
}

func (h *History) Ordering() ordering {
	return h.ordering
}

func (h *History) Order(o ordering) {
	h.ordering = o
}

//NewHistformat
func newHistoryFormat(timestamp, entry, expr string) *HistoryFormat {
	f := &HistoryFormat{
		timestamp: timestamp,
		entry:     entry,
	}
	return f
}

func NewHistory() *History {
	viper.AutomaticEnv()

	// Get the user home directory and the shell and infer the history filename
	home := viper.GetString("HOME")
	shellPath := strings.Split(viper.GetString("SHELL"), "/")
	shell := shellPath[len(shellPath)-1]
	histfileName := home + "/." + shell + "_history"

	// Read the history file!
	historyFileContents, _ := ioutil.ReadFile(histfileName)

	h := &History{
		ordering: DateDesc,
		fmt:      newHistoryFormat("DD/MM/YYYY:hh:mm:ss", "", ""),
	}

	lines := strings.Split(string(historyFileContents), "\n")
	nonEmptyLines := make([]string, 0)
	for _, l := range lines {
		if len(strings.TrimSpace(l)) > 0 {
			nonEmptyLines = append(nonEmptyLines, l)
		}
	}
	h.lines = nonEmptyLines
	h.createItems()
	return h
}

func (h *History) createItems() {
	for i := len(h.lines) - 1; i >= 0; i-- {
		v := h.lines[i]
		if !validHistLine(v) {
			continue
		}
		i := newItem(v, h.fmt)
		h.allItems = append(h.allItems, i)
	}
	h.allVisibleItems = h.allItems
}

func validHistLine(l string) bool {
	if len(l) <= 1 || strings.IndexRune(l, ':') != 0 {
		return false
	}
	m, err := regexp.MatchString(": \\d+:\\d+;.+", l)
	return m == false || err == nil
}
