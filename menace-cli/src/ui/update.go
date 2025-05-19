package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all incoming messages (keypresses, etc.).
//
// Part of Bubble Tea Model interface
// runs every time a key is pressed
//
// each key press is represented by msg Type tea.Msg
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//flag to check if input or cursor position was modified during handling of keypress
	changed := false

	switch msg := msg.(type) {

	// Styles to fit terminal size
	case tea.WindowSizeMsg:
		m.ResizeWindow(msg)

	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			// fmt.Println("Mouse wheel up detected")
			m.HandleScroll(1)
		case tea.MouseButtonWheelDown:
			// fmt.Println("Mouse wheel down detected")
			m.HandleScroll(-1)
		}

	// Handle key presses
	case tea.KeyMsg:
		switch msg.String() {
		//if ctrl+c is pressed, quit the program
		case tea.KeyCtrlC.String():
			return m, tea.Quit

		//Case for Enter key press
		//Should send message to LLM
		case tea.KeyEnter.String():
			if m.Input == "" {
				return m, nil
			}
			m.AddUserMessage(m.Input)
			m.ClearState()

		//Case for horizontal cursor movement
		case tea.KeyLeft.String(), tea.KeyRight.String():
			m.HandleHorizontalCursorMovement(msg.String())
			changed = true

		//Case for backspace key press
		case tea.KeyBackspace.String():
			m.HandleBackSpace()
			changed = true

		//Case for delete key press
		case tea.KeyDelete.String(), tea.KeyCtrlD.String():
			m.HandleDelete()
			changed = true

		//general key press
		//Inserts single character input into the cursor position
		default:
			if len(msg.String()) == 1 {
				m.InsertCharacter(msg.String())
				changed = true
			}
		}
	}

	if changed {
		m.UpdateWindowStart(m.GetMaxInputWidth())
	}

	return m, nil
}
