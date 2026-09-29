package tui

import (
	"strconv"
	"strings"

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


	b.WriteString(helpStyle.Render("\n↑/↓ to navigate • enter to select • r to rename a profile • [DEL]/'d' to delete a profile • q to quit"))
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
	case "q", "esc":
		return m, tea.Quit
	}
	return m, nil
}
