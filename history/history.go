package history

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"regexp"
	"strings"
)

type History struct {
	lines []string
	Fmt   *HistoryFormat
}

// HistoryFormat is the currently selected format of a history Item
type HistoryFormat struct {
	timestamp string
	entry     string
	cmd       string
}

//
func newHistoryFormat(timestamp, entry string) *HistoryFormat {
	f := &HistoryFormat{
		timestamp: timestamp,
		entry:     entry,

	}
	return f
}

func (h *History) Lines() []string{
	return h.lines
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
		Fmt: newHistoryFormat("DD/MM/YYYY:hh:mm:ss", "" ),
	}

	lines := strings.Split(string(historyFileContents), "\n")
	nonEmptyLines := make([]string, 0)
	for _, l := range lines {
		if len(strings.TrimSpace(l)) > 0 {
			nonEmptyLines = append(nonEmptyLines, l)
		}
	}
	h.lines = nonEmptyLines
	return h
}

func validHistoryLine(l string) bool {
	if len(l) <= 1 || strings.IndexRune(l, ':') != 0 {
		return false
	}
	m, err := regexp.MatchString(": \\d+:\\d+;.+", l)
	return m == false || err == nil
}
