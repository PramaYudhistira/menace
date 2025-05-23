package ui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
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
	// Handle thinking animation
	case ThinkingMsg:
		return m.UpdateThinking()

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
		case tea.MouseButtonLeft:
			if msg.Action == tea.MouseActionRelease {
				// Clear any existing selection
				m.IsHighlighting = false
				m.SelectionStartX = 0
				m.SelectionStartY = 0
				m.SelectionEndX = 0
				m.SelectionEndY = 0

				if zone.Get("help").InBounds(msg) {
					m.AddSystemMessage("testing system message")
					return m, nil
				}
				if zone.Get("config").InBounds(msg) {
					m.OpenConfig()
					return m, nil
				}
			}
		}

	// Handle key presses
	case tea.KeyMsg:

		if m.IsConfigOpen {
			switch msg.String() {
			case tea.KeyEnter.String():
				m.SelectModel()
				changed = true
			case tea.KeyEsc.String():
				m.CloseConfig()
				changed = true
			case tea.KeyUp.String(), tea.KeyDown.String():
				m.HandleConfigNavigation(msg.String())
				changed = true
			}
			return m, nil
		}

		if m.AwaitingCommandApproval {
			switch msg.String() {
			case "y":
				// Execute the command
				m.AwaitingCommandApproval = false
				output, err := runShellCommand(m.PendingCommand.Command)
				if err != nil {
					// TODO: Consider implementing functionality to send error details to the LLM for self-improvement.
					m.AddSystemMessage(fmt.Sprintf("Error: %s", err))
				}
				cleanOutput := strings.ReplaceAll(output, "\r\n", "\n")
				cleanOutput = strings.ReplaceAll(cleanOutput, "\r", "\n")
				cleanOutput = strings.ReplaceAll(cleanOutput, "\t", "    ") // replace tabs with spaces
				m.AddSystemMessage(fmt.Sprintf("Command output:\n%s", cleanOutput))
				m.StartThinking()
				return m, tea.Batch(
					func() tea.Msg {
						// truncate output?
						response, cmdSuggestion, err := m.agent.SendMessage(context.Background(), output)
						if err != nil {
							return SystemMessage{Content: "Error: " + err.Error()}
						}
						if cmdSuggestion != nil {
							return CommandSuggestionMsg{Command: cmdSuggestion.Command, Reason: cmdSuggestion.Reason}
						}
						return LLMResponseMsg{Content: response}
					},
					thinkingTick(),
				)

			case "n":
				// Cancel
				m.AwaitingCommandApproval = false
				m.PendingCommand = nil
				m.AddAgentMessage("Command Cancelled.")
				m.StartThinking()
				return m, tea.Batch(
					func() tea.Msg {

						response, cmdSuggestion, err := m.agent.SendMessage(context.Background(), "No, stop for now.")
						if err != nil {
							return SystemMessage{Content: "Error: " + err.Error()}
						}
						if cmdSuggestion != nil {
							// Try not to get to this case...
							return CommandSuggestionMsg{Command: cmdSuggestion.Command, Reason: cmdSuggestion.Reason}
						}
						return LLMResponseMsg{Content: response}
					},
					thinkingTick(),
				)
			case "e":
				// Switch to edit mode (maybe put command in input box)
				m.Input = m.PendingCommand.Command
				m.AwaitingCommandApproval = false
				m.PendingCommand = nil

				return m, nil
			}
			return m, nil
		}

		switch msg.String() {
		//if ctrl+c is pressed, quit the program
		case tea.KeyCtrlC.String():
			if m.IsHighlighting {
				selectedText := m.GetSelectedText()
				m.CopyToClipboard(selectedText)
				// Clear selection after copying
				m.IsHighlighting = false
				m.SelectionStartX = 0
				m.SelectionStartY = 0
				m.SelectionEndX = 0
				m.SelectionEndY = 0
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyCtrlX.String():
			if m.IsHighlighting {
				m.CutSelectedText()
				changed = true
			}

		case tea.KeyCtrlV.String():
			clipboardContent := m.GetClipboardContent()
			if clipboardContent != "" {
				// Normalize clipboard content
				clipboardContent = strings.ReplaceAll(clipboardContent, "\r\n", "\n")
				clipboardContent = strings.ReplaceAll(clipboardContent, "\r", "\n")

				// If text is selected, replace it
				if m.IsHighlighting {
					m.Input = ""
					m.CursorX = 0
					m.CursorY = 0
					m.IsHighlighting = false
					m.SelectionStartX = 0
					m.SelectionStartY = 0
					m.SelectionEndX = 0
					m.SelectionEndY = 0
				}
				// Insert clipboard content
				lines := strings.Split(clipboardContent, "\n")
				for i, line := range lines {
					for _, char := range line {
						m.InsertCharacter(string(char))
					}
					// Insert a new line if not the last line
					if i < len(lines)-1 {
						m.InsertNewLine()
					}
				}
				changed = true
			}

		//Case for Enter key press
		//Should send message to LLM
		case tea.KeyEnter.String():
			if m.Input == "" {
				return m, nil
			}

			// Add user message to UI
			m.AddUserMessage(m.Input)

			// Start thinking animation
			m.StartThinking()

			// Capture input before clearing
			userInput := m.Input

			// Clear input
			m.ClearState()

			// Send to agent and get response asynchronously via Bubble Tea command
			return m, tea.Batch(
				func() tea.Msg {
					response, cmdSuggestion, err := m.agent.SendMessage(context.Background(), userInput)
					if err != nil {
						return SystemMessage{Content: "Error: " + err.Error()}
					}
					if cmdSuggestion != nil {
						return CommandSuggestionMsg{Command: cmdSuggestion.Command, Reason: cmdSuggestion.Reason}
					}
					return LLMResponseMsg{Content: response}
				},
				thinkingTick(),
			)

		//Case for horizontal cursor movement
		case tea.KeyLeft.String(), tea.KeyRight.String():
			m.HandleHorizontalCursorMovement(msg.String())
			changed = true

		//Case for vertical cursor movement
		case tea.KeyUp.String(), tea.KeyDown.String():
			m.HandleVerticalCursorMovement(msg.String())
			changed = true

		//Case for backspace key press
		case tea.KeyBackspace.String():
			m.HandleBackSpace()
			changed = true

		//Case for delete key press
		case tea.KeyDelete.String(), tea.KeyCtrlD.String():
			m.HandleDelete()
			changed = true

		//Case for newline key press
		case tea.KeyShiftDown.String():
			m.InsertNewLine()
			changed = true

		//Case for ctrl A key press
		case tea.KeyCtrlA.String():
			m.IsHighlighting = true
			m.SelectAll()
			changed = true

		//general key press
		//Inserts single character input into the cursor position
		default:
			if len(msg.String()) == 1 {
				if m.IsHighlighting {
					m.Input = ""
					m.CursorX = 0
					m.CursorY = 0
					m.IsHighlighting = false
					m.SelectionStartX = 0
					m.SelectionStartY = 0
					m.SelectionEndX = 0
					m.SelectionEndY = 0
				}
				m.InsertCharacter(msg.String())
				changed = true
			}
		}

	case CommandSuggestionMsg:
		m.PendingCommand = &msg
		m.AwaitingCommandApproval = true
		m.StopThinking()
		m.AddAgentMessage(fmt.Sprintf("Explanation: %s", msg.Reason))
		//TODO: Better UX please
		m.AddSystemMessage(fmt.Sprintf("Command suggestion: %s\nExecute command? (y/n/e)", msg.Command))
		return m, nil

	case LLMResponseMsg:
		m.StopThinking()
		m.AddAgentMessage(msg.Content)
		return m, nil

	case SystemMessage:
		m.StopThinking()
		m.AddSystemMessage(msg.Content)
		return m, nil
	}

	if changed {
		m.UpdateWindowStart(m.GetMaxInputWidth())
	}

	return m, nil
}
