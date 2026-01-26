// Package regexinput for get the values to be expressed.
package regexinput

import (
	"errors"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	InputStyle        = lipgloss.NewStyle().Padding(4)
	PromptStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#8B7EC8"))
	WrapperInputStyle = lipgloss.NewStyle().Padding(1, 1, 0)
)

func validateInput(s string) error {
	if len(s) == 0 {
		return errors.New("input cannot be empty")
	}
	return nil
}

func New() textinput.Model {
	textInput := textinput.New()
	textInput.Placeholder = "insert your regular expression here"
	textInput.Prompt = "┃ "
	textInput.CharLimit = 200
	textInput.PromptStyle = PromptStyle
	textInput.Focus()
	textInput.Width = 30
	textInput.TextStyle = InputStyle

	textInput.Validate = validateInput

	return textInput
}
