package ui

import (
	"menace-go/llmServer"
	"menace-go/model"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/cellbuf"
	"github.com/mattn/go-runewidth"
)

// LLMResponseMsg represents a message from the LLM
type LLMResponseMsg struct {
	Content string
}

// LoadingMsg represents a loading animation frame
type LoadingMsg struct {
	Frame int
}

// getLLMResponse is a command that fetches a response from the LLM
func getLLMResponse(input string, agent *llmServer.Agent) tea.Cmd {
	return func() tea.Msg {
		// Run the agent to handle the input
		if err := agent.Run(input); err != nil {
			return LLMResponseMsg{Content: "Error: " + err.Error()}
		}

		// Get the final response from the LLM
		llm := llmServer.GetInstance()
		response, err := llm.SendMessage(input)
		if err != nil {
			return LLMResponseMsg{Content: "Error: " + err.Error()}
		}

		if len(response.Choices) > 0 {
			return LLMResponseMsg{Content: response.Choices[0].Message.Content}
		}
		return LLMResponseMsg{Content: "No response from LLM"}
	}
}

// loadingAnimation returns a command that sends loading animation frames
func loadingAnimation() tea.Cmd {
	return tea.Tick(time.Millisecond*300, func(t time.Time) tea.Msg {
		return LoadingMsg{Frame: int((t.UnixNano() / int64(time.Millisecond*300)) % 4)}
	})
}

