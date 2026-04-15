package main

import (
	"bytes"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/muesli/termenv"
)

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
}

func TestInitialRender(t *testing.T) {
	tm := teatest.NewTestModel(
		t,
		initModel(),
		teatest.WithInitialTermSize(100, 40),
	)

	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("RGOX"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second),
	)

	tm.Send(tea.Quit())

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestRegexMatch(t *testing.T) {
	tm := teatest.NewTestModel(
		t,
		initModel(),
		teatest.WithInitialTermSize(100, 40),
	)

	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("RGOX"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second),
	)

	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestFinalModelContent(t *testing.T) {
	tm := teatest.NewTestModel(
		t,
		initModel(),
		teatest.WithInitialTermSize(100, 40),
	)

	tm.Send(tea.Quit())

	fm := tm.FinalModel(t)
	m, ok := fm.(Model)
	if !ok {
		t.Fatalf("final model has wrong type: %T", fm)
	}

	if m.content != "" {
		t.Errorf("expected empty content, got: %s", m.content)
	}

	if m.groupsContent != "" {
		t.Errorf("expected empty groupsContent, got: %s", m.groupsContent)
	}
}
