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
	screenGoodbye
)

// Model is the TUI model driving the application
type Model struct {
	store *store.Store

	screen screen

	profileMenu choiceList

	newProfileInput textinput.Model
	newProfileErr string

	profileToDelete string
	delProfileErr string // error for when choosing wrong entry for deletion
	invalidDeletePhraseErr string // error for when user typed invalid deletion phrase
	deleteProfileConfirmMenu choiceList
	specialDeletePhraseInput   textinput.Model


	existingProfiles []store.Profile // all active profiles
}

