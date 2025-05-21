package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ThinkingMsg is sent every 500ms to update the thinking animation
type ThinkingMsg struct{}

// thinkingTick returns a command that sends a ThinkingMsg every 500ms
func thinkingTick() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return ThinkingMsg{}
	})
}

// StartThinking starts the thinking animation by setting the IsThinking flag to true,
// resetting the ThinkingDots counter, and adding a "thinking" system message.
// This function is typically used to indicate that the system is processing or waiting.
// The "thinking" system message added by this function can be removed by calling StopThinking.
func (m *Model) StartThinking() {
	m.IsThinking = true
	m.ThinkingDots = 0
	m.AddSystemMessage("thinking")
}

// StopThinking stops the thinking animation and removes the thinking message
func (m *Model) StopThinking() {
	m.IsThinking = false
	// Remove the thinking message if it exists
	if len(m.Messages) > 0 && m.Messages[len(m.Messages)-1].Content == "thinking" {
		m.Messages = m.Messages[:len(m.Messages)-1]
	}
}

// UpdateThinking handles the thinking animation update
func (m *Model) UpdateThinking() (tea.Model, tea.Cmd) {
	if m.IsThinking {
		m.ThinkingDots = (m.ThinkingDots + 1) % 4
		return m, thinkingTick()
	}
	return m, nil
}
