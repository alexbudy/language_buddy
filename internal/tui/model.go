package tui

import (
	"github.com/alexbudy/go_spanish_rewrite/internal/store"
)

type screen int

const (
	screenProfileSelect screen = iota

	screenGoodbye
)

// Model is the TUI model driving the application
type Model struct {
	store *store.Store

	screen screen

	profileMenu choiceList

	existingProfiles []store.Profile // all active profiles
}

