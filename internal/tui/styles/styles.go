package styles

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("213")).
			Bold(true)

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

	InputPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("111")).
				Bold(true)

	InputPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))

	InputTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("231"))

	InputCursorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("213")).
				Bold(true)
)

func CursorPrefix(selected bool) string {
	if selected {
		return cursorSelectedStyle.Render("> ")
	}
	return cursorStyle.Render("  ")
}

func ItemTitle(title string, selected bool) string {
	if selected {
		return itemSelectedStyle.Render(title)
	}
	return itemStyle.Render(title)
}

func Description(text string) string {
	return descriptionStyle.Render(text)
}

func Intro(text string) string {
	return IntroStyle.Render(text)
}

func InputView(value string, selected bool) string {
	if selected {
		return inputSelectedStyle.Render(value)
	}
	return inputStyle.Render(value)
}

func SelectValue(label string, selected bool) string {
	if selected {
		return selectValueSelectedStyle.Render(label)
	}
	return selectValueStyle.Render(label)
}

func SelectOption(label string, active bool) string {
	if active {
		return selectOptionActiveStyle.Render(label)
	}
	return selectOptionStyle.Render(label)
}

func Error(text string) string {
	return errorStyle.Render(text)
}
