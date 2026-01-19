package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

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
		textInput:      regexinput.New(),
		textareaInput:  caseinput.New(),
		err:            nil,
		keys:           ihelp.InternalKeys,
		help:           help.New(),
		content:        "",
		containerStyle: lipgloss.NewStyle().Width(25).Height(5).MarginLeft(1).Padding(1).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")),
	}
}

func (m *Model) validateExpression() {
	if m.textInput.Value() == "" {
		m.err = errors.New("expression cannot be empty")
		return
	}
	re := regexp.MustCompile(m.textInput.Value())

	m.content = strings.Join(re.FindAllString(m.textareaInput.Value(), -1), "\n")
	// if !re.MatchString(m.textareaInput.Value()) {
	// 	m.err = errors.New("expression does not match")
	// 	return
	// }
	m.err = nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, textarea.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
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

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)
	m.textareaInput, cmd = m.textareaInput.Update(msg)
	cmds = append(cmds, cmd)

	m.content = m.textInput.Value()

	m.validateExpression()

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	helpView := m.help.View(m.keys)
	errView := ""
	if m.err != nil {
		errView = fmt.Sprintf("Error: %s", m.err.Error())
	}

	return fmt.Sprintf(
		lipgloss.NewStyle().Padding(1).Render("R%sX %s\n%s\n%s\n%s\n%s\n\n%s"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#8B7EC8")).Bold(true).Render("GO"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("- try your expressions..."),
		regexinput.WrapperInputStyle.Render(m.textInput.View()),
		lipgloss.NewStyle().Foreground(lipgloss.Color("9")).PaddingTop(1).PaddingLeft(1).Render(errView),
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
