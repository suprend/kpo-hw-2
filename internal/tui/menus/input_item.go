package menus

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/styles"
)

type inputItem struct {
	key         string
	title       string
	description string
	errorText   string

	model    textinput.Model
	onChange func(string)
}

func (i *inputItem) Key() string         { return i.key }
func (i *inputItem) Title() string       { return i.title }
func (i *inputItem) Description() string { return i.description }
func (i *inputItem) Kind() ItemKind      { return ItemInput }

func (i *inputItem) Value() string { return i.model.Value() }

func (i *inputItem) SetValue(value string) {
	i.model.SetValue(value)
	i.errorText = ""
}

func (i *inputItem) Focus() tea.Cmd {
	return i.model.Focus()
}

func (i *inputItem) Blur() {
	i.model.Blur()
}

func (i *inputItem) Handle(msg tea.Msg, _ tui.ScreenContext, values Values) (tui.Result, bool) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "down":
			return tui.Result{}, false
		}
	}

	if msg == nil {
		return tui.Result{}, false
	}

	var cmd tea.Cmd
	i.model, cmd = i.model.Update(msg)

	value := i.model.Value()
	if values != nil {
		values[i.key] = value
	}

	if i.onChange != nil {
		i.onChange(value)
	}

	return tui.Result{Cmd: cmd}, true
}

func (i *inputItem) View(selected bool) string {
	var b strings.Builder

	label := styles.ItemTitle(i.title, selected)
	b.WriteString(label)
	b.WriteString(":\n")

	inputView := i.model.View()
	inputView = styles.InputView(inputView, selected)
	b.WriteString("  " + inputView + "\n")

	if i.errorText != "" {
		errLine := styles.Error("Ошибка: " + i.errorText)
		b.WriteString("  " + errLine + "\n")
	}

	if i.description != "" {
		desc := styles.Description(i.description)
		b.WriteString("  " + desc + "\n")
	}

	return b.String()
}

func (i *inputItem) SetError(message string) {
	i.errorText = message
}

func (i *inputItem) ClearError() {
	i.errorText = ""
}
