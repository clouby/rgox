package testinput

import (
	"github.com/charmbracelet/bubbles/textarea"
)

func New() textarea.Model {
	textareaInput := textarea.New()
	textareaInput.Placeholder = "insert your test string here"
	textareaInput.Prompt = "┃ "
	textareaInput.ShowLineNumbers = false
	textareaInput.SetWidth(30)
	textareaInput.SetHeight(3)

	return textareaInput
}
