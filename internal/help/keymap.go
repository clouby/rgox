package help

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Regex key.Binding
	Test  key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Help,
		k.Regex,
		k.Test,
		k.Quit,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Regex, k.Test, k.Up, k.Down},
		{k.Help, k.Quit},
	}
}

var InternalKeys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Regex: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "go to regex input"),
	),
	Test: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "go to test input"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "show help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc/ctrl+c", "quit"),
	),
}
