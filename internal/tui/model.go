package tui

import (
	"github.com/alexbudy/go_spanish_rewrite/internal/store"
	"github.com/charmbracelet/bubbles/textinput"
)

type screen int

const (
	screenProfileSelect screen = iota
	screenNewProfile
	screenDeleteProfileConfirm
	screenChooseQuizMode
	screenProfileSettings
	screenManageWords
	screenWordDetails

	screenQuestion
	screenResults // results screen after completing a quiz

	screenGoodbye
)

// Model is the TUI model driving the application
type Model struct {
	store *store.Store

	screen screen

	profileMenu choiceList

	newProfileInput textinput.Model // input for creating new profile
	newProfileErr string

	profileToDelete string
	invalidDeletePhraseErr string // error for when user typed invalid deletion phrase
	deleteProfileConfirmMenu choiceList
	specialDeletePhraseInput   textinput.Model
	profileActionErr string // gerneric error for profile action

	profileRenameInput textinput.Model
	profileRenaming bool // are we in renaming mode
	profileRenameIndex int // index of item being renamed, or -1

	quizModeMenu choiceList
	quiz quizState
	answerMenu choiceList


	profileSettings profileSettingsConfig
	updateProfileSuccess string

	manageWordsMenu manageWordsList
	wordDetailsMenu wordDetails

	// profile states here
	existingProfiles []store.Profile // all active profiles
	selectedProfile store.Profile // Profile selected from the ProfileSelect screen
	allWords []store.Word // all words for the selected profile
}

