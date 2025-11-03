package tui

import tea "github.com/charmbracelet/bubbletea"

type Screen interface {
	Name() string
	Init(ctx ScreenContext) tea.Cmd
	Update(msg tea.Msg, ctx ScreenContext) Result
	View() string
}

type Result struct {
	Push    Screen
	Replace Screen
	Pop     bool
	Cmd     tea.Cmd
}
