// Package regexinput for get the values to be expressed.
package regexinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	InputStyle        = lipgloss.NewStyle().Padding(4)
	PromptStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#8B7EC8"))
	WrapperInputStyle = lipgloss.NewStyle().Padding(1, 1, 0)
)

func New() textinput.Model {
	textInput := textinput.New()
	textInput.Placeholder = "insert your regular expression here"
	textInput.Prompt = "┃ "
	textInput.CharLimit = 200 // Default Char Limit Base
	textInput.PromptStyle = PromptStyle
	textInput.Focus()
	textInput.Width = 30 // Default Width Base
	textInput.TextStyle = InputStyle

	return textInput
}
