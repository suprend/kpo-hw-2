package menus

import (
	"github.com/charmbracelet/bubbles/textinput"

	"kpo-hw-2/internal/tui"
	"kpo-hw-2/internal/tui/styles"
)

// InputConfig describes configuration for creating input items.
type InputConfig struct {
	Placeholder string
	Prompt      string
	Initial     string
	Validate    func(string) error
	OnChange    func(string)
}

// SelectConfig describes configuration for select items.
type SelectConfig struct {
	InitialIndex int
	Validate     func(string) error
	OnChange     func(string)
}

// NewInputItem constructs an input menu entry with ready-to-use text input model.
func NewInputItem(key, title, description string, cfg InputConfig) Item {
	model := textinput.New()
	if cfg.Prompt != "" {
		model.Prompt = cfg.Prompt
	}
	if cfg.Placeholder != "" {
		model.Placeholder = cfg.Placeholder
		model.Width = len([]rune(cfg.Placeholder))
	}
	if cfg.Initial != "" {
		model.SetValue(cfg.Initial)
	}
	model.PromptStyle = styles.InputPromptStyle
	model.PlaceholderStyle = styles.InputPlaceholderStyle
	model.TextStyle = styles.InputTextStyle
	model.Cursor.Style = styles.InputCursorStyle

	return Item{
		Key:         key,
		Title:       title,
		Description: description,
		Kind:        ItemInput,
		Input: &InputField{
			Model:    model,
			Validate: cfg.Validate,
			OnChange: cfg.OnChange,
		},
	}
}

// NewSelectItem constructs a select menu entry with predefined options.
func NewSelectItem(key, title, description string, options []SelectOption, cfg SelectConfig) Item {
	field := NewSelectField(options, cfg.InitialIndex, cfg.Validate, cfg.OnChange)

	return Item{
		Key:         key,
		Title:       title,
		Description: description,
		Kind:        ItemSelect,
		Select:      field,
	}
}

// NewActionItem constructs a clickable menu entry with provided handler.
func NewActionItem(
	key, title, description string,
	action func(ctx tui.ScreenContext, values Values) tui.Result,
) Item {
	return Item{
		Key:         key,
		Title:       title,
		Description: description,
		Kind:        ItemAction,
		Action:      action,
	}
}

// NewPopItem returns action item that performs back navigation.
func NewPopItem(title, description string) Item {
	return NewActionItem("back", title, description, func(tui.ScreenContext, Values) tui.Result {
		return tui.Result{Pop: true}
	})
}
