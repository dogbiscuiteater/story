package gistviewer

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"regexp"
	"strings"
)

var commonCommandText = map[string]commonCommand{
	"vim":   VIM,
	"nvim":  NVIM,
	"cd":    CD,
	"ls":    LS,
	"mkdir": MKDIR,
}


type commonCommand int

const (
	VIM commonCommand = iota
	NVIM
	CD
	LS
	MKDIR
)

type ordering int

const (
	DateAsc ordering = iota
	DateDesc
	FrequecyDesc
)

// Create a cheaty cache
var none = newFilterSettings(DateAsc, false)
var nvim = newFilterSettings(DateAsc, true)
var vimAndNvim = newFilterSettings(DateAsc, true, true)
var cdLsMkdir = newFilterSettings(DateAsc, false, false, true, true, true)

type History struct {
	allItems [] 	*Item
	wordsToItems 	map[string][]*Item
	allVisibleItems [] *Item
	allItemsRunes   [][] rune
	lines           [] string
	filters         map[commonCommand]bool
	ordering        ordering
	fmt             *HistoryFormat

	filteredItemCache	filteredItemCache
}

type filteredItemCache struct {
	filteredItems map[filterSettings][]*Item
	filteredItemRunes map[filterSettings][][]rune
}

type HistoryFormat struct {
	timestamp string
	entry string
	cmd string

}

type filterSettings struct {
	filters  [5]bool
	ordering ordering
}

func (h *History) Filters (filters...commonCommand){
	for _, f := range filters{
		h.filters[f] = true
	}
}

func (h *History) Ordering() ordering {
	return h.ordering
}

func (h *History) Order (o ordering){
	h.ordering = o
}

//NewHistformat
func newHistoryFormat(timestamp, entry, expr string) *HistoryFormat {
	f := &HistoryFormat{
		timestamp: timestamp,
		entry: entry,

	}
	return f
}

func NewHistory() *History {
	viper.AutomaticEnv()

	// Get the user home directory and the shell and infer the history filename
	home := viper.GetString("HOME")
	shellPath := strings.Split(viper.GetString("SHELL"),"/")
	shell := shellPath[len(shellPath)-1]
	histfileName := home+ "/." + shell + "_history"

	// Read the history file!
	historyFileContents, _ := ioutil.ReadFile(histfileName)

	h := &History{
		ordering: DateAsc,
		fmt:      newHistoryFormat("DD/MM/YYYY:hh:mm:ss", "", ""),
		lines :   strings.Split(string(historyFileContents), "\n"),
		filters : make(map[commonCommand]bool),
		filteredItemCache: filteredItemCache{
			filteredItems: make(map[filterSettings][]*Item, 0),
			filteredItemRunes:	 make(map[filterSettings][][]rune, 0),
		},
	}

	h.create()
	return h
}

 func (h *History) create() {
	for _, v := range h.lines {
		if !validHistLine(v) {
			continue
		}
		i := newItem(v, h.fmt)
		h.allItems = append(h.allItems, i)
		h.allItemsRunes = append(h.allItemsRunes, i.runes)
	}
	h.allVisibleItems = h.allItems
	h.filter(none)
	h.filter(nvim)
	h.filter(vimAndNvim)
	h.filter(cdLsMkdir)
}

func newFilterSettings (o ordering, commandSwitches ... bool) filterSettings{
	var c [5]bool
	for i := 0; i<len(commandSwitches) && i<5; i++ {

		if commandSwitches[i] {
			c[i] = true
		}
	}
	return filterSettings{c, o}
}

func (h *History) filter(s filterSettings){
	filteredItems := make([]*Item, 0)
	filteredItemRunes := make([][]rune, 0)

	for _, v := range h.allItems{
		show := !s.filters[commonCommandText[v.cmd]]
		if show{
			filteredItems = append(filteredItems, v)
			filteredItemRunes = append(filteredItemRunes, v.runes)
		}
	}
	h.filteredItemCache.filteredItemRunes[s] = filteredItemRunes
	h.filteredItemCache.filteredItems[s] = filteredItems
}

func validHistLine(l string) bool {
	if len(l) <= 1 || strings.IndexRune(l, ':') != 0 {
		return false
	}
	m, err := regexp.MatchString( ": \\d+:\\d+;.+", l)
	return m == false || err == nil
}
