package tui

import (
	"context"
	"strings"

	"github.com/alexbudy/go_spanish_rewrite/internal/log"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// components for deleting profile screen

func createSpecialDeletePhraseInput() textinput.Model {
	t1 := textinput.New()
	t1.Placeholder = "type deletion phrase"
	t1.Width = 50
	t1.CharLimit = 40

	t1.Focus()

	return t1
}

func (m Model) viewDeleteProfileConfirm() string {
	deletionPhrase := "delete " + m.profileToDelete // phrase to type to delete profile

	var b strings.Builder
	b.WriteString(deleteConfirmStyle.Render("Type '"))
	b.WriteString(deleteConfirmSpecialWordStyle.Render(deletionPhrase))
	b.WriteString(deleteConfirmStyle.Render("'to confirm deletion"))
	b.WriteString("\n\n")
	b.WriteString(m.specialDeletePhraseInput.View())

	if m.invalidDeletePhraseErr != "" {
		b.WriteString("\n\n")
		b.WriteString(errorStyle.Render(m.invalidDeletePhraseErr))
	}

	b.WriteString(helpStyle.Render("\n\nenter to confirm • esc to go back"))
	return b.String()
}
// TODO
func (m *Model) updateDeleteProfileConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	deletionPhrase := "delete " + m.profileToDelete

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "enter":
		if deletionPhrase != m.specialDeletePhraseInput.Value() {
			m.invalidDeletePhraseErr = "Incorrect deletion phrase"
			return m, nil
		}

		if err := m.store.DeleteProfile(context.Background(), m.profileToDelete); err != nil {
			log.Error("Failed deleting profile %s",m.profileToDelete)
			return m, nil
		}

		// remove the deleted profile from existing profiles
		for i, existingProfile := range m.existingProfiles {
			if existingProfile.Name == m.profileToDelete {
				m.existingProfiles = append(m.existingProfiles[:i], m.existingProfiles[i+1:]...)
				break
			}
		}

		m.buildProfileMenu() // rebuild profile menu in case using new profile
		m.screen = screenProfileSelect

		return m, nil
	case "esc":
		m.buildProfileMenu()

		m.screen = screenProfileSelect
		return m, nil
	}

	var cmd tea.Cmd
	m.specialDeletePhraseInput, cmd = m.specialDeletePhraseInput.Update(msg)

	return m, cmd
}