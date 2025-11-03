package menus

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/styles"
)

type selectItem struct {
	key         string
	title       string
	description string
	errorText   string

	options  []SelectOption
	index    int
	cursor   int
	expanded bool

	onChange func(string)
}

func (s *selectItem) Key() string         { return s.key }
func (s *selectItem) Title() string       { return s.title }
func (s *selectItem) Description() string { return s.description }
func (s *selectItem) Kind() ItemKind      { return ItemSelect }

func (s *selectItem) Value() string {
	if s.index < 0 || s.index >= len(s.options) {
		return ""
	}
	return s.options[s.index].Value
}

func (s *selectItem) SetValue(value string) {
	if len(s.options) == 0 {
		s.index = 0
		s.cursor = 0
		return
	}

	match := -1
	for idx, opt := range s.options {
		if opt.Value == value {
			match = idx
			break
		}
	}

	if match == -1 {
		match = 0
	}

	s.index = match
	s.cursor = match
}

func (s *selectItem) Focus() tea.Cmd {
	return nil
}

func (s *selectItem) Blur() {
	s.expanded = false
	s.cursor = s.index
}

func (s *selectItem) Handle(msg tea.Msg, _ tui.ScreenContext, values Values) (tui.Result, bool) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return tui.Result{}, false
	}

	key := keyMsg.String()

	if !s.expanded {
		if key == "enter" {
			if len(s.options) == 0 {
				return tui.Result{}, true
			}
			s.expanded = true
			s.cursor = s.index
			return tui.Result{}, true
		}
		return tui.Result{}, false
	}

	switch key {
	case "up":
		s.moveCursor(-1)
		return tui.Result{}, true
	case "down":
		s.moveCursor(1)
		return tui.Result{}, true
	case "enter":
		s.commitSelection(values)
		return tui.Result{}, true
	case "esc":
		s.expanded = false
		s.cursor = s.index
		return tui.Result{}, true
	default:
		return tui.Result{}, true
	}
}

func (s *selectItem) View(selected bool) string {
	var b strings.Builder

	label := styles.ItemTitle(s.title, selected)
	b.WriteString(label)
	b.WriteString(":\n")

	if len(s.options) == 0 {
		placeholder := styles.Description("Нет вариантов")
		b.WriteString("  " + placeholder + "\n")
	} else if s.expanded {
		for idx, opt := range s.options {
			active := idx == s.cursor
			cursor := styles.CursorPrefix(active)
			line := cursor + styles.SelectOption(opt.Label, active)
			b.WriteString("  " + line + "\n")
		}
	} else {
		labelText := s.currentLabel()
		valueView := styles.SelectValue(labelText, selected)
		b.WriteString("  " + valueView + "\n")
	}

	if s.errorText != "" {
		errLine := styles.Error("Ошибка: " + s.errorText)
		b.WriteString("  " + errLine + "\n")
	}

	if s.description != "" {
		desc := styles.Description(s.description)
		b.WriteString("  " + desc + "\n")
	}

	return b.String()
}

func (s *selectItem) SetError(message string) {
	s.errorText = message
}

func (s *selectItem) ClearError() {
	s.errorText = ""
}

func (s *selectItem) moveCursor(delta int) {
	count := len(s.options)
	if count == 0 {
		s.cursor = 0
		return
	}

	s.cursor = (s.cursor + delta + count) % count
}

func (s *selectItem) commitSelection(values Values) {
	s.expanded = false
	if len(s.options) == 0 {
		return
	}

	if s.cursor < 0 || s.cursor >= len(s.options) {
		s.cursor = s.index
		return
	}

	if s.index != s.cursor {
		s.index = s.cursor
		if s.onChange != nil {
			s.onChange(s.Value())
		}
	}

	if values != nil {
		values[s.key] = s.Value()
	}
}

func (s *selectItem) currentLabel() string {
	if len(s.options) == 0 || s.index < 0 || s.index >= len(s.options) {
		return "Не выбрано"
	}
	return s.options[s.index].Label
}
