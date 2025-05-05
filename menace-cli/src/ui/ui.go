package ui

import (
	"menace-go/model"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	model.Model
	Width  int
	Height int
	// Scroll offset (0 = bottom of chat, increase to scroll up)
	Scroll int
}

// Init is called when the program starts.
func (m Model) Init() tea.Cmd {
	_ = model.Model{}
	return nil
}

// main entry point for the UI
func NewModel() Model {
	return Model{
		Model: model.NewModel(),
	}
}
