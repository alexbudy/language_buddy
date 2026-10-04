package tui

import (
	"fmt"
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

/* Components for the feedback screen */

func (m *Model) buildResultsMenu(){
	items := []choiceItem{}

	for _, w := range m.quiz.incorrectWords {
		wordMap := map[string]string{"en": w.English, "es": w.Spanish, "uk": w.Ukrainian}
		label := wordMap[m.selectedProfile.Lang1] + " -> " + wordMap[m.selectedProfile.Lang2]
		value := m.selectedProfile.Lang1 + "_" + m.selectedProfile.Lang2
		items = append(items, newChoiceItem(label, value))
	}
	m.resultsMenu = newChoiceList(" ", items)

}

func (m *Model) viewResults() string {
	var b strings.Builder
	pct := int(math.Round(float64(m.quiz.questionsCorrect) / float64(m.selectedProfile.DefaultNumQuestions) * 100.0))
	b.WriteString(titleStyle.Render(fmt.Sprintf("Quiz complete! You got %d/%d = %d%% correct.", m.quiz.questionsCorrect,
		 	m.selectedProfile.DefaultNumQuestions, pct)))
	b.WriteString("\n\n")

	if len(m.quiz.incorrectWords) == 0 {
		b.WriteString(successStyle.Render("Looks like you know all the words given!"))
	} else {
		b.WriteString(promptStyle.Render("Here are the words you got wrong for review:\n"))
		b.WriteString(m.resultsMenu.view())
	}


	b.WriteString(helpStyle.Render("\nesc/enter to continue • p to hear a word phrase"))
	return b.String()
}

func (m *Model) updateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg); 
	if !ok {
		return m, tea.Quit
	}
	switch keyMsg.String() {
	case "enter", "esc":
		m.buildQuizModeMenu()
		m.screen = screenChooseQuizMode
		return m, nil
	case "up", "k":
		m.resultsMenu.up()
	case "down", "j":
		m.resultsMenu.down()
	case "p":
		words := strings.Split(m.resultsMenu.selected().label, " -> ") // ex "time -> el tiempo"
		langs := strings.Split(m.resultsMenu.selected().value, "_") // ex "en_es"

		pronounceWordPairStringsOnly(strings.TrimSpace(words[0]), strings.TrimSpace(words[1]), langs[0], langs[1])

	} 


	return m, nil
}