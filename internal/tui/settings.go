package tui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// components for the settings screen


type saveCancelSelection int
const (
	saveSelected saveCancelSelection = iota
	cancelSelected
)

type profileSetting int
const (
	settingTTS profileSetting = iota
	settingLanguage1
	settingLanguage2
	settingNumQuestions
	settingNumAnswerOptions
	settingSaveCancel
)

type profileSettingsConfig struct {
	title string
	profile string // name of profile
	enableTTS bool
	language1 string
	language2 string
	defaultNumQuestions int
	defaultNumAnswers int
	
	selectedSetting profileSetting

	saveCancelSelection saveCancelSelection
}

var languages = []string{"en", "es", "uk"}
var langCodeToName = map[string]string{
	"en": "English",
	"es": "Spanish",
	"uk": "Ukrainian",
}

func newProfileSettings(title string, profile string, allowTTS bool, defaultNumQuestions int, 
						defaultNumAnswerOptions int, lang1 string, lang2 string) profileSettingsConfig {
	return profileSettingsConfig{
		title: title, profile: profile, enableTTS: allowTTS,
		defaultNumQuestions:     defaultNumQuestions,
		defaultNumAnswers: defaultNumAnswerOptions,
		selectedSetting:         settingTTS,
		language1:         lang1,
		language2:         lang2,
	}
}

func (psc *profileSettingsConfig) toggleEnableTTS() {
	psc.enableTTS = !psc.enableTTS
}

func (psc *profileSettingsConfig) cycleLanguage(lang1OrLang2 string) {
	currentLang := psc.language1
	otherLang := psc.language2 // assume lang1 is being cycled
	if lang1OrLang2 == "lang2" {
		otherLang = psc.language1
		currentLang = psc.language2
	}

	curIdx := 0
	for i, lang := range languages {
		if lang == currentLang {
			curIdx = i 
			break // found the 
		}
	}

	updatedLang := ""
	for i := 0; i < len(languages); i++ {
		nextIdx := (curIdx + 1 + i) % len(languages)

		if languages[nextIdx] != otherLang && languages[nextIdx] != lang1OrLang2 {

			updatedLang = languages[nextIdx]
			break
		}
	}

	if lang1OrLang2 == "lang1" {
		psc.language1 = updatedLang
	} else {
		psc.language2 = updatedLang
	}
}


func (psc *profileSettingsConfig) up() {
	switch psc.selectedSetting {
	case settingTTS:
		psc.selectedSetting = settingSaveCancel
	
	case settingLanguage1:
		psc.selectedSetting = settingTTS

	case settingLanguage2:
		psc.selectedSetting = settingLanguage1

	case settingNumQuestions:
		psc.selectedSetting = settingLanguage2

	case settingNumAnswerOptions:
		psc.selectedSetting = settingNumQuestions

	case settingSaveCancel:
		psc.selectedSetting = settingNumAnswerOptions
	}
}

func (psc *profileSettingsConfig) down() {
	switch psc.selectedSetting {
	case settingTTS:
		psc.selectedSetting = settingLanguage1

	case settingLanguage1:
		psc.selectedSetting = settingLanguage2

	case settingLanguage2:
		psc.selectedSetting = settingNumQuestions

	case settingNumQuestions:
		psc.selectedSetting = settingNumAnswerOptions

	case settingNumAnswerOptions:
		psc.selectedSetting = settingSaveCancel

	case settingSaveCancel:
		psc.selectedSetting = settingTTS
	}
}


func (psc *profileSettingsConfig) increaseQuestions() {
	if psc.defaultNumQuestions < 20 {
		psc.defaultNumQuestions++
	}
}

func (psc *profileSettingsConfig) decreaseQuestions() {
	if psc.defaultNumQuestions > 1 {
		psc.defaultNumQuestions--
	}
}

func (psc *profileSettingsConfig) increaseNumAnswers() {
	if psc.defaultNumAnswers < 6 {
		psc.defaultNumAnswers++
	}
}

func (psc *profileSettingsConfig) decreaseNumAnswers() {
	if psc.defaultNumAnswers > 2 {
		psc.defaultNumAnswers--
	}
}

func (m *Model) buildSettingsMenu() {
	m.profileSettings = newProfileSettings(
		"Profile Settings",
		m.selectedProfile.Name,
		m.selectedProfile.EnableSpeech,
		m.selectedProfile.DefaultNumQuestions,
		m.selectedProfile.DefaultNumAnswers,
		m.selectedProfile.Lang1,
		m.selectedProfile.Lang2,
	)
}

