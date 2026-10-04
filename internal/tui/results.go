package tui

import (
	"fmt"
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

/* Components for the feedback screen */

func (m *Model) viewResults() string {
	var b strings.Builder
	pct := int(math.Round(float64(m.quiz.questionsCorrect) / float64(m.selectedProfile.DefaultNumQuestions) * 100.0))
	b.WriteString(titleStyle.Render(fmt.Sprintf("Quiz complete! You got %d/%d = %d%% correct.", m.quiz.questionsCorrect,
		 	m.selectedProfile.DefaultNumQuestions, pct)))
	b.WriteString("\n\n")

	if len(m.quiz.incorrectWords) == 0 {
		b.WriteString(successStyle.Render("Looks like you know all the words given!"))
	} else {
		b.WriteString(promptStyle.Render("Here are the words you got wrong for review:"))
		b.WriteString("\n")
		for _, w := range m.quiz.incorrectWords {
			wordMap := map[string]string{"en": w.English, "es": w.Spanish, "uk": w.Ukrainian}
			b.WriteString(fmt.Sprintf("    %s -> %s\n", wordMap[m.selectedProfile.Lang1], wordMap[m.selectedProfile.Lang2]))
		}
	}


	b.WriteString(helpStyle.Render("\nesc/enter to continue"))
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
	} // TODO cycle between wrong answers, pronounce

	return m, nil
}