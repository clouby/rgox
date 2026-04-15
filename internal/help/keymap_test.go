package help

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/key"
)

func TestInternalKeysDefinition(t *testing.T) {
	if len(InternalKeys.Up.Keys()) == 0 {
		t.Error("Up binding should be defined")
	}

	if len(InternalKeys.Down.Keys()) == 0 {
		t.Error("Down binding should be defined")
	}

	if len(InternalKeys.Regex.Keys()) == 0 {
		t.Error("Regex binding should be defined")
	}

	if len(InternalKeys.Test.Keys()) == 0 {
		t.Error("Test binding should be defined")
	}

	if len(InternalKeys.Help.Keys()) == 0 {
		t.Error("Help binding should be defined")
	}

	if len(InternalKeys.Quit.Keys()) == 0 {
		t.Error("Quit binding should be defined")
	}
}

func TestKeyMapShortHelp(t *testing.T) {
	shortHelp := InternalKeys.ShortHelp()

	if len(shortHelp) != 3 {
		t.Errorf("expected 3 bindings in ShortHelp, got %d", len(shortHelp))
	}
}

func TestKeyMapFullHelp(t *testing.T) {
	fullHelp := InternalKeys.FullHelp()

	if len(fullHelp) != 2 {
		t.Errorf("expected 2 rows in FullHelp, got %d", len(fullHelp))
	}

	if len(fullHelp[0]) != 4 {
		t.Errorf("expected 4 bindings in first row, got %d", len(fullHelp[0]))
	}

	if len(fullHelp[1]) != 1 {
		t.Errorf("expected 1 binding in second row, got %d", len(fullHelp[1]))
	}
}

func TestKeyBindingsHaveKeys(t *testing.T) {
	tests := []struct {
		name    string
		binding key.Binding
		keys    []string
	}{
		{"Up", InternalKeys.Up, []string{"up", "↑"}},
		{"Down", InternalKeys.Down, []string{"down", "↓"}},
		{"Regex", InternalKeys.Regex, []string{"ctrl+r"}},
		{"Test", InternalKeys.Test, []string{"ctrl+t"}},
		{"Help", InternalKeys.Help, []string{"?"}},
		{"Quit", InternalKeys.Quit, []string{"esc", "ctrl+c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			help := tt.binding.Help()
			found := false
			for _, wantKey := range tt.keys {
				if strings.Contains(help.Key, wantKey) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected one of %v in binding %s, got %s", tt.keys, tt.name, help.Key)
			}
		})
	}
}

func TestKeyBindingsHaveHelp(t *testing.T) {
	bindings := []struct {
		name    string
		binding key.Binding
	}{
		{"Up", InternalKeys.Up},
		{"Down", InternalKeys.Down},
		{"Regex", InternalKeys.Regex},
		{"Test", InternalKeys.Test},
		{"Help", InternalKeys.Help},
		{"Quit", InternalKeys.Quit},
	}

	for _, tt := range bindings {
		t.Run(tt.name, func(t *testing.T) {
			help := tt.binding.Help()
			if help.Key == "" {
				t.Errorf("expected non-empty help key for %s", tt.name)
			}
			if help.Desc == "" {
				t.Errorf("expected non-empty help description for %s", tt.name)
			}
		})
	}
}
