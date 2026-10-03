package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/alexbudy/go_spanish_rewrite/internal/store"
	tea "github.com/charmbracelet/bubbletea"
)

/* components needed for word details screen */

type wordDetails struct {
	title         string
	word          store.Word
	lang1 string
	lang2 string
	profileId     int64        // profile id for associated rankings
	profile       string       // profile name for associated rankings
	wordUpdateMsg string       // if word was updated (score reset), show this msg
}

func (wd *wordDetails) toggleSelectedLocale() {
	wd.lang1, wd.lang2 = wd.lang2, wd.lang1
}

func (wd *wordDetails) resetScore(s *store.Store) store.Word {
	switch wd.lang1 + "->" + wd.lang2 {
		case "en->es":
			wd.word.EnToEsScore = 0.0
		case "es->en":
			wd.word.EsToEnScore = 0.0
		case "en->uk":
			wd.word.EnToUkScore = 0.0
		case "uk->en":
			wd.word.UkToEnScore = 0.0
		case "es->uk":
			wd.word.EsToUkScore = 0.0
		case "uk->es":
			wd.word.UkToEsScore = 0.0
	}
	
	wd.wordUpdateMsg = "Score reset successfully"
	s.ResetRankingsForWord(context.Background(), wd.word.ID, wd.lang1, wd.lang2, wd.profileId)
	return wd.word
}

func (wd *wordDetails) swapLangs() {
	wd.lang1, wd.lang2 = wd.lang2, wd.lang1
}


func newWordDetails(title string, word store.Word, profile store.Profile) wordDetails {
	return wordDetails{title: title, word: word, lang1: profile.Lang1, lang2: profile.Lang2, profile: profile.Name, profileId: profile.ID}
}

func (m *Model) buildWordDetailsMenu() {
	m.wordDetailsMenu = newWordDetails("Word Details for ", m.manageWordsMenu.items[m.manageWordsMenu.cursor], m.selectedProfile) // default to English 'from' selection
}