func (psc profileSettingsConfig) view() string {
	var b strings.Builder

	b.WriteString(promptStyle.Render(psc.title))
	b.WriteString("\n\n")

	label := "Allow TTS: "

	if psc.selectedSetting == settingTTS {
		b.WriteString(selectedStyle.Render(" > " + label))
	} else {
		b.WriteString("   " + label)
	}
	if psc.enableTTS {
		b.WriteString(settingSelectionStyle.Render("YES"))
	} else {
		b.WriteString(settingSelectionStyle.Render("NO"))
	}

	b.WriteString("\n")

	label = "Language 1: "
	if psc.selectedSetting == settingLanguage1 {
		b.WriteString(selectedStyle.Render(" > " + label))
	} else {
		b.WriteString("   " + label)
	}
	
	b.WriteString(settingSelectionStyle.Render(langCodeToName[psc.language1]))
	b.WriteString("\n")

	label = "Language 2: "
	if psc.selectedSetting == settingLanguage2 {
		b.WriteString(selectedStyle.Render(" > " + label))
	} else {
		b.WriteString("   " + label)
	}
	b.WriteString(settingSelectionStyle.Render(langCodeToName[psc.language2]))
	b.WriteString("\n")

	label = "Number of Questions per Quiz: "
	if psc.selectedSetting == settingNumQuestions {
		b.WriteString(selectedStyle.Render(" > " + label))
	} else {
		b.WriteString("   " + label)
	}
	b.WriteString(settingSelectionStyle.Render(strconv.Itoa(psc.defaultNumQuestions)))

	b.WriteString("\n")

	label = "Number of Answer Choices per Question: "
	if psc.selectedSetting == settingNumAnswerOptions {
		b.WriteString(selectedStyle.Render(" > " + label))
	} else {
		b.WriteString("   " + label)
	}
	b.WriteString(settingSelectionStyle.Render(strconv.Itoa(psc.defaultNumAnswers)))

	b.WriteString("\n\n")

	if psc.selectedSetting == settingSaveCancel {
		if psc.saveCancelSelection == saveSelected {
			b.WriteString(selectedStyle.Render("> SAVE "))
			b.WriteString("   CANCEL")
		} else {
			b.WriteString("  SAVE  ")
			b.WriteString(selectedStyle.Render("> CANCEL"))
		}
	} else {
		b.WriteString("  SAVE    CANCEL")
	}

	return b.String()
}

func (m *Model) viewProfileSettings() string {
	return m.profileSettings.view() + 
		helpStyle.Render("\n↑/↓/⇆ to navigate • enter to select • esc to go back")
}

func (m *Model) updateProfileSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)

	if !ok {
		return m, nil
	}

		switch keyMsg.String() {
	case "up", "k":
		m.profileSettings.up()
	case "down", "j":
		m.profileSettings.down()
	// case "left", "h":
	// 	if m.profileSettings.selectedSetting == settingNumQuestions {
	// 		m.profileSettings.decreaseQuestions()
	// 	} else if m.profileSettings.selectedSetting == settingTTS {
	// 		m.profileSettings.toggleEnableTTS()
	// 	} else if m.profileSettings.selectedSetting == settingLanguage {
	// 		m.profileSettings.cycleDefaultLanguage()
	// 	} else if m.profileSettings.selectedSetting == settingNumAnswerOptions {
	// 		m.profileSettings.decreaseNumAnswers()
	// 	} else if m.profileSettings.selectedSetting == settingSaveCancel {
	// 		m.profileSettings.saveCancelSelection = saveSelected
	// 	}
	// case "right", "l":
	// 	if m.profileSettings.selectedSetting == settingNumQuestions {
	// 		m.profileSettings.increaseQuestions()
	// 	} else if m.profileSettings.selectedSetting == settingTTS {
	// 		m.profileSettings.toggleEnableTTS()
	// 	} else if m.profileSettings.selectedSetting == settingLanguage {
	// 		m.profileSettings.cycleDefaultLanguage()
	// 	} else if m.profileSettings.selectedSetting == settingNumAnswerOptions {
	// 		m.profileSettings.increaseNumAnswers()
	// 	} else if m.profileSettings.selectedSetting == settingSaveCancel {
	// 		m.profileSettings.saveCancelSelection = cancelSelected
	// 	}
	// case "tab":
	// 	if m.profileSettings.selectedSetting == settingTTS {
	// 		m.profileSettings.toggleEnableTTS()
	// 	} else if m.profileSettings.selectedSetting == settingLanguage {
	// 		m.profileSettings.cycleDefaultLanguage()
	// 	}
	// case "enter":
	// 	if m.profileSettings.selectedSetting == settingSaveCancel {
	// 		if m.profileSettings.saveCancelSelection == saveSelected {
	// 			err := m.store.UpdateProfile(m.ctx,
	// 				store.Profile{
	// 					Name: m.profile, EnableSpeech: m.profileSettings.enableTTS,
	// 					DefaultNumQuestions: m.profileSettings.defaultNumQuestions,
	// 					DefaultNumAnswers:   m.profileSettings.defaultNumAnswerOptions,
	// 					DefaultLanguage:     m.profileSettings.defaultLanguage,
	// 				})
	// 			if err != nil {
	// 				fmt.Errorf("store: check some_new_column: %w", err)
	// 				return m, nil
	// 			}

	// 			m.updateProfileSuccess = "Profile settings updated successfully"
	// 			m.screen = screenDirectionSelect
	// 		} else {
	// 			m.screen = screenDirectionSelect
	// 		}
	// 	}
	case "esc":
		// m.updateProfileSuccess = ""
		m.buildQuizModeMenu()
		m.screen = screenChooseQuizMode
	}
	return m, nil
}