package tui

import (
	"context"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	exitValue = "__exit__"
	newProfileValue = "__new_profile__"
	maxProfileSlots = 8 // leave one for the exit, and limit to single digit selection
)



// Methods for profile management

func (m *Model) buildProfileMenu() {
	var items []choiceItem
	for _, profile := range m.existingProfiles {
		items = append(items, newChoiceItem(profile.Name, profile.Name))
	}

	// add a visual line break to the last element, if it exists
	if len(items) > 0 {
		items = append(items, newSeparatorItem())
	}

	// add one new profile entry for dynamic profile creation
	if len(items) < maxProfileSlots {
		items = append(items, newChoiceItem("New Profile", newProfileValue))
	}

	items = append(items, newChoiceItem("Exit :(", exitValue))
	m.profileMenu = newChoiceList("Select or create a profile", items)

}

func (m *Model) viewProfileSelect() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Welcome to Spanish Buddy!"))
	b.WriteString("\n")
	b.WriteString(subtleStyle.Render("Practice your Spanish, English, and Ukrainian translations."))
	b.WriteString("\n\n")
	b.WriteString(m.profileMenu.view())


	b.WriteString(helpStyle.Render("↑/↓ to navigate • enter to select • r to rename a profile • [DEL]/'d' to delete a profile • q to quit"))
	return b.String()
 }

 // updateProfileSelect is the handler for the profile selection menu
func (m *Model) updateProfileSelect(msg tea.Msg) (tea.Model, tea.Cmd) { 
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	// Allow selecting a profile or option (new profile, exit) by number (1-based), skip non-selectable items
	n, err := strconv.Atoi(keyMsg.String())
	if err == nil { // 1-9 was pressed
		if n >= 1 && n <= len(m.profileMenu.items) {
			selectableNumber := 0

			for i, item := range m.profileMenu.items {
				if !item.selectable {
					continue
				}
				selectableNumber++
				if selectableNumber == n {
					m.profileMenu.cursor = i
				}
			}
		}
	}

	switch keyMsg.String() {
	case "up", "k":
		m.profileMenu.up()
	case "down", "j":
		m.profileMenu.down()
	case "enter":
		switch selected := m.profileMenu.selected(); selected.value {
		case exitValue:
			m.screen = screenGoodbye
		case newProfileValue:
			m.newProfileInput.Reset()
			m.newProfileInput.Focus()
			m.newProfileErr = ""
			m.screen = screenNewProfile
			return m, textinput.Blink
		}
	case "q", "esc":
		return m, tea.Quit
	}
	return m, nil
}


// screenNewProfile functions
func (m Model) viewNewProfile() string {
	var b strings.Builder
	b.WriteString(promptStyle.Render("What would you like to call this profile?"))
	b.WriteString("\n\n")
	b.WriteString(m.newProfileInput.View())
	if m.newProfileErr != "" {
		b.WriteString("\n\n")
		b.WriteString(errorStyle.Render(m.newProfileErr))
	}
	b.WriteString(helpStyle.Render("\n\nenter to confirm • esc to go back"))
	return b.String()
}

func (m *Model) updateNewProfile(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
			case "enter":
				name := m.newProfileInput.Value()
				if !isValidProfileName(name) {
					m.newProfileErr = "Please enter a valid name for your profile" + name
					return m, nil
				}

				profile, err := m.store.StoreProfile(context.Background(), name)

				if err != nil {
					m.newProfileErr = "Something went wrong creating the profile: " + err.Error()
					return m, nil
				}

				m.existingProfiles = append(m.existingProfiles, profile)
				m.buildProfileMenu() // rebuild profile menu with new profile
				m.screen = screenProfileSelect
				return m, nil
			case "esc":
				m.screen = screenProfileSelect
				return m, nil
			}
		}

	var cmd tea.Cmd
	m.newProfileInput, cmd = m.newProfileInput.Update(msg)

	return m, cmd
}

// Somewhat arbitrary profile name validation:
// 3 letters min, have one char, spaces, dashes, underscores ok, everything else not
func isValidProfileName(name string) bool {
	if len(name) < 3 {
		return false
	}
	hasLetter := false

	for _, c := range name {
		switch {
		case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z'):
			hasLetter = true
		case c >= '0' && c <= '9':
			// allowed
		case c == ' ', c == '-', c == '_':
			// allowed
		default:
			return false
		}

	}

	return hasLetter
}