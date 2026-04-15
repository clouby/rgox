package main

import (
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

var (
	errorInputStyle  = lipgloss.Color("9")
	brandColor       = lipgloss.Color("#8B7EC8")
	grayColor        = lipgloss.Color("240")
	greenColor       = lipgloss.Color("10")
	blackColor       = lipgloss.Color("0")
	matchStyle       = lipgloss.NewStyle().Foreground(blackColor).Background(greenColor)
	brandStyle       = lipgloss.NewStyle().Foreground(brandColor).Bold(true)
	dimStyle         = lipgloss.NewStyle().Foreground(grayColor)
	errorStyle       = lipgloss.NewStyle().Foreground(errorInputStyle)
	matchHeaderStyle = lipgloss.NewStyle().Foreground(grayColor).MarginLeft(1)
	groupPartStyle   = lipgloss.NewStyle().MarginLeft(1)
	containerPadding = lipgloss.NewStyle().Padding(1)
)

type Model struct {
	tea.Model
	textInput      textinput.Model
	textareaInput  textarea.Model
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
		keys:          ihelp.InternalKeys,
		help:          help.New(),
		content:       "",
		groupsContent: "",
		containerStyle: lipgloss.NewStyle().
			MaxWidth(35).
			Height(5).
			MarginLeft(1).
			Padding(1).
			BorderStyle(lipgloss.HiddenBorder()).
			BorderForeground(lipgloss.Color("63")),
	}
}

var groupColors = []string{
	"#9C27B0", "#673AB7", "#3F51B5", "#2196F3",
	"#00BCD4", "#009688", "#4CAF50", "#8BC34A", "#CDDC39",
	"#FFC107", "#FF9800", "#FF5722", "#795548", "#607D8B",
}

func (m *Model) validateExpression() {
	text := m.textareaInput.Value()

	pattern, err := regexp.Compile(m.textInput.Value())
	if err != nil {
		m.textInput.Err = err
		m.content = ""
		m.groupsContent = ""
		return
	}

	m.content = pattern.ReplaceAllString(text, matchStyle.Render("$0"))

	matches := pattern.FindAllStringSubmatch(text, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		colors := make([]string, len(groupColors))
		copy(colors, groupColors)
		rand.Shuffle(len(colors), func(i, j int) { colors[i], colors[j] = colors[j], colors[i] })

		var matchGroups []string
		for matchIdx, match := range matches {
			matchGroups = append(matchGroups, matchHeaderStyle.Render(fmt.Sprintf("Match %d", matchIdx+1)))
			var parts []string
			for i, group := range match[1:] {
				color := colors[(matchIdx+i)%len(colors)]
				parts = append(parts, groupPartStyle.Foreground(lipgloss.Color(color)).Render(fmt.Sprintf("$%d: %s ", i+1, group)))
			}
			if len(parts) > 0 {
				matchGroups = append(matchGroups, lipgloss.JoinHorizontal(0, parts...))
			}
		}
		m.groupsContent = "Groups:\n" + lipgloss.JoinVertical(0, matchGroups...)
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

	header := containerPadding.Render(fmt.Sprintf("R%sX %s", brandStyle.Render("GO"), dimStyle.Render("- try your expressions...")))

	return fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n\n%s\n%s\n",
		header,
		regexinput.WrapperInputStyle.Render(m.textInput.View()),
		errorStyle.Render(errorInputView),
		regexinput.WrapperInputStyle.Render(m.textareaInput.View()),
		m.containerStyle.Render(m.content),
		dimStyle.Render(m.groupsContent),
		helpView,
	)
}

func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
