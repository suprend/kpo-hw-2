package menus

import (
	tea "github.com/charmbracelet/bubbletea"

	"kpo-hw-2/internal/tui"
)

// ItemKind describes how a menu item behaves.
type ItemKind int

const (
	// ItemAction represents clickable command in menu.
	ItemAction ItemKind = iota
	// ItemInput represents editable input field.
	ItemInput
	// ItemSelect represents option selector with predefined values.
	ItemSelect
)

// MenuItem defines behaviour shared between different menu entries.
type MenuItem interface {
	Key() string
	Title() string
	Description() string
	Kind() ItemKind

	Value() string
	SetValue(string)

	Focus() tea.Cmd
	Blur()

	Handle(msg tea.Msg, ctx tui.ScreenContext, values Values) (tui.Result, bool)
	View(selected bool) string

	SetError(string)
	ClearError()
}

// Values keeps current values of menu inputs.
type Values map[string]string

// SelectOption describes single selectable value.
type SelectOption struct {
	Label string
	Value string
}
