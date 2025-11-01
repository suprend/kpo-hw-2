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
	Width       int
	OnChange    func(string)
}

const defaultInputWidth = 20

// SelectConfig describes configuration for select items.
type SelectConfig struct {
	InitialIndex int
	OnChange     func(string)
}

// NewInputItem constructs an input menu entry with ready-to-use text input model.
func NewInputItem(key, title, description string, cfg InputConfig) MenuItem {
	model := textinput.New()
	if cfg.Prompt != "" {
		model.Prompt = cfg.Prompt
	}
	width := cfg.Width
	if width < defaultInputWidth {
		width = defaultInputWidth
	}
	if cfg.Placeholder != "" {
		model.Placeholder = cfg.Placeholder
		if phLen := len([]rune(cfg.Placeholder)); phLen > width {
			width = phLen
		}
	}
	if cfg.Initial != "" {
		model.SetValue(cfg.Initial)
		if initialLen := len([]rune(cfg.Initial)); initialLen > width {
			width = initialLen
		}
	}
	model.Width = width
	model.CharLimit = 0
	model.PromptStyle = styles.InputPromptStyle
	model.PlaceholderStyle = styles.InputPlaceholderStyle
	model.TextStyle = styles.InputTextStyle
	model.Cursor.Style = styles.InputCursorStyle

	return &inputItem{
		key:         key,
		title:       title,
		description: description,
		model:       model,
		onChange:    cfg.OnChange,
	}
}

// NewSelectItem constructs a select menu entry with predefined options.
func NewSelectItem(key, title, description string, options []SelectOption, cfg SelectConfig) MenuItem {
	item := &selectItem{
		key:         key,
		title:       title,
		description: description,
		options:     options,
		onChange:    cfg.OnChange,
	}

	if len(options) == 0 {
		item.index = 0
		item.cursor = 0
	} else {
		idx := cfg.InitialIndex
		if idx < 0 || idx >= len(options) {
			idx = 0
		}
		item.index = idx
		item.cursor = idx
	}
	return item
}

// NewActionItem constructs a clickable menu entry with provided handler.
func NewActionItem(
	key, title, description string,
	action func(ctx tui.ScreenContext, values Values) tui.Result,
) MenuItem {
	return &actionItem{
		key:         key,
		title:       title,
		description: description,
		handler:     action,
	}
}

// NewPopItem returns action item that performs back navigation.
func NewPopItem(title, description string) MenuItem {
	return NewActionItem("back", title, description, func(tui.ScreenContext, Values) tui.Result {
		return tui.Result{Pop: true}
	})
}