func (wd *wordDetails) view() string {
	var b strings.Builder
	var langCodeToWord = map[string]string{"en": wd.word.English, "es": wd.word.Spanish, "uk": wd.word.Ukrainian}

	wd1 := langCodeToWord[wd.lang1]
	wd2 := langCodeToWord[wd.lang2]

	fmt.Fprintf(&b, "Word details for '%s' <-> '%s'\n", wd1, wd2)
	fmt.Fprintf(&b, "ID: %d\n", wd.word.ID)

	// TODO - simplify
	if wd.lang1 == "en" && wd.lang2 == "es" {
		b.WriteString(cursorStyle.Render("> "))
		b.WriteString(scoreToKnowledgeLevelStyle(wd.word.EnToEsScore).Render(wd.word.English + " -> " + wd.word.Spanish + ": " + strconv.FormatFloat(wd.word.EnToEsScore, 'f', 2, 64)))
		b.WriteString("\n")
		b.WriteString(scoreToKnowledgeLevelStyle(wd.word.EsToEnScore).Render("  " + wd.word.Spanish + " -> " + wd.word.English + ": " + strconv.FormatFloat(wd.word.EsToEnScore, 'f', 2, 64)))
		b.WriteString("\n")
	} else  if wd.lang1 == "es" && wd.lang2 == "en" {
		b.WriteString(cursorStyle.Render("> "))
		b.WriteString(scoreToKnowledgeLevelStyle(wd.word.EsToEnScore).Render(wd.word.Spanish + " -> " + wd.word.English + ": " + strconv.FormatFloat(wd.word.EsToEnScore, 'f', 2, 64)))
		b.WriteString("\n")
		b.WriteString(scoreToKnowledgeLevelStyle(wd.word.EnToEsScore).Render("  " + wd.word.English + " -> " + wd.word.Spanish + ": " + strconv.FormatFloat(wd.word.EnToEsScore, 'f', 2, 64)))
		b.WriteString("\n")
	} else if wd.lang1 == "en" && wd.lang2 == "uk" {
		b.WriteString(cursorStyle.Render("> "))
		b.WriteString(scoreToKnowledgeLevelStyle(wd.word.EnToUkScore).Render(wd.word.English + " -> " + wd.word.Ukrainian + ": " + strconv.FormatFloat(wd.word.EnToUkScore, 'f', 2, 64)))
		b.WriteString("\n")
		b.WriteString(scoreToKnowledgeLevelStyle(wd.word.UkToEnScore).Render("  " + wd.word.Ukrainian + " -> " + wd.word.English + ": " + strconv.FormatFloat(wd.word.UkToEnScore, 'f', 2, 64)))
		b.WriteString("\n")
	} else if wd.lang1 == "uk" && wd.lang2 == "en" {
		b.WriteString(cursorStyle.Render("> "))
		b.WriteString(scoreToKnowledgeLevelStyle(wd.word.UkToEnScore).Render(wd.word.Ukrainian + " -> " + wd.word.English + ": " + strconv.FormatFloat(wd.word.UkToEnScore, 'f', 2, 64)))
		b.WriteString("\n")
		b.WriteString(scoreToKnowledgeLevelStyle(wd.word.EnToUkScore).Render("  " + wd.word.English + " -> " + wd.word.Ukrainian + ": " + strconv.FormatFloat(wd.word.EnToUkScore, 'f', 2, 64)))
		b.WriteString("\n")
	} else if wd.lang1 == "es" && wd.lang2 == "uk" {
		b.WriteString(cursorStyle.Render("> "))
		b.WriteString(scoreToKnowledgeLevelStyle(wd.word.EsToUkScore).Render(wd.word.Spanish + " -> " + wd.word.Ukrainian + ": " + strconv.FormatFloat(wd.word.EsToUkScore, 'f', 2, 64)))
		b.WriteString("\n")
		b.WriteString(scoreToKnowledgeLevelStyle(wd.word.UkToEsScore).Render("  " + wd.word.Ukrainian + " -> " + wd.word.Spanish + ": " + strconv.FormatFloat(wd.word.UkToEsScore, 'f', 2, 64)))
		b.WriteString("\n")
	} else if wd.lang1 == "uk" && wd.lang2 == "es" {
		b.WriteString(cursorStyle.Render("> "))
		b.WriteString(scoreToKnowledgeLevelStyle(wd.word.UkToEsScore).Render(wd.word.Ukrainian + " -> " + wd.word.Spanish + ": " + strconv.FormatFloat(wd.word.UkToEsScore, 'f', 2, 64)))
		b.WriteString("\n")
		b.WriteString(scoreToKnowledgeLevelStyle(wd.word.EsToUkScore).Render("  " + wd.word.Spanish + " -> " + wd.word.Ukrainian + ": " + strconv.FormatFloat(wd.word.EsToUkScore, 'f', 2, 64)))
		b.WriteString("\n")
	}

	if wd.wordUpdateMsg != "" {
		b.WriteString(successStyle.Render(wd.wordUpdateMsg))
		b.WriteString("\n")
	}

	return b.String()
}

// Manage words screen - show all words, allow for reset, removal (TODO?)
func (m Model) viewWordDetails() string {
	return m.wordDetailsMenu.view() + helpStyle.Render("\n ↑/↓ to navigate • esc to go back • p to pronounce • ctrl+r to reset selected score")
}

func (m *Model) updateWordDetails(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "up", "k", "down", "j", "tab":
		m.wordDetailsMenu.swapLangs()
	case "p":
		pronounceWordPair(m.wordDetailsMenu.word, m.wordDetailsMenu.lang1, m.wordDetailsMenu.lang2)
	case "ctrl+r":
		updatedWord := m.wordDetailsMenu.resetScore(m.store)

		m.manageWordsMenu.items[m.manageWordsMenu.cursor] = updatedWord
	case "esc":
		m.screen = screenManageWords
	}

	return m, nil
}