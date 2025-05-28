package ui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"menace-go/llmServer"

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
		// handle execution of command when awaiting command approval
		if m.AwaitingCommandApproval {
			switch msg.String() {
			case "y":
				m.AwaitingCommandApproval = false
				var output string
				var err error
				if m.PendingFunctionCall != nil {
					switch m.PendingFunctionCall.Name {
					case "ReadFileWithLineNumbers":
						path, _ := m.PendingFunctionCall.Args["path"].(string)
						output, err = ReadFileWithLineNumbers(path)
					case "CreateAndApplyDiffs":
						path, _ := m.PendingFunctionCall.Args["path"].(string)
						diffsRaw, _ := m.PendingFunctionCall.Args["diffs"].([]interface{})
						var diffs []LineDiff
						for _, d := range diffsRaw {
							if diffMap, ok := d.(map[string]interface{}); ok {
								diff := LineDiff{}
								if t, ok := diffMap["Type"].(float64); ok {
									diff.Type = DiffType(int(t))
								}
								if idx, ok := diffMap["LineIndex"].(float64); ok {
									diff.LineIndex = int(idx)
								}
								if oldC, ok := diffMap["OldContent"].(string); ok {
									diff.OldContent = oldC
								}
								if newC, ok := diffMap["NewContent"].(string); ok {
									diff.NewContent = newC
								}
								diffs = append(diffs, diff)
							}
						}
						err = CreateAndApplyDiffs(path, diffs)
						if err == nil {
							output = "Diffs applied successfully."
						}
					}
					m.PendingFunctionCall = nil
				} else if m.PendingCommand != nil {
					output, err = runShellCommand(m.PendingCommand.Command)
					m.PendingCommand = nil
				}
				if err != nil {
					m.AddSystemMessage(fmt.Sprintf("Error: %s", err))
				} else {
					cleanOutput := strings.ReplaceAll(output, "\r\n", "\n")
					cleanOutput = strings.ReplaceAll(cleanOutput, "\r", "\n")
					cleanOutput = strings.ReplaceAll(cleanOutput, "\t", "    ")
					m.AddSystemMessage(fmt.Sprintf("Output:\n%s", cleanOutput))
				}
				m.StartThinking()
				return m, tea.Batch(
					func() tea.Msg {
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

		// Case for Enter key press -- START OF DEBUGGING
		// Should send message to LLM
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
			// All the return types are re-entered into this large switch statement
			// Check case CommandSuggestionMsg, LLMResponseMsg, and SystemMessage for more details
			return m, tea.Batch(
				func() tea.Msg {
					response, cmdSuggestion, err := m.agent.SendMessage(
						context.Background(), 
						userInput,
					)
					// fmt.Println("userInput: ", userInput)
					// fmt.Println("response: ", response)
					if err != nil {
						return SystemMessage{Content: "Error: " + err.Error()}
					}
					if cmdSuggestion != nil {
						return CommandSuggestionMsg{
							Command: cmdSuggestion.Command, 
							Reason: cmdSuggestion.Reason, 
							AwaitingCommandApproval: cmdSuggestion.AwaitingCommandApproval == "true",
						}
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

	// Adding extra context prior to actually executing the commands, think of this as pre-run add-ons
	case CommandSuggestionMsg:
		m.PendingCommand = &msg
		m.PendingFunctionCall = nil
		m.AwaitingCommandApproval = msg.AwaitingCommandApproval
		m.StopThinking()
		m.AddAgentMessage(fmt.Sprintf("Explanation: %s", msg.Reason))

		// For git commands, the LLM gets extra context to guide it to its next step
		if strings.HasPrefix(msg.Command, "git add") {
			_, adds, _ := llmServer.HasChanges()
			m.agent.AddToMessageChain(fmt.Sprintf("Your next step should be to commit, only if the user asks to commit or beyond (push or pr). Here are the changes so far: %s", adds), "")
		} else if strings.HasPrefix(msg.Command, "git commit") {
			m.agent.AddToMessageChain("Your next step should be to push, only if the user asks to push or beyond (pr)", "")
		} else if strings.HasPrefix(msg.Command, "git push") {
			m.agent.AddToMessageChain("Your next step should be to create a pull request, only if the user asks to create a pull request", "")
		}

		// Not all commands needs human intervention, so we can skip the command execution handled in case SkipStepMsg
		if m.AwaitingCommandApproval {
			m.AddSystemMessage(fmt.Sprintf("Command suggestion: %s\nExecute command? (y/n/e)", msg.Command))
			return m, nil
		} else {
			return m, tea.Batch(
				func() tea.Msg {
					m.StartThinking()
					return SkipStepMsg{Command_to_execute: &msg, Function_to_execute: nil}
				},
				thinkingTick(),
			)
		}

	case LLMResponseMsg:
		// Try to parse for a function call
		if fnCall := parseFunctionCall(msg.Content); fnCall != nil {
			m.PendingFunctionCall = fnCall
			m.PendingCommand = nil
			m.AwaitingCommandApproval = true
			m.StopThinking()
			m.AddAgentMessage(fmt.Sprintf("Explanation: %s", fnCall.Reason))
			// If it's a diff, show the diff preview
			if fnCall.Name == "CreateAndApplyDiffs" {
				path, _ := fnCall.Args["path"].(string)
				diffsRaw, _ := fnCall.Args["diffs"].([]interface{})
				var diffs []LineDiff
				for _, d := range diffsRaw {
					if diffMap, ok := d.(map[string]interface{}); ok {
						diff := LineDiff{}
						if t, ok := diffMap["Type"].(float64); ok {
							diff.Type = DiffType(int(t))
						}
						if idx, ok := diffMap["LineIndex"].(float64); ok {
							diff.LineIndex = int(idx)
						}
						if oldC, ok := diffMap["OldContent"].(string); ok {
							diff.OldContent = oldC
						}
						if newC, ok := diffMap["NewContent"].(string); ok {
							diff.NewContent = newC
						}
						diffs = append(diffs, diff)
					}
				}
				// Format the diff for display
				var preview strings.Builder
				preview.WriteString(fmt.Sprintf("Proposed changes to %s:\n", path))
				for _, d := range diffs {
					preview.WriteString(FormatDiff(d) + "\n")
				}
				m.AddSystemMessage(preview.String())
			} else if fnCall.Name == "createPullRequest" {
				branchName, _ := fnCall.Args["branch_name"].(string)
				title, _ := fnCall.Args["title"].(string)
				summary, _ := fnCall.Args["summary"].(string)
				err := llmServer.CreatePullRequest(branchName, title, summary)
				if err != nil {
					m.AddSystemMessage(fmt.Sprintf("Error: %s", err))
					m.agent.AddToMessageChain(fmt.Sprintf("Oops! An error occured. Error: %s. Please fix this and try again", err), "")
				}
			}
			m.AddSystemMessage(fmt.Sprintf("Function call suggestion: %s\nExecute function? (y/n/e)", fnCall.Name))
			return m, nil
		}
		m.StopThinking()
		m.AddAgentMessage(msg.Content)
		return m, nil

	case SystemMessage:
		m.StopThinking()
		m.AddSystemMessage(msg.Content)
		return m, nil

	case SkipStepMsg:
		m.StopThinking()
		m.AddSystemMessage("Skipping human intervention...")
		if msg.Command_to_execute != nil {
			m.PendingCommand = msg.Command_to_execute
			output, err := runShellCommand(m.PendingCommand.Command)
			cleanOutput := strings.ReplaceAll(output, "\r\n", "\n")
			cleanOutput = strings.ReplaceAll(cleanOutput, "\r", "\n")
			cleanOutput = strings.ReplaceAll(cleanOutput, "\t", "    ")
			if err != nil {
				m.AddSystemMessage(fmt.Sprintf("Error: %s", err))
			} else {
				m.AddSystemMessage(fmt.Sprintf("Output:\n%s", cleanOutput))
			}
			m.StartThinking()
			m.PendingCommand = nil
			return m, tea.Batch(
				func() tea.Msg {
					response, cmdSuggestion, err := m.agent.SendMessage(
						context.Background(), 
						fmt.Sprintf("Command %s executed. Output: %s", msg.Command_to_execute.Command, output),
					)
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
		} else if msg.Function_to_execute != nil {
			m.PendingFunctionCall = msg.Function_to_execute
		}
		m.PendingFunctionCall = nil
		
	}

	if changed {
		m.UpdateWindowStart(m.GetMaxInputWidth())
	}

	return m, nil
}
