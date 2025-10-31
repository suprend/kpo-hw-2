package menus

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/styles"
)

// Screen renders menu composed of items (actions or inputs).
type Screen struct {
	title    string
	intro    string
	items    []Item
	cursor   int
	values   Values
	emptyMsg string
}

// NewScreen constructs new menu screen instance.
func NewScreen(title, intro string, items []Item) *Screen {
	values := make(Values)
	for _, item := range items {
		if item.Kind == ItemInput && item.Input != nil {
			values[item.Key] = item.Input.Model.Value()
		}
		if item.Kind == ItemSelect && item.Select != nil {
			values[item.Key] = item.Select.Value()
		}
	}

	return &Screen{
		title:  title,
		intro:  intro,
		items:  items,
		values: values,
	}
}

// Name implements tui.Screen.
func (s *Screen) Name() string { return s.title }

// Init implements tui.Screen.
func (s *Screen) Init(tui.ScreenContext) tea.Cmd {
	if len(s.items) == 0 {
		return nil
	}

	current := s.items[s.cursor]
	if current.Kind == ItemInput && current.Input != nil {
		current.Input.Model.Focus()
		return textinput.Blink
	}

	return nil
}

// Update implements tui.Screen.
func (s *Screen) Update(msg tea.Msg, ctx tui.ScreenContext) tui.Result {
	if len(s.items) == 0 {
		return tui.Result{}
	}

	current := s.items[s.cursor]
	if handled := s.handleExpandedSelect(current, msg); handled != nil {
		return *handled
	}

	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "up":
			s.moveCursor(-1)
			return tui.Result{}
		case "down":
			s.moveCursor(1)
			return tui.Result{}
		case "enter":
			if current.Kind == ItemInput {
				return tui.Result{}
			}
			if current.Kind == ItemSelect && current.Select != nil {
				current.Select.Expand()
				s.values[current.Key] = current.Select.Value()
				return tui.Result{}
			}
			return s.activateCurrent(ctx)
		case "esc":
			return tui.Result{Pop: true}
		}
	}

	if current.Kind == ItemInput && current.Input != nil {
		return s.updateInput(current.Input, current.Key, msg)
	}

	return tui.Result{}
}

func (s *Screen) handleExpandedSelect(item Item, msg tea.Msg) *tui.Result {
	if item.Kind != ItemSelect || item.Select == nil || !item.Select.Expanded {
		return nil
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		res := tui.Result{}
		return &res
	}

	switch keyMsg.String() {
	case "up":
		item.Select.Move(-1)
		return &tui.Result{}
	case "down":
		item.Select.Move(1)
		return &tui.Result{}
	case "enter":
		item.Select.Commit()
		s.values[item.Key] = item.Select.Value()
		item.Select.Collapse()
		return &tui.Result{}
	case "esc":
		item.Select.Collapse()
		s.values[item.Key] = item.Select.Value()
		return &tui.Result{}
	default:
		return &tui.Result{}
	}
}

// WithEmptyMessage configures placeholder text shown when menu has no items.
func (s *Screen) WithEmptyMessage(message string) *Screen {
	s.emptyMsg = message
	return s
}

