package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ihelp "github.com/clouby/rgox/internal/help"
	"github.com/clouby/rgox/internal/regexinput"
	"github.com/clouby/rgox/internal/testinput"
)

const (
	gap           = "\n\n"
	generalOffset = 6
)

type (
	errMsg error
)

type Model struct {
	tea.Model
	textInput     textinput.Model
	textareaInput textarea.Model
	err           error
	keys          ihelp.KeyMap
	help          help.Model
	cursor        int
	selected      any
}

func initModel() Model {
	ti := regexinput.New()
	ta := testinput.New()

	return Model{
		textInput:     ti,
		textareaInput: ta,
		err:           nil,
		keys:          ihelp.InternalKeys,
		help:          help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, textarea.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textInput.Width = msg.Width - generalOffset
		m.textareaInput.SetWidth(msg.Width - generalOffset)
		m.help.Width = msg.Width - generalOffset
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Regex):
			if m.textareaInput.Focused() {
				m.textareaInput.Blur()
			}
			m.textInput.PromptStyle = regexinput.PromptStyle
			return m, m.textInput.Focus()
		case key.Matches(msg, m.keys.Test):
			if m.textInput.Focused() {
				m.textInput.Blur()
				m.textInput.PromptStyle = lipgloss.NewStyle()
			}

			m.textareaInput.FocusedStyle.Prompt = regexinput.PromptStyle
			return m, m.textareaInput.Focus()
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	m.textareaInput, cmd = m.textareaInput.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	var keyMap help.KeyMap = m.keys
	helpView := m.help.View(keyMap)

	return fmt.Sprintf(
		"R%sX %s\n%s\n%s\n\n%s",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#8B7EC8")).Bold(true).Render("GO"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("- try your expressions..."),
		regexinput.WrapperInputStyle.Render(m.textInput.View()),
		regexinput.WrapperInputStyle.Render(m.textareaInput.View()),
		helpView,
	) + "\n"
}

func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
