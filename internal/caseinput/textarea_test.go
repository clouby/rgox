package caseinput

import (
	"testing"
)

func TestNew(t *testing.T) {
	textarea := New()

	if textarea.Placeholder != "insert your test string here" {
		t.Errorf("expected placeholder 'insert your test string here', got %q", textarea.Placeholder)
	}

	if textarea.Prompt != "┃ " {
		t.Errorf("expected prompt '┃ ', got %q", textarea.Prompt)
	}

	if !textarea.ShowLineNumbers {
		t.Error("expected ShowLineNumbers to be true")
	}

	if textarea.Height() != 3 {
		t.Errorf("expected height 3, got %d", textarea.Height())
	}
}
