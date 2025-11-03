package menus

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/styles"
)

type actionItem struct {
	key         string
	title       string
	description string
	handler     func(tui.ScreenContext, Values) tui.Result
}

func (a *actionItem) Key() string         { return a.key }
func (a *actionItem) Title() string       { return a.title }
func (a *actionItem) Description() string { return a.description }
func (a *actionItem) Kind() ItemKind      { return ItemAction }

func (a *actionItem) Value() string   { return "" }
func (a *actionItem) SetValue(string) {}
func (a *actionItem) Focus() tea.Cmd  { return nil }
func (a *actionItem) Blur()           {}
func (a *actionItem) SetError(string) {}
func (a *actionItem) ClearError()     {}

func (a *actionItem) Handle(msg tea.Msg, ctx tui.ScreenContext, values Values) (tui.Result, bool) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return tui.Result{}, false
	}

	if keyMsg.String() != "enter" {
		return tui.Result{}, false
	}

	if a.handler == nil {
		return tui.Result{}, true
	}

	return a.handler(ctx, values), true
}

func (a *actionItem) View(selected bool) string {
	label := styles.ItemTitle(a.title, selected)

	var b strings.Builder
	b.WriteString(label)

	if a.description != "" {
		b.WriteString("\n")
		desc := styles.Description(a.description)
		b.WriteString("    ")
		b.WriteString(desc)
	}

	return b.String()
}