// View implements tui.Screen.
func (s *Screen) View() string {
	if len(s.items) == 0 {
		if s.emptyMsg != "" {
			return s.emptyMsg
		}
		return "Нет пунктов меню"
	}

	var b strings.Builder

	if s.intro != "" {
		b.WriteString(styles.Intro(s.intro))
		b.WriteString("\n\n")
	}

	for i, item := range s.items {
		selected := i == s.cursor
		cursor := styles.CursorPrefix(selected)

		switch item.Kind {
		case ItemInput:
			label := styles.ItemTitle(item.Title, selected)
			b.WriteString(cursor + label + ":\n")
			if item.Input != nil {
				inputView := item.Input.Model.View()
				inputView = styles.InputView(inputView, selected)
				b.WriteString("  " + inputView + "\n")
				if item.Input.ErrorText != "" {
					errLine := styles.Error("Ошибка: " + item.Input.ErrorText)
					b.WriteString("  " + errLine + "\n")
				}
			}
			if item.Description != "" {
				desc := styles.Description(item.Description)
				b.WriteString("  " + desc + "\n")
			}
		case ItemSelect:
			label := styles.ItemTitle(item.Title, selected)
			b.WriteString(cursor + label + ":\n")
			if item.Select != nil {
				if item.Select.Expanded {
					for idx, opt := range item.Select.Options {
						active := idx == item.Select.cursor
						optCursor := styles.CursorPrefix(active)
						line := "  " + optCursor + styles.SelectOption(opt.Label, active)
						b.WriteString(line + "\n")
					}
				} else {
					labelText := item.Select.Label()
					if labelText == "" {
						labelText = "Не выбрано"
					}
					valueView := styles.SelectValue(labelText, selected)
					b.WriteString("  " + valueView + "\n")
				}
				if item.Select.ErrorText != "" {
					errLine := styles.Error("Ошибка: " + item.Select.ErrorText)
					b.WriteString("  " + errLine + "\n")
				}
			}
			if item.Description != "" {
				desc := styles.Description(item.Description)
				b.WriteString("  " + desc + "\n")
			}
		default:
			label := styles.ItemTitle(item.Title, selected)
			b.WriteString(cursor + label + "\n")
			if item.Description != "" {
				desc := styles.Description(item.Description)
				b.WriteString("    " + desc + "\n")
			}
		}
	}

	return b.String()
}

var _ tui.Screen = (*Screen)(nil)

func (s *Screen) moveCursor(delta int) {
	count := len(s.items)
	if count == 0 {
		return
	}

	prev := s.items[s.cursor]
	if prev.Kind == ItemInput && prev.Input != nil {
		prev.Input.Model.Blur()
	}

	s.cursor = (s.cursor + delta + count) % count
	next := s.items[s.cursor]
	if next.Kind == ItemInput && next.Input != nil {
		next.Input.Model.Focus()
	}
}

func (s *Screen) activateCurrent(ctx tui.ScreenContext) tui.Result {
	current := s.items[s.cursor]
	switch current.Kind {
	case ItemAction:
		if current.Action == nil {
			return tui.Result{}
		}

		return current.Action(ctx, s.values)
	default:
		return tui.Result{}
	}
}

func (s *Screen) updateInput(field *InputField, key string, msg tea.Msg) tui.Result {
	if msg == nil {
		return tui.Result{}
	}

	var cmd tea.Cmd
	field.Model, cmd = field.Model.Update(msg)
	value := field.Model.Value()
	s.values[key] = value
	if field.OnChange != nil {
		field.OnChange(value)
	}

	return tui.Result{Cmd: cmd}
}

// SetFieldError updates validation error for input identified by key.
func (s *Screen) SetFieldError(key, message string) {
	for i := range s.items {
		item := &s.items[i]
		if item.Key == key && item.Kind == ItemInput && item.Input != nil {
			item.Input.ErrorText = message
			return
		}
		if item.Key == key && item.Kind == ItemSelect && item.Select != nil {
			item.Select.ErrorText = message
			return
		}
	}
}

// SetValue updates underlying value of input identified by key.
func (s *Screen) SetValue(key, value string) {
	for i := range s.items {
		item := &s.items[i]
		if item.Key == key && item.Kind == ItemInput && item.Input != nil {
			item.Input.Model.SetValue(value)
			item.Input.ErrorText = ""
			s.values[key] = value
			return
		}
		if item.Key == key && item.Kind == ItemSelect && item.Select != nil {
			item.Select.SetValue(value)
			item.Select.ErrorText = ""
			s.values[key] = item.Select.Value()
			return
		}
	}
}

// Value returns latest stored value for given key.
func (s *Screen) Value(key string) string {
	return s.values[key]
}