// Update handles all incoming messages (keypresses, etc.).
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LLMResponseMsg:
		// Remove loading message and add LLM response
		if len(m.Messages) > 0 && strings.HasPrefix(m.Messages[len(m.Messages)-1].Content, "Thinking") {
			m.Messages = m.Messages[:len(m.Messages)-1]
		}
		m.Messages = append(m.Messages, model.Message{Sender: "llm", Content: msg.Content})
		return m, nil

	case LoadingMsg:
		// Update loading animation
		if len(m.Messages) > 0 && strings.HasPrefix(m.Messages[len(m.Messages)-1].Content, "Thinking") {
			dots := strings.Repeat(".", msg.Frame+1)
			m.Messages[len(m.Messages)-1].Content = "Thinking" + dots
			return m, loadingAnimation()
		}
		return m, nil

	case tea.WindowSizeMsg:
		// Update terminal size and reset scroll to bottom
		m.Width = msg.Width
		m.Height = msg.Height
		m.Scroll = 0
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String():
			return m, tea.Quit

		case tea.KeyShiftUp.String(), tea.KeyShiftDown.String():
			// Move cursor up or down one visual (wrapped) row; Enter new line if at bottom
			// Compute inner wrap width
			wrapWidth := m.Width - 24 - 2
			indent := runewidth.StringWidth("> ")
			wrapWidth -= indent
			if wrapWidth < 1 {
				wrapWidth = 1
			}
			// Build visual rows
			type inputRow struct{ LineIndex, Offset, Length int }
			var rows []inputRow
			parts := strings.Split(m.Input, "\n")
			for li, part := range parts {
				runes := []rune(part)
				if len(runes) == 0 {
					rows = append(rows, inputRow{li, 0, 0})
					continue
				}
				for off := 0; off < len(runes); {
					w, start := 0, off
					for ; off < len(runes); off++ {
						rw := runewidth.StringWidth(string(runes[off]))
						if w+rw > wrapWidth {
							break
						}
						w += rw
					}
					if off == start {
						off = start + 1
					}
					rows = append(rows, inputRow{li, start, off - start})
				}
			}
			// Locate current visual row
			oldVis := 0
			for i, r := range rows {
				if r.LineIndex == m.CursorY && m.CursorX >= r.Offset && m.CursorX <= r.Offset+r.Length {
					oldVis = i
					break
				}
			}
			// SHIFT UP
			if msg.String() == tea.KeyShiftUp.String() {
				if oldVis == 0 {
					return m, nil
				}
				newVis := oldVis - 1
				old := rows[oldVis]
				target := rows[newVis]
				rel := m.CursorX - old.Offset
				if rel < 0 {
					rel = 0
				}
				if rel > target.Length {
					rel = target.Length
				}
				m.CursorY = target.LineIndex
				m.CursorX = target.Offset + rel
				return m, nil
			}
			// SHIFT DOWN
			if oldVis == len(rows)-1 {
				// Insert newline at cursor in bottom row
				lines := strings.Split(m.Input, "\n")
				runes := []rune(lines[m.CursorY])
				before := string(runes[:m.CursorX])
				after := string(runes[m.CursorX:])
				newLines := append([]string{}, lines[:m.CursorY+1]...)
				newLines[len(newLines)-1] = before
				newLines = append(newLines, after)
				if m.CursorY+1 < len(lines) {
					newLines = append(newLines, lines[m.CursorY+1:]...)
				}
				m.Input = strings.Join(newLines, "\n")
				m.CursorY++
				m.CursorX = 0
				return m, nil
			}
			newVis := oldVis + 1
			old := rows[oldVis]
			target := rows[newVis]
			rel := m.CursorX - old.Offset
			if rel < 0 {
				rel = 0
			}
			if rel > target.Length {
				rel = target.Length
			}
			m.CursorY = target.LineIndex
			m.CursorX = target.Offset + rel
			return m, nil
		case tea.KeyLeft.String():
			// Move cursor left or up to previous line end
			lines := strings.Split(m.Input, "\n")
			if m.CursorX > 0 {
				m.CursorX--
			} else if m.CursorY > 0 {
				m.CursorY--
				m.CursorX = len([]rune(lines[m.CursorY]))
			}
			return m, nil
		case tea.KeyRight.String():
			// Move cursor right or down to next line start
			lines := strings.Split(m.Input, "\n")
			runes := []rune(lines[m.CursorY])
			if m.CursorX < len(runes) {
				m.CursorX++
			} else if m.CursorY < len(lines)-1 {
				m.CursorY++
				m.CursorX = 0
			}
			return m, nil

		case tea.KeyUp.String():
			// Scroll up the chat history by lines
			// Calculate how many lines fit (adjusted for chat frame and input height)
			visible := m.Height - 7
			if visible < 0 {
				visible = 0
			}
			// Calculate wrap width for content lines
			chatWidth := m.Width - 24
			wrapWidth := chatWidth - 2 // account for ChatStyle padding
			if wrapWidth < 1 {
				wrapWidth = 1
			}
			// Count total lines across all messages after wrapping
			totalLines := 0
			for _, msg := range m.Messages {
				for _, part := range strings.Split(msg.Content, "\n") {
					wrapped := cellbuf.Wrap(part, wrapWidth, "")
					totalLines += strings.Count(wrapped, "\n") + 1
				}
			}
			maxOff := totalLines - visible
			if maxOff < 0 {
				maxOff = 0
			}
			if m.Scroll < maxOff {
				m.Scroll++
			}
			return m, nil
		case tea.KeyDown.String():
			// Scroll down
			if m.Scroll > 0 {
				m.Scroll--
			}
			return m, nil
		case tea.KeyEnter.String():
			if m.Input == "" {
				return m, nil
			}
			// Check if the last message is "Thinking"
			if len(m.Messages) > 0 && strings.HasPrefix(m.Messages[len(m.Messages)-1].Content, "Thinking") {
				return m, nil // Don't allow new message while thinking
			}
			// Append user message
			m.Messages = append(m.Messages, model.Message{Sender: "user", Content: m.Input})
			// Add loading message
			m.Messages = append(m.Messages, model.Message{Sender: "llm", Content: "Thinking"})
			// Get LLM response asynchronously
			cmd := getLLMResponse(m.Input, m.agent)
			// Start loading animation
			loadingCmd := loadingAnimation()
			// Reset input and cursor
			m.Input = ""
			m.CursorX = 0
			m.CursorY = 0
			// Reset scroll to bottom
			m.Scroll = 0
			return m, tea.Batch(cmd, loadingCmd)

		case tea.KeyBackspace.String():
			// Delete character before cursor or merge lines
			lines := strings.Split(m.Input, "\n")
			if m.CursorX > 0 {
				runes := []rune(lines[m.CursorY])
				lines[m.CursorY] = string(runes[:m.CursorX-1]) + string(runes[m.CursorX:])
				m.CursorX--
			} else if m.CursorY > 0 {
				// Merge with previous line
				prev := lines[m.CursorY-1]
				curr := lines[m.CursorY]
				m.CursorX = len([]rune(prev))
				lines[m.CursorY-1] = prev + curr
				lines = append(lines[:m.CursorY], lines[m.CursorY+1:]...)
				m.CursorY--
			}
			m.Input = strings.Join(lines, "\n")
			return m, nil

		case tea.KeyDelete.String():
			// Delete character at cursor or merge lines
			lines := strings.Split(m.Input, "\n")
			if m.CursorX < len([]rune(lines[m.CursorY])) {
				runes := []rune(lines[m.CursorY])
				lines[m.CursorY] = string(runes[:m.CursorX]) + string(runes[m.CursorX+1:])
			} else if m.CursorY < len(lines)-1 {
				// Merge with next line
				curr := lines[m.CursorY]
				next := lines[m.CursorY+1]
				lines[m.CursorY] = curr + next
				lines = append(lines[:m.CursorY+1], lines[m.CursorY+2:]...)
			}
			m.Input = strings.Join(lines, "\n")
			return m, nil

		default:
			if len(msg.String()) == 1 {
				// Insert character at cursor position
				lines := strings.Split(m.Input, "\n")
				runes := []rune(lines[m.CursorY])
				ch := []rune(msg.String())[0]
				newLine := string(runes[:m.CursorX]) + string(ch) + string(runes[m.CursorX:])
				lines[m.CursorY] = newLine
				m.CursorX++
				m.Input = strings.Join(lines, "\n")
			}
			return m, nil
		}
	}
	return m, nil
}
