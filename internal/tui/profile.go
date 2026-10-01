package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	exitValue = "__exit__"
	newProfileValue = "__new_profile__"
	maxProfileSlots = 8 // leave one for the exit, and limit to single digit selection
	profileSettings = "__profile_settings__"
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

	helpTxt := "↑/↓ to navigate • enter to select • r to rename a profile • [DEL]/'d' to delete a profile • q to quit"
	if m.profileRenaming {
		b.WriteString(m.viewProfileRename())
		helpTxt = "enter to confirm • esc to go back"
	} else {
		b.WriteString(m.profileMenu.view())
	}


	if m.profileActionErr != "" {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(m.profileActionErr))
	}

	b.WriteString(helpStyle.Render(helpTxt))
	return b.String()
 }

func (m *Model) viewProfileRename() string {
	var b strings.Builder

	for i, profile := range m.existingProfiles {
		cursor := "  "
		if i == m.profileRenameIndex {
			b.WriteString(selectedStyle.Render(">  " + strconv.Itoa(i+1) + "."))
		}

		if i == m.profileRenameIndex {
			fmt.Fprintf(&b, " %s\n", m.profileRenameInput.View(),
            )
		} else {
			fmt.Fprintf(&b, "%s %d. %s\n", cursor, i+1, profile.Name,
            )
		}
	}
	return b.String()
}

 // updateProfileSelect is the handler for the profile selection menu
func (m *Model) updateProfileSelect(msg tea.Msg) (tea.Model, tea.Cmd) { 
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if m.profileRenaming {
    	return m.updateProfileRename(msg)
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
	case "delete", "d":
		selected := m.profileMenu.selected()
		if selected.value == exitValue || selected.value == newProfileValue {
			m.profileActionErr = "Invalid deletion option selected"
			break
		}

		m.profileToDelete = selected.value        // profile to delete
		m.specialDeletePhraseInput.Reset()
		m.screen = screenDeleteProfileConfirm
	case "r":
		selected := m.profileMenu.selected()
		if selected.value == exitValue || selected.value == newProfileValue {
			m.profileActionErr = "Invalid rename option selected"
			break
		}

		m.profileRenaming = true
		m.profileRenameIndex = m.profileMenu.cursor

        m.profileRenameInput = textinput.New()
		m.profileRenameInput.Prompt = ""
        m.profileRenameInput.Placeholder = "New profile name"
        m.profileRenameInput.SetValue("")
        m.profileRenameInput.Focus()

		// m.profileMenu.renameIndex = m.profileRenameIndex
		// m.profileMenu.renameInput = &m.profileRenameInput

		return m, textinput.Blink
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


func (m *Model) updateProfileRename(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
			case "enter":
				name := strings.TrimSpace(m.profileRenameInput.Value())
				if !isValidProfileName(name) {
					m.profileActionErr = "Please enter a valid name for your profile"
					return m, nil
				}

				oldName := m.existingProfiles[m.profileRenameIndex].Name

				updatedProfile, err := m.store.RenameProfile(context.Background(), oldName, name)
				if err != nil {
					m.profileActionErr = err.Error()
					return m, nil
				}

				m.existingProfiles[m.profileRenameIndex] = updatedProfile
				m.profileRenaming = false
				m.buildProfileMenu()
				return m, nil
			case "esc":
				m.profileRenaming = false
				m.profileActionErr = ""

				return m, nil
		}
	}

	var cmd tea.Cmd
	m.profileRenameInput, cmd = m.profileRenameInput.Update(msg)
	return m, cmd
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