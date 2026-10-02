package tui

import (
	"github.com/alexbudy/go_spanish_rewrite/internal/store"
	tea "github.com/charmbracelet/bubbletea"
)

/* components needed for word manage screen */

const entriesPerPage = 25

type manageWordsList struct {
	title     string       // the title for managing words
	items     []store.Word // all words
	pageStart int          // 0 base start of words to show
	cursor    int          // between pageStart and pageEnd, inclusive
	pageEnd   int          // 0 base end of words to show (TODO - can do without)
	lang1 	string 		   //'from' language
	lang2 	string  	   //'to' language
}


func newManageWordsList(title string, items []store.Word, lang1 string, lang2 string) manageWordsList {
	return manageWordsList{
		title: title, items: items, pageStart: 0, cursor: 0,
		pageEnd: entriesPerPage, lang1: lang1, lang2: lang2,
	}
}

func (m *Model) buildManageWordsMenu() {
	// m.manageWordsMenu = newManageWordsList("Select a word", m.allWords, m.locale) // TODO
}

// Manage words screen - show all words, allow for reset, removal (TODO)
func (m Model) viewManageWords() string {
	// TODO
	return ""
}

func (m Model) updateManageWords(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO
	return m, nil
}
