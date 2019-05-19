package gistviewer

import (
	"strconv"
	"strings"
	"time"
)

// Item is an entry in a shell history. It contains the timestamp, command expression, search terms and highlighted terms
type item struct {

	timestamp time.Time
	fmt       *HistoryFormat
	entry     string
	formatted string
	grouped	  string
	cmdexpr   string
	cmd       string
	cmdArgs   string
	words     []string
	highlights []highlights
}

func newItem(entry string, fmt *HistoryFormat) *item {
	h := &item{
		entry: entry,
		fmt:   fmt,
	}
	h.split()
	return h
}

func (i *item) split() {
	elements := strings.Split(i.entry, ";")

	// Get the timestamp element
	s := strings.Split(elements[0], ":")
	if len(s) > 1 {
		t, _ := strconv.ParseInt(strings.TrimSpace(s[1]), 10, 64)
		i.timestamp = time.Unix(t, 0)
	}

	// Get the command element
	if len(elements) == 1 {
		i.cmdexpr = ""
	} else {
		i.cmdexpr = elements[1]
		i.cmd = strings.TrimSpace(strings.Split(i.cmdexpr, " ")[0])
		i.cmdArgs = strings.TrimSpace(strings.TrimPrefix(i.cmdexpr, i.cmd))
	}
	i.formatted = i.timestamp.String() + " : " + i.cmdexpr
	i.words = strings.Split(i.cmdArgs, " ")
}
