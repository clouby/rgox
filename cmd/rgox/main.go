package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/clouby/rgox/internal/regexinput"
)

const gap = "\n\n"

type (
	errMsg error
)

type Model struct {
	tea.Model
	textInput textinput.Model
	err       error
}

func initModel() Model {
	ti := regexinput.New()

	return Model{
		textInput: ti,
		err:       nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textInput.Width = msg.Width - 6
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"R%sX %s\n%s\n\n%s",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#8B7EC8")).Bold(true).Render("GO"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("- try your expressions here"),
		regexinput.WrapperInputStyle.Render(m.textInput.View()),
		"(esc to quit)",
	) + "\n"
}

func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
