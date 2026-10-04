package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/alexbudy/go_spanish_rewrite/internal/store"
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
	settingQuizMode
	settingNumQuestions
	settingNumAnswerOptions
	settingSaveCancel
)
const totalSettings profileSetting = 7

type profileSettingsConfig struct {
	title string
	profile string // name of profile
	enableTTS bool
	language1 string
	language2 string
	quizMode quizMode
	defaultNumQuestions int
	defaultNumAnswers int
	
	selectedSetting profileSetting

	saveCancelSelection saveCancelSelection
}

type quizMode string

const (
	// QuizMode keeps track of the quiz mode setting
	quizModeAny quizMode = "any"
	quizModeWellKnown quizMode = "well_known"
	quizModeLeastKnown quizMode = "least_known"
)
var orderedQuizMode = []quizMode{quizModeAny, quizModeWellKnown, quizModeLeastKnown}
var quizModeToString = map[quizMode]string{
	quizModeAny: "Any Words",
	quizModeWellKnown: "Well Known Words",
	quizModeLeastKnown: "Least Known Words",
}

var languages = []string{"en", "es", "uk"}
var langCodeToName = map[string]string{
	"en": "English",
	"es": "Spanish",
	"uk": "Ukrainian",
}

func newProfileSettings(title string, profile string, allowTTS bool, defaultNumQuestions int, 
						defaultNumAnswerOptions int, lang1 string, lang2 string, quizMode quizMode) profileSettingsConfig {
	return profileSettingsConfig{
		title: title, profile: profile, enableTTS: allowTTS,
		defaultNumQuestions:     defaultNumQuestions,
		defaultNumAnswers: defaultNumAnswerOptions,
		selectedSetting:         settingTTS,
		language1:         lang1,
		language2:         lang2,
		quizMode: quizMode,
	}
}

func (psc *profileSettingsConfig) toggleEnableTTS() {
	psc.enableTTS = !psc.enableTTS
}

func (psc *profileSettingsConfig) toggleSaveCancelSelection() {
	psc.saveCancelSelection = (psc.saveCancelSelection + 1) % 2
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

func (psc *profileSettingsConfig) cycleQuizMode(direction string) {
	curIdx := 0
	
	for i, qm := range orderedQuizMode {
		if qm == psc.quizMode {
			curIdx = i 
			break // found the 
		}
	} // found index where we currently are
	if direction == "left" {
		curIdx = (curIdx - 1 + len(orderedQuizMode)) % len(orderedQuizMode)
	} else {
		curIdx = (curIdx + 1) % len(orderedQuizMode)
	}
	psc.quizMode = orderedQuizMode[curIdx]
}


func (psc *profileSettingsConfig) down() {
	psc.selectedSetting = (psc.selectedSetting + 1) % totalSettings // 7 settings
}

func (psc *profileSettingsConfig) up() {
	psc.selectedSetting = (psc.selectedSetting - 1 + totalSettings) % totalSettings
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
		"Quiz Settings",
		m.selectedProfile.Name,
		m.selectedProfile.EnableSpeech,
		m.selectedProfile.DefaultNumQuestions,
		m.selectedProfile.DefaultNumAnswers,
		m.selectedProfile.Lang1,
		m.selectedProfile.Lang2,
		quizMode(m.selectedProfile.QuizMode),
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

	label = "Quiz Mode: "
	if psc.selectedSetting == settingQuizMode {
		b.WriteString(selectedStyle.Render(" > " + label))
	} else {
		b.WriteString("   " + label)
	}

	b.WriteString(settingSelectionStyle.Render(quizModeToString[psc.quizMode]))
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
		helpStyle.Render("\n↑/↓/⇆ to navigate • enter to select • tab to cycle • s/c to save/cancel • esc to go back")
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
	case "left", "h":
		if m.profileSettings.selectedSetting == settingNumQuestions {
			m.profileSettings.decreaseQuestions()
		} else if m.profileSettings.selectedSetting == settingTTS {
			m.profileSettings.toggleEnableTTS()
		} else if m.profileSettings.selectedSetting == settingLanguage1 {
			m.profileSettings.cycleLanguage("lang1")
		} else if m.profileSettings.selectedSetting == settingLanguage2 {
			m.profileSettings.cycleLanguage("lang2")
		} else if m.profileSettings.selectedSetting == settingQuizMode {
			m.profileSettings.cycleQuizMode("left")
		} else if m.profileSettings.selectedSetting == settingNumAnswerOptions {
			m.profileSettings.decreaseNumAnswers()
		} else if m.profileSettings.selectedSetting == settingSaveCancel {
			m.profileSettings.saveCancelSelection = saveSelected
		}
	case "right", "l":
		if m.profileSettings.selectedSetting == settingNumQuestions {
			m.profileSettings.increaseQuestions()
		} else if m.profileSettings.selectedSetting == settingTTS {
			m.profileSettings.toggleEnableTTS()
		} else if m.profileSettings.selectedSetting == settingLanguage1 {
			m.profileSettings.cycleLanguage("lang1")
		} else if m.profileSettings.selectedSetting == settingLanguage2 {
			m.profileSettings.cycleLanguage("lang2")
		} else if m.profileSettings.selectedSetting == settingQuizMode {
			m.profileSettings.cycleQuizMode("right")
		} else if m.profileSettings.selectedSetting == settingNumAnswerOptions {
			m.profileSettings.increaseNumAnswers()
		} else if m.profileSettings.selectedSetting == settingSaveCancel {
			m.profileSettings.saveCancelSelection = cancelSelected
		}
	case "tab":
		if m.profileSettings.selectedSetting == settingTTS {
			m.profileSettings.toggleEnableTTS()
		} else if m.profileSettings.selectedSetting == settingLanguage1 {
			m.profileSettings.cycleLanguage("lang1")
		} else if m.profileSettings.selectedSetting == settingLanguage2 {
			m.profileSettings.cycleLanguage("lang2")
		} else if m.profileSettings.selectedSetting == settingSaveCancel {
			m.profileSettings.toggleSaveCancelSelection()
		}
	case "enter":
		if m.profileSettings.selectedSetting == settingSaveCancel {
			if m.profileSettings.saveCancelSelection == saveSelected {
				saveSettings(m)
			} else {
				m.screen = screenChooseQuizMode
			}
		}
	case "s":
		saveSettings(m)
		return m, nil
	case "esc", "c": // 'c' for cancel
		m.updateProfileSuccess = ""
		m.buildQuizModeMenu()
		m.screen = screenChooseQuizMode
	}
	return m, nil
}

func saveSettings(m *Model) error {
	updatedProfile := store.Profile{
			Name: m.selectedProfile.Name, EnableSpeech: m.profileSettings.enableTTS,
			DefaultNumQuestions: m.profileSettings.defaultNumQuestions,
			DefaultNumAnswers:   m.profileSettings.defaultNumAnswers,
			Lang1:     m.profileSettings.language1,
			Lang2:     m.profileSettings.language2,
			QuizMode:  string(m.profileSettings.quizMode),
		}
	err := m.store.UpdateProfile(context.Background(), updatedProfile)
	if err != nil {
		fmt.Errorf("store: check some_new_column: %w", err)
		return err
	}
	// update the proifle on the model
	m.selectedProfile = updatedProfile
	m.updateProfileSuccess = "Profile settings updated successfully"
	m.screen = screenChooseQuizMode

	return nil
}