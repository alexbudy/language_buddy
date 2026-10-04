// Package tui provides the terminal user interface for the application.
package tui

import (
	"context"

	"github.com/alexbudy/go_spanish_rewrite/internal/log"
	"github.com/alexbudy/go_spanish_rewrite/internal/store"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// New creates a TUI model backed by s.
func New(s *store.Store) Model {
	m := Model{store: s, newProfileInput:
		createProfileInput(), 
		specialDeletePhraseInput: createSpecialDeletePhraseInput()}

	// load profiles into m.profileMenu
	m.loadProfiles()

	return m
}

func createProfileInput() textinput.Model {
	t1 := textinput.New()
	t1.Placeholder = "Enter a name for your profile"
	t1.CharLimit = 40
	t1.Width = 40

	return t1
}

// loads profiles from DB, attaches to the model, builds the menu, and sets the screen
func (m *Model) loadProfiles() {
	ctx := context.Background()

	profiles, err := m.store.GetProfiles(ctx)
	if err != nil {
		log.Error("Failed to get profiles: %v", err)
		return
	}
	log.Debug("Got %d profiles from the DB", len(profiles))
	m.existingProfiles = profiles
	m.buildProfileMenu()
	m.screen = screenProfileSelect
}

// Init implements [tea.Model].
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements [tea.Model].
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+c":
			log.Debug("Ctrl+c pressed, exiting program")
			return m, tea.Quit
		}
	}

	switch m.screen {
		case screenProfileSelect:
			return m.updateProfileSelect(msg)
		case screenNewProfile:
			return m.updateNewProfile(msg)
		case screenDeleteProfileConfirm:
			return m.updateDeleteProfileConfirm(msg)
		case screenChooseQuizMode:
			return m.updateChooseQuizMode(msg)
		case screenQuestion:
			return m.updateQuestion(msg)
		case screenProfileSettings:
			return m.updateProfileSettings(msg)
		case screenManageWords:
			return m.updateManageWords(msg)
		case screenWordDetails:
			return m.updateWordDetails(msg)
		case screenResults:
			return m.updateResults(msg)
		case screenGoodbye:
			return m, tea.Quit
	}
	return m, nil
}

// View implements [tea.Model].
func (m Model) View() string {
	switch m.screen {
	case screenProfileSelect:
		return m.viewProfileSelect()
	case screenNewProfile:
		return m.viewNewProfile()
	case screenDeleteProfileConfirm:
		return m.viewDeleteProfileConfirm()
	case screenChooseQuizMode:
		return m.viewChooseQuizMode()
	case screenQuestion:
		return m.viewQuestion()
	case screenProfileSettings:
		return m.viewProfileSettings()
	case screenManageWords:
		return m.viewManageWords()
	case screenWordDetails:
		return m.viewWordDetails()
	case screenResults:
		return m.viewResults()
	case screenGoodbye:
		return titleStyle.Render("Exiting the program, thanks for training!") + "\n"
	}
	return "Shouldn't come here" // shouldnt come here
} 