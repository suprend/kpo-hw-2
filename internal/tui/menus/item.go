package menus

import (
	"github.com/charmbracelet/bubbles/textinput"

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

// Item defines single entry inside menu-driven screen.
type Item struct {
	// Key uniquely identifies item and used for collecting values of inputs.
	Key string
	// Title is short label displayed to a user.
	Title string
	// Description provides additional hint text under the item.
	Description string
	// Kind defines whether item is action or input.
	Kind ItemKind
	// Action handles activation of action item (e.g. pressing Enter).
	Action func(ctx tui.ScreenContext, values Values) tui.Result
	// Input holds configuration/state for editable item.
	Input *InputField
	// Select holds configuration/state for selectable item.
	Select *SelectField
}

// Values keeps current values of menu inputs.
type Values map[string]string

// InputField stores state of single text input within menu.
type InputField struct {
	Model     textinput.Model
	Validate  func(string) error
	ErrorText string
	OnChange  func(string)
}

// SelectOption describes single selectable value.
type SelectOption struct {
	Label string
	Value string
}

// SelectField stores state of select input within menu.
type SelectField struct {
	Options   []SelectOption
	Index     int
	cursor    int
	Expanded  bool
	Validate  func(string) error
	ErrorText string
	OnChange  func(string)
}

func (s *SelectField) optionAt(idx int) (SelectOption, bool) {
	if s == nil || len(s.Options) == 0 {
		return SelectOption{}, false
	}
	if idx < 0 || idx >= len(s.Options) {
		return SelectOption{}, false
	}
	return s.Options[idx], true
}

// Value returns currently selected option value or empty string when out of range.
func (s *SelectField) Value() string {
	if opt, ok := s.optionAt(s.Index); ok {
		return opt.Value
	}
	return ""
}

// Label returns human-readable label of selected option.
func (s *SelectField) Label() string {
	if opt, ok := s.optionAt(s.Index); ok {
		return opt.Label
	}
	return ""
}

// SetValue updates selected option by matching provided value.
func (s *SelectField) SetValue(value string) {
	if s == nil || len(s.Options) == 0 {
		return
	}
	if value == "" {
		s.setIndex(0)
		return
	}

	for i, opt := range s.Options {
		if opt.Value == value {
			s.setIndex(i)
			return
		}
	}

	s.setIndex(0)
}

// Expand prepares selector for navigating options.
func (s *SelectField) Expand() {
	if s == nil || len(s.Options) == 0 {
		return
	}
	s.Expanded = true
	s.cursor = s.Index
}

// Collapse exits option navigation while keeping current selection.
func (s *SelectField) Collapse() {
	if s == nil {
		return
	}
	s.Expanded = false
	s.cursor = s.Index
}

// Move shifts cursor across available options.
func (s *SelectField) Move(delta int) {
	if s == nil || len(s.Options) == 0 {
		return
	}

	count := len(s.Options)
	s.cursor = (s.cursor + delta + count) % count
}

// Commit saves cursor selection as active option and triggers callbacks.
func (s *SelectField) Commit() {
	if s == nil || len(s.Options) == 0 {
		return
	}

	if s.cursor < 0 || s.cursor >= len(s.Options) {
		return
	}

	if s.Index != s.cursor {
		s.Index = s.cursor
		if s.OnChange != nil {
			s.OnChange(s.Options[s.Index].Value)
		}
	}
}

func (s *SelectField) setIndex(index int) {
	if len(s.Options) == 0 {
		s.Index = 0
		s.cursor = 0
		return
	}
	if index < 0 || index >= len(s.Options) {
		index = 0
	}
	s.Index = index
	s.cursor = index
}

// NewSelectField constructs select field with normalized configuration.
func NewSelectField(options []SelectOption, initialIndex int, validate func(string) error, onChange func(string)) *SelectField {
	field := &SelectField{
		Options:  options,
		Validate: validate,
		OnChange: onChange,
	}
	field.setIndex(initialIndex)
	return field
}
