package caseinput

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

func New() textarea.Model {
	textareaInput := textarea.New()
	textareaInput.Placeholder = "insert your test string here"
	textareaInput.Prompt = "┃ "
	textareaInput.ShowLineNumbers = true
	textareaInput.FocusedStyle.CursorLine = lipgloss.NewStyle()
	textareaInput.SetWidth(30)
	textareaInput.SetHeight(3)

	return textareaInput
}
