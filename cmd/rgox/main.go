package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
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

var errorInputStyle = lipgloss.Color("9")

type Model struct {
	tea.Model
	textInput      textinput.Model
	textareaInput  textarea.Model
	err            error
	keys           ihelp.KeyMap
	help           help.Model
	content        string
	groupsContent  string
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
		groupsContent: "",
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

var groupColors = []string{
	"#9C27B0", "#673AB7", "#3F51B5", "#2196F3",
	"#00BCD4", "#009688", "#4CAF50", "#8BC34A", "#CDDC39",
	"#FFC107", "#FF9800", "#FF5722", "#795548", "#607D8B",
}

func (m *Model) validateExpression() {
	const (
		Green = lipgloss.Color("10")
		Black = lipgloss.Color("0")
	)

	text := m.textareaInput.Value()

	pattern, err := regexp.Compile(m.textInput.Value())
	if err != nil {
		m.textInput.Err = err
		return
	}

	m.content = pattern.ReplaceAllString(text, lipgloss.NewStyle().Foreground(Black).Background(Green).Render("$0"))

	matches := pattern.FindAllStringSubmatch(text, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		colors := make([]string, len(groupColors))
		copy(colors, groupColors)
		rand.Shuffle(len(colors), func(i, j int) { colors[i], colors[j] = colors[j], colors[i] })

		var matchGroups []string
		for matchIdx, match := range matches {
			if matchIdx >= 0 {
				matchGroups = append(matchGroups, lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginLeft(1).Render(fmt.Sprintf("Match %d", matchIdx+1)))
			}
			var parts []string
			for i, group := range match[1:] {
				color := colors[(matchIdx+i)%len(colors)]
				parts = append(parts, lipgloss.NewStyle().MarginLeft(1).Foreground(lipgloss.Color(color)).Render(fmt.Sprintf("$%d: %s ", i+1, group)))
			}
			if len(parts) > 0 {
				matchGroups = append(matchGroups, lipgloss.JoinHorizontal(0, parts...))
			}
		}
		if len(matchGroups) > 0 {
			m.groupsContent = "Groups:\n" + lipgloss.JoinVertical(0, matchGroups...)
		} else {
			m.groupsContent = ""
		}
	} else {
		m.groupsContent = ""
	}
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
		case key.Matches(msg, m.keys.Regex, m.keys.Up):
			if m.textareaInput.Focused() {
				m.textareaInput.Blur()
			}
			m.textInput.PromptStyle = regexinput.PromptStyle
			cmd = m.textInput.Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keys.Test, m.keys.Down):
			if m.textInput.Focused() {
				m.textInput.PromptStyle = lipgloss.NewStyle()
				m.textInput.Blur()
			}
			m.textareaInput.FocusedStyle.Prompt = regexinput.PromptStyle
			cmd = m.textareaInput.Focus()
			cmds = append(cmds, cmd)
		}
		m.validateExpression()
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	helpView := m.help.View(m.keys)
	errorInputView := ""
	if m.textInput.Err != nil {
		m.textInput.PromptStyle = lipgloss.NewStyle().Foreground(errorInputStyle)
		errorInputView = m.textInput.Err.Error()
	} else {
		m.textInput.PromptStyle = regexinput.PromptStyle
	}

	_ = errorInputView

	return fmt.Sprintf(
		lipgloss.NewStyle().Padding(1).Render("R%sX %s\n%s\n%s\n%s\n\n%s\n%s\n%s"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#8B7EC8")).Bold(true).Render("GO"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("- try your expressions..."),
		regexinput.WrapperInputStyle.Render(m.textInput.View()),
		lipgloss.NewStyle().Foreground(errorInputStyle).Render(errorInputView),
		regexinput.WrapperInputStyle.Render(m.textareaInput.View()),
		m.containerStyle.Render(m.content),
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(m.groupsContent),
		helpView,
	) + "\n"
}

func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
