package tui

import tea "github.com/charmbracelet/bubbletea"

// Screen represents a single interactive TUI view.
type Screen interface {
	// Name returns human-readable title for the screen (used in headers, logs).
	Name() string
	// Init runs when screen becomes active. Returned command is executed by Bubble Tea.
	Init(ctx ScreenContext) tea.Cmd
	// Update handles incoming Bubble Tea messages and may request navigation changes.
	Update(msg tea.Msg, ctx ScreenContext) Result
	// View renders the visual representation of the screen.
	View() string
}

// Result describes navigation or command actions requested by screen Update.
type Result struct {
	Push    Screen // push a new screen on top of stack
	Replace Screen // replace current screen
	Pop     bool   // pop current screen
	Cmd     tea.Cmd
}
