package menus

import (
	tea "github.com/charmbracelet/bubbletea"

	"kpo-hw-2/internal/tui"
)

type ItemKind int

const (
	ItemAction ItemKind = iota
	ItemInput
	ItemSelect
)

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

type Values map[string]string

type SelectOption struct {
	Label string
	Value string
}
