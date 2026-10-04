package tui

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/alexbudy/go_spanish_rewrite/internal/log"
	"github.com/alexbudy/go_spanish_rewrite/internal/store"
	tea "github.com/charmbracelet/bubbletea"
)

/* quiz components */

type quizState struct {
	profile store.Profile

	questionIndex     int // 1-based index of the question currently shown/answered
	questionsCorrect  int
	questionedWordIDs []int64
	incorrectWords    []store.Word
	answerPool        []string

	target         store.Word
	targetText     string
	correctAnswer  string
	selectedAnswer string
	wasCorrect     bool
}

func (m *Model) startQuiz() error {
	answerPool := make([]string, 0, len(m.allWords))

	for _, word := range m.allWords {
		answerPool = append(answerPool, word.GetTextFromLang(m.selectedProfile.Lang2))
	}

	m.quiz = quizState{
		profile: m.selectedProfile,
		questionsCorrect: 0,
		questionIndex: 1,
		answerPool: answerPool,
	}

	return m.loadNextQuestion() 

}

func (m *Model) loadNextQuestion() error {
	// grab ten words by quiz mode (difficulty)
	words, err := m.store.GetWordsForQuestion(
		context.Background(), m.selectedProfile.Name,
		m.selectedProfile.QuizMode, m.selectedProfile.Lang1, m.quiz.questionedWordIDs)

	if err != nil {
		return err
	}
	if len(words) == 0 {
		return fmt.Errorf("no more words available for %q", m.selectedProfile.Name)
	}

	target := words[rand.Intn(len(words))] // choose random Word to quiz on
	targetWordText := target.GetTextFromLang(m.selectedProfile.Lang1) // word to question, in target language
    correctAnswer := target.GetTextFromLang(m.selectedProfile.Lang2)

	answerOptions := sampleDistinct(m.quiz.answerPool, m.selectedProfile.DefaultNumQuestions)
	answerOptions = removeFirst(answerOptions, correctAnswer) // want to remove the correct answer if its there

	if len(answerOptions) > m.selectedProfile.DefaultNumAnswers - 1 {
		answerOptions = answerOptions[:m.selectedProfile.DefaultNumAnswers-1] // remove one word to allow for the correct answer
	}

	answerSet := append(append([]string{}, answerOptions...), correctAnswer) // append the correct answer
	rand.Shuffle(len(answerSet), func(i, j int) { answerSet[i], answerSet[j] = answerSet[j], answerSet[i] }) // randomly shuffle answer set

	m.quiz.target = target // target Word
	m.quiz.targetText = targetWordText
	m.quiz.correctAnswer = correctAnswer
	title := fmt.Sprintf("%d/%d Please select the translation for %q", m.quiz.questionIndex, m.selectedProfile.DefaultNumQuestions, targetWordText)
	m.answerMenu = newChoiceList(title, labeledItems(answerSet...))

	pronounceWord(target, m.selectedProfile.Lang1)

	m.screen = screenQuestion
	return err
}	

// sampleDistinct returns up to n distinct random elements from pool, without
// replacement, mirroring Python's random.sample.
func sampleDistinct(pool []string, n int) []string {
	if n > len(pool) {
		n = len(pool)
	}
	perm := rand.Perm(len(pool))[:n]
	out := make([]string, n)
	for i, idx := range perm {
		out[i] = pool[idx]
	}
	return out
}

// removeFirst returns a copy of s with the first occurrence of val removed.
func removeFirst(s []string, val string) []string {
	for i, v := range s {
		if v == val {
			out := make([]string, 0, len(s)-1)
			out = append(out, s[:i]...)
			out = append(out, s[i+1:]...)
			return out
		}
	}
	return s
}



func (m *Model) updateQuestion(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

		// Allow selecting an answer by number (1-based)
	n, err := strconv.Atoi(keyMsg.String())
	if err == nil && n >= 1 && n <= len(m.answerMenu.items) {
		m.answerMenu.cursor = n - 1
		keyMsg = tea.KeyMsg{Type: tea.KeyEnter} // continue as if "enter" was pressed
	}

	switch keyMsg.String() {
	case "up", "k":
		m.answerMenu.up()
	case "down", "j":
		m.answerMenu.down()
	case "enter":
		selected := m.answerMenu.selected().value
		correct := selected == m.quiz.correctAnswer

		adjustment := 0.15 * float64(m.selectedProfile.DefaultNumAnswers)
		if !correct {
			adjustment *= -1
		}
		if err := m.store.UpdateRankingForWord(context.Background(), m.quiz.target.ID, 
				m.selectedProfile.ID, m.selectedProfile.Lang1, m.selectedProfile.Lang2, adjustment); err != nil {
			log.Error("Failed to update word ranking: %v", err) // TODO
			return m, nil
		}

		if correct {
			m.quiz.questionsCorrect++

			// pronounce the correct translation once if correct
			pronounceWord(m.quiz.target, m.selectedProfile.Lang2)
		} else {
			m.quiz.incorrectWords = append(m.quiz.incorrectWords, m.quiz.target)
		}

		// TODO update quizState
		m.screen = screenQuestion
	case "p":
		pronounceWord(m.quiz.target, m.selectedProfile.Lang1)
	case "esc":
		m.screen = screenChooseQuizMode
	}

	return m, nil
}

func (m *Model) viewQuestion() string {
	return m.answerMenu.view() + helpStyle.Render("\n↑/↓ to navigate • enter to select • p to hear word again • esc to exit quiz")

}