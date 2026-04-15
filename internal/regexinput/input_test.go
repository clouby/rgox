package regexinput

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	input := New()

	if input.Placeholder != "insert your regular expression here" {
		t.Errorf("expected placeholder 'insert your regular expression here', got %q", input.Placeholder)
	}

	if input.Prompt != "┃ " {
		t.Errorf("expected prompt '┃ ', got %q", input.Prompt)
	}

	if input.CharLimit != 200 {
		t.Errorf("expected CharLimit 200, got %d", input.CharLimit)
	}

	if input.Width != 30 {
		t.Errorf("expected width 30, got %d", input.Width)
	}
}

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty string", "", true},
		{"single character", "a", false},
		{"valid regex", "[a-z]+", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateInputIntegration(t *testing.T) {
	input := New()

	_, cmd := input.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no command on key enter")
	}
}

func TestInputStyles(t *testing.T) {
	_ = InputStyle.Inline(true)
	_ = PromptStyle.Inline(true)
	_ = WrapperInputStyle.Inline(true)
}
