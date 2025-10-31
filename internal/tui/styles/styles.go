package styles

import "github.com/charmbracelet/lipgloss"

var (
	// TitleStyle used for screen titles or section headers.
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("213")).
			Bold(true)

	// IntroStyle used for introductory text.
	IntroStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	cursorSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("213")).
				Bold(true)

	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	itemSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("230")).
				Background(lipgloss.Color("57")).
				Bold(true)

	descriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("244"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203"))

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("231"))

	inputSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("231")).
				Background(lipgloss.Color("60"))

	selectValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("231"))

	selectValueSelectedStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("231")).
					Background(lipgloss.Color("60"))

	selectOptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	selectOptionActiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("231")).
				Background(lipgloss.Color("60"))

	// InputPromptStyle applied to textinput prompt.
	InputPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("111")).
				Bold(true)

	// InputPlaceholderStyle applied to placeholder text.
	InputPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))

	// InputTextStyle applied to textinput content.
	InputTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("231"))

	// InputCursorStyle applied to textinput cursor.
	InputCursorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("213")).
				Bold(true)
)

// CursorPrefix renders leading marker depending on selection state.
func CursorPrefix(selected bool) string {
	if selected {
		return cursorSelectedStyle.Render("> ")
	}
	return cursorStyle.Render("  ")
}

// ItemTitle renders item title based on selection state.
func ItemTitle(title string, selected bool) string {
	if selected {
		return itemSelectedStyle.Render(title)
	}
	return itemStyle.Render(title)
}

// Description renders helper or hint text.
func Description(text string) string {
	return descriptionStyle.Render(text)
}

// Intro renders introductory text block.
func Intro(text string) string {
	return IntroStyle.Render(text)
}

// InputView renders text input value with optional highlight.
func InputView(value string, selected bool) string {
	if selected {
		return inputSelectedStyle.Render(value)
	}
	return inputStyle.Render(value)
}

// SelectValue renders collapsed select current value.
func SelectValue(label string, selected bool) string {
	if selected {
		return selectValueSelectedStyle.Render(label)
	}
	return selectValueStyle.Render(label)
}

// SelectOption renders option row inside expanded select.
func SelectOption(label string, active bool) string {
	if active {
		return selectOptionActiveStyle.Render(label)
	}
	return selectOptionStyle.Render(label)
}

// Error renders validation message.
func Error(text string) string {
	return errorStyle.Render(text)
}
