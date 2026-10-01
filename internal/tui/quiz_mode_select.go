package tui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)



func languageName(langCode string) string {
	switch langCode {
	case "en":
		return "English"
	case "es":
		return "Spanish"
	case "uk":
		return "Ukrainian"
	default:
		return langCode
	}
}

// components for choosing the quiz mode

func (m *Model) buildQuizModeMenu() {
	lang1 := languageName(m.selectedProfile.Lang1)
	lang2 := languageName(m.selectedProfile.Lang2)


	// TODO fill in with correct directions
	m.quizModeMenu = newChoiceList("Select a training direction", 
	[]choiceItem {
		newChoiceItem(fmt.Sprintf("Provide %s words, answer with %s words", lang1, lang2), "lang1_to_lang2"),
		newChoiceItem(fmt.Sprintf("Provide %s words, answer with %s words", lang2, lang1), "lang2_to_lang1"),
		newSeparatorItem(),
		newChoiceItem("Profile settings", profileSettings),
	})
}

func (m *Model) viewChooseQuizMode() string {
	var b strings.Builder
	b.WriteString(m.quizModeMenu.view())

	if m.updateProfileSuccess != "" {
		b.WriteString(successStyle.Render(m.updateProfileSuccess))
	}

	b.WriteString(helpStyle.Render("\n↑/↓ to navigate • enter to select • esc to go back"))

	return b.String()
}

func (m *Model) updateChooseQuizMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	// Allow selecting by number (1-based)
	n, err := strconv.Atoi(keyMsg.String())
	if err == nil && (n >= 1 && n <= 3) { // Add one for the Exit option

		m.quizModeMenu.cursor = n - 1

		keyMsg = tea.KeyMsg{Type: tea.KeyEnter} // continue as if "enter" was pressed
	}

	switch keyMsg.String() {
	case "up", "k":
		m.quizModeMenu.up()
	case "down", "j":
		m.quizModeMenu.down()
	case "enter":
		// m.updateProfileSuccess = "" // reset success message if it already exists

		selectableNumber := 0

		for i, item := range m.quizModeMenu.items {
			if !item.selectable {
				continue
			}

			selectableNumber++

			if selectableNumber == n {
				m.quizModeMenu.cursor = i
			}
		}

		if m.quizModeMenu.selected().value == profileSettings {
			m.buildSettingsMenu()
			m.screen = screenProfileSettings
			return m, nil
		}

		// quizModeLanguageDirection := m.quizModeMenu.selected().value

		m.buildQuizModeMenu()
		m.screen = screenProfileSelect // TODO
		return m, nil // TODO change to textinput.Blink
	case "esc":
		m.buildProfileMenu() // rebuild profile menu in case using new profile
		// m.updateProfileSuccess = ""
		m.screen = screenProfileSelect
	}

	return m, nil
}