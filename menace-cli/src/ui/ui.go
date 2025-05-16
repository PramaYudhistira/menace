package ui

import (
	"menace-go/llmServer"
	"menace-go/model"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	model.Model
	agent  *llmServer.Agent
	Width  int
	Height int
	// Scroll offset (0 = bottom of chat, increase to scroll up)
	Scroll int
	// Cursor position in input (column, row)
	CursorX           int
	CursorY           int
	waitingForCommand bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

// main entry point for the UI
func NewModel(agent *llmServer.Agent) Model {
	return Model{
		Model:   model.NewModel(),
		agent:   agent,
		CursorX: 0,
		CursorY: 0,
	}
}
