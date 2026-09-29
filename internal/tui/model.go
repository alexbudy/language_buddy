package tui

import (
	"github.com/alexbudy/go_spanish_rewrite/internal/store"
	"github.com/charmbracelet/bubbles/textinput"
)

type screen int

const (
	screenProfileSelect screen = iota
	screenNewProfile
	screenGoodbye
)

// Model is the TUI model driving the application
type Model struct {
	store *store.Store

	screen screen

	profileMenu choiceList

	newProfileInput textinput.Model
	newProfileErr string

	existingProfiles []store.Profile // all active profiles
}

