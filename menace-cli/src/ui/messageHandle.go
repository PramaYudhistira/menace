package ui

import (
	"menace-go/model"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/cellbuf"
)

// Update handles all incoming messages (keypresses, etc.).
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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

		case tea.KeyShiftDown.String():
			// Handle Ctrl+C
			m.Input += "\n"
			return m, nil

		case "up":
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
		case "down":
			// Scroll down
			if m.Scroll > 0 {
				m.Scroll--
			}
			return m, nil
		case tea.KeyEnter.String():
			// Append user message
			//TODO: Call LLM, send request via m.Input
			m.Messages = append(m.Messages, model.Message{Sender: "user", Content: m.Input})
			response := "Echo (LLM output here...) haha: " + m.Input
			m.Messages = append(m.Messages, model.Message{Sender: "llm", Content: response})
			m.Input = ""
			// reset scroll to bottom on new message
			m.Scroll = 0
			return m, nil

		case "backspace":
			if len(m.Input) > 0 {
				m.Input = m.Input[:len(m.Input)-1]
			}
			return m, nil

		default:
			if len(msg.String()) == 1 {
				m.Input += msg.String()
			}
			return m, nil
		}
	}
	return m, nil
}
