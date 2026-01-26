package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/clouby/rgox/internal/caseinput"
	ihelp "github.com/clouby/rgox/internal/help"
	"github.com/clouby/rgox/internal/regexinput"
)

const (
	generalOffset = 6
)

type (
	errMsg error
)

type Model struct {
	tea.Model
	textInput      textinput.Model
	textareaInput  textarea.Model
	err            error
	keys           ihelp.KeyMap
	help           help.Model
	content        string
	containerStyle lipgloss.Style
}

func initModel() Model {
	return Model{
		textInput:     regexinput.New(),
		textareaInput: caseinput.New(),
		err:           errors.New(""),
		keys:          ihelp.InternalKeys,
		help:          help.New(),
		content:       "",
		containerStyle: lipgloss.NewStyle().
			MaxWidth(35).
			Height(5).
			MarginLeft(1).
			Padding(1).
			BorderStyle(
				lipgloss.HiddenBorder(),
			).
			BorderForeground(
				lipgloss.Color("63"),
			),
	}
}

func (m *Model) validateExpression() {
	const (
		Yellow  = "\033[1;33m"
		NoColor = "\033[0m"
		Red     = lipgloss.Color("#ffde9829")
	)

	text := m.textareaInput.Value()

	pattern, err := regexp.Compile(m.textInput.Value())
	if err != nil {
		m.textInput.Err = err
	}

	m.content = pattern.ReplaceAllString(text, lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(Red).Render("$0"))
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.textInput.Cursor.BlinkCmd(), m.textareaInput.Cursor.BlinkCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)
	m.textareaInput, cmd = m.textareaInput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textInput.Width = msg.Width - generalOffset
		m.textareaInput.SetWidth(msg.Width - generalOffset)
		m.help.Width = msg.Width - generalOffset
		m.containerStyle.Width(msg.Width)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Regex):
			if m.textareaInput.Focused() {
				m.textareaInput.Blur()
			}
			m.textInput.PromptStyle = regexinput.PromptStyle
			cmd = m.textInput.Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keys.Test):
			if m.textInput.Focused() {
				m.textInput.PromptStyle = lipgloss.NewStyle()
				m.textInput.Blur()
			}
			m.textareaInput.FocusedStyle.Prompt = regexinput.PromptStyle
			cmd = m.textareaInput.Focus()
			cmds = append(cmds, cmd)
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.validateExpression()

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	helpView := m.help.View(m.keys)
	if m.textInput.Err != nil {
		m.textInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	} else {
		m.textInput.PromptStyle = regexinput.PromptStyle
	}

	return fmt.Sprintf(
		lipgloss.NewStyle().Padding(1).Render("R%sX %s\n%s\n%s\n%s\n\n%s"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#8B7EC8")).Bold(true).Render("GO"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("- try your expressions..."),
		regexinput.WrapperInputStyle.Render(m.textInput.View()),
		regexinput.WrapperInputStyle.Render(m.textareaInput.View()),
		m.containerStyle.Render(m.content),
		helpView,
	) + "\n"
}

func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
