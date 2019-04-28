package gistviewer

import (
	"strconv"
	"strings"
	"time"
)

type Item struct {
	timestamp time.Time
	fmt 	  *HistoryFormat
	entry     string
	formatted  string
	runes     []rune
	cmdexpr   string
	cmd       string
	cmdArgs	  string
	words		[]string
}


func newItem(entry string, fmt *HistoryFormat) *Item {
	h := &Item{
		entry: entry,
		fmt: fmt,
	}
	h.split()
	return h
}

func (i *Item) split() {

	elements := strings.Split(i.entry, ";")

	// Get the timestamp element
	s := strings.Split(elements[0], ":")
	if len(s) > 1	{
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
	i.runes = []rune(i.formatted)
	i.words = strings.Split(i.cmdArgs, " ")
}
