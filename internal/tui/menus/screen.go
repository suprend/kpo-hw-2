package menus

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/styles"
)

// Screen renders menu composed of items (actions or inputs).
type Screen struct {
	title    string
	intro    string
	items    []MenuItem
	cursor   int
	values   Values
	emptyMsg string
}

// NewScreen constructs new menu screen instance.
func NewScreen(title, intro string, items []MenuItem) *Screen {
	values := make(Values)
	for _, item := range items {
		values[item.Key()] = item.Value()
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

	return s.items[s.cursor].Focus()
}

// Update implements tui.Screen.
func (s *Screen) Update(msg tea.Msg, ctx tui.ScreenContext) tui.Result {
	if len(s.items) == 0 {
		return tui.Result{}
	}

	current := s.items[s.cursor]

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up":
			if res, handled := current.Handle(msg, ctx, s.values); handled {
				s.values[current.Key()] = current.Value()
				return res
			}
			return s.moveCursor(-1)
		case "down":
			if res, handled := current.Handle(msg, ctx, s.values); handled {
				s.values[current.Key()] = current.Value()
				return res
			}
			return s.moveCursor(1)
		case "esc":
			if res, handled := current.Handle(msg, ctx, s.values); handled {
				s.values[current.Key()] = current.Value()
				return res
			}
			return tui.Result{Pop: true}
		}
	}

	result, handled := current.Handle(msg, ctx, s.values)
	if !handled {
		return tui.Result{}
	}
	s.values[current.Key()] = current.Value()
	return result
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

		view := item.View(selected)
		if view != "" {
			lines := strings.Split(view, "\n")
			for idx, line := range lines {
				if idx == 0 {
					b.WriteString(cursor + line + "\n")
				} else if line != "" {
					b.WriteString("  " + line + "\n")
				} else {
					b.WriteString("\n")
				}
			}
		} else {
			b.WriteString(cursor + "\n")
		}
	}

	return b.String()
}

var _ tui.Screen = (*Screen)(nil)

func (s *Screen) moveCursor(delta int) tui.Result {
	count := len(s.items)
	if count == 0 {
		return tui.Result{}
	}

	prev := s.items[s.cursor]
	prev.Blur()

	s.cursor = (s.cursor + delta + count) % count
	next := s.items[s.cursor]
	cmd := next.Focus()

	if cmd != nil {
		return tui.Result{Cmd: cmd}
	}

	return tui.Result{}
}

// SetFieldError updates validation error for input identified by key.
func (s *Screen) SetFieldError(key, message string) {
	item := s.findItem(key)
	if item == nil {
		return
	}
	if message == "" {
		item.ClearError()
		return
	}
	item.SetError(message)
}

// SetValue updates underlying value of input identified by key.
func (s *Screen) SetValue(key, value string) {
	item := s.findItem(key)
	if item == nil {
		return
	}
	item.SetValue(value)
	s.values[key] = item.Value()
}

// Value returns latest stored value for given key.
func (s *Screen) Value(key string) string {
	return s.values[key]
}

func (s *Screen) findItem(key string) MenuItem {
	for _, item := range s.items {
		if item.Key() == key {
			return item
		}
	}
	return nil
}
