// Package tui provides the terminal user interface for the application.
package tui

import (
	"context"

	"github.com/alexbudy/go_spanish_rewrite/internal/log"
	"github.com/alexbudy/go_spanish_rewrite/internal/store"
	tea "github.com/charmbracelet/bubbletea"
)

// New creates a TUI model backed by s.
func New(s *store.Store) Model {
	m := Model{store: s}

	// load profiles into m.profileMenu
	m.loadProfiles()

	return m
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
	case screenGoodbye:
		return titleStyle.Render("Exiting the program, thanks for training!") + "\n"
	}
	return "Shouldn't come here" // shouldnt come here
} 