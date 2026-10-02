package tui

import (
	"strconv"
	"strings"

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
	pageEnd  int          // number of words to show per page
	lang1 	string 		   //'from' language
	lang2 	string  	   //'to' language
}


func newManageWordsList(title string, items []store.Word, lang1 string, lang2 string) manageWordsList {
	return manageWordsList{
		title: title, items: items, pageStart: 0, cursor: 0,
		pageEnd: entriesPerPage, lang1: lang1, lang2: lang2,
	}
}

func (mwl manageWordsList) view() string {
	var b strings.Builder
	if mwl.title != "" {
		b.WriteString(promptStyle.Render(mwl.title) + "\n\n")
	}

	for i := mwl.pageStart; i <= mwl.pageEnd && i < len(mwl.items); i++ {
		item := mwl.items[i]

		if i == mwl.cursor {
			b.WriteString(cursorStyle.Render("> "))
		} else {
			b.WriteString("  ")
		}
		wordLine := strconv.FormatInt(item.ID, 10) + ". "
		var score float64 // how well you know word in given direction

		switch mwl.lang1 + "->" + mwl.lang2 {
		case "en->es":
			wordLine += item.English + " -> " + item.Spanish
			score = item.EnToEsScore
		case "es->en":
			wordLine += item.Spanish + " -> " + item.English
			score = item.EsToEnScore
		case "en->uk":
			wordLine += item.English + " -> " + item.Ukrainian
			score = item.EnToUkScore
		case "uk->en":
			wordLine += item.Ukrainian + " -> " + item.English
			score = item.UkToEnScore
		case "es->uk":
			wordLine += item.Spanish + " -> " + item.Ukrainian
			score = item.EsToUkScore
		case "uk->es":
			wordLine += item.Ukrainian + " -> " + item.Spanish
			score = item.UkToEsScore
		}

		b.WriteString(scoreToKnowledgeLevelStyle(score).Render(wordLine))

		b.WriteString("\n")

	}

	// Pages subtext - // offset by 1 since page numbers start at 1, round up when dividing
	pageCounterSubText := ("Page " + strconv.Itoa((mwl.pageStart/entriesPerPage)+1) +
		" of " + strconv.Itoa((len(mwl.items)+entriesPerPage-1)/entriesPerPage))
	b.WriteString(pageCounterStyle.Render(pageCounterSubText))

	return b.String()
}


func (m *Model) buildManageWordsMenu() {
	m.manageWordsMenu = newManageWordsList("Select a word", m.allWords, m.selectedProfile.Lang1, m.selectedProfile.Lang2) 
}

// Manage words screen - show all words, allow for reset, removal (TODO)
func (m Model) viewManageWords() string {
	return m.manageWordsMenu.view() + helpStyle.Render("\n↑/↓ to navigate • enter to select • esc to go back • [PGUP/PGDN/HOME/END] to page by "+strconv.Itoa(entriesPerPage)+" words • [TAB] to swap direction • p to pronounce")
}

func (m Model) updateManageWords(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.screen = 	screenChooseQuizMode

			return m, nil
		}
	}
	return m, nil
}
