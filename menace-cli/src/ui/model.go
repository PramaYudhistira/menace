package ui

import (
	"menace-go/llmServer"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/cellbuf"
	"github.com/mattn/go-runewidth"
)

type Model struct {
	Input    string
	Messages []Message
	agent    *llmServer.Agent
	Width    int
	Height   int
	// Scroll offset (0 = bottom of chat, increase to scroll up)
	Scroll int
	// Cursor position in input (column, row)
	CursorX           int
	CursorY           int
	waitingForCommand bool
	WindowStart       int // Global horizontal window start for all lines
	IsHighlighting    bool
	// Selection state
	SelectionStartX int
	SelectionStartY int
	SelectionEndX   int
	SelectionEndY   int
	HasSelection    bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

// Clears the input state
func (m *Model) ClearState() {
	m.Input = ""
	m.CursorX = 0
	m.CursorY = 0
	m.Scroll = 0
}

// Resets window size and scroll position when terminal is resized.
//
// Returns new Model.
func (m *Model) ResizeWindow(msg tea.Msg) {
	m.Width = msg.(tea.WindowSizeMsg).Width
	m.Height = msg.(tea.WindowSizeMsg).Height
	m.Scroll = 0
}

// UpdateWindowStart ensures the input window is scrolled so the cursor is always visible (applies to all lines).
func (m *Model) UpdateWindowStart(maxInputW int) {
	start := m.WindowStart
	if m.CursorX < start {
		start = m.CursorX
	} else if m.CursorX >= start+maxInputW {
		start = m.CursorX - maxInputW + 1
	}
	if start < 0 {
		start = 0
	}
	m.WindowStart = start
}

// InsertNewLine inserts a new line at the cursor position and moves the cursor to the start of the new line.
func (m *Model) InsertNewLine() {
	lines := strings.Split(m.Input, "\n")
	// Split current line at cursor
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
}

// InsertCharacter inserts a character at the cursor position.
func (m *Model) InsertCharacter(character string) {
	lines := strings.Split(m.Input, "\n")
	// Ensure CursorY is a valid index
	for len(lines) <= m.CursorY {
		lines = append(lines, "")
	}
	//Get current line and convert to runes (Unicode handling)
	runes := []rune(lines[m.CursorY])
	ch := []rune(character)[0]
	newLine := string(runes[:m.CursorX]) + string(ch) + string(runes[m.CursorX:])
	lines[m.CursorY] = newLine
	m.CursorX++
	m.Input = strings.Join(lines, "\n")
}

// Adds a user message to the chat history
func (m *Model) AddUserMessage(message string) {
	m.Messages = append(m.Messages, Message{Sender: "user", Content: message})
}

// Adds a system message to the chat history
func (m *Model) SystemMessage(message string) {
	m.Messages = append(m.Messages, Message{Sender: "system", Content: message})
}

// Handle mouse scrolling
func (m *Model) HandleScroll(direction int) {
	if direction > 0 { // Scroll up
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
	} else { // Scroll down
		if m.Scroll > 0 {
			m.Scroll--
		}
	}
}

// Handles backspace key press
func (m *Model) HandleBackSpace() {
	// Delete character before cursor or merge lines
	lines := strings.Split(m.Input, "\n")
	if m.CursorX > 0 {
		// Delete one character before cursor
		runes := []rune(lines[m.CursorY])
		lines[m.CursorY] = string(runes[:m.CursorX-1]) + string(runes[m.CursorX:])
		m.CursorX--
	} else if m.CursorY > 0 {
		// Merge with previous line
		prev := lines[m.CursorY-1]
		m.CursorX = len([]rune(prev))
		lines[m.CursorY-1] += lines[m.CursorY]
		// Remove current line
		lines = append(lines[:m.CursorY], lines[m.CursorY+1:]...)
		m.CursorY--
	}
	m.Input = strings.Join(lines, "\n")
}

// Handles horizontal cursor movement.
//
// The width of the terminal is not always the same, but it doesn't matter since it purely handles horizontal traversal.
func (m *Model) HandleHorizontalCursorMovement(direction string) {
	// this might not even be the best way to handle lines
	lines := strings.Split(m.Input, "\n")
	if direction == tea.KeyLeft.String() {
		if m.CursorX > 0 {
			m.CursorX--
		} else if m.CursorY > 0 {
			m.CursorY--
			m.CursorX = len([]rune(lines[m.CursorY]))
		}
	} else if direction == tea.KeyRight.String() {
		runes := []rune(lines[m.CursorY])
		if m.CursorX < len(runes) {
			m.CursorX++
		} else if m.CursorY < len(lines)-1 {
			m.CursorY++
			m.CursorX = 0
		}
	}
}

// Handles delete key press
func (m *Model) HandleDelete() {
	lines := strings.Split(m.Input, "\n")
	runes := []rune(lines[m.CursorY])
	if m.CursorX < len(runes) {
		// Remove the character at the cursor position
		lines[m.CursorY] = string(runes[:m.CursorX]) + string(runes[m.CursorX+1:])
		m.Input = strings.Join(lines, "\n")
	} else if m.CursorY < len(lines)-1 {
		// If at the end of the line, merge with the next line
		lines[m.CursorY] += lines[m.CursorY+1]
		lines = append(lines[:m.CursorY+1], lines[m.CursorY+2:]...)
		m.Input = strings.Join(lines, "\n")
	}
}

// GetMaxInputWidth returns the maximum width of the input field.
func (m *Model) GetMaxInputWidth() int {
	prefix := "> "
	boxW := m.Width - 24
	prefixW := runewidth.StringWidth(prefix)
	return boxW - 2 - prefixW
}

// Handles vertical cursor movement.
//
// Maintains the cursor's X position when moving between lines when possible.
func (m *Model) HandleVerticalCursorMovement(direction string) {
	lines := strings.Split(m.Input, "\n")
	if direction == tea.KeyUp.String() {
		if m.CursorY > 0 {
			m.CursorY--
			// Try to maintain X position, but don't exceed line length
			runes := []rune(lines[m.CursorY])
			if m.CursorX > len(runes) {
				m.CursorX = len(runes)
			}
		}
	} else if direction == tea.KeyDown.String() {
		if m.CursorY < len(lines)-1 {
			m.CursorY++
			// Try to maintain X position, but don't exceed line length
			runes := []rune(lines[m.CursorY])
			if m.CursorX > len(runes) {
				m.CursorX = len(runes)
			}
		}
	}
}

// SelectAll selects all text in the input area
func (m *Model) SelectAll() {
	lines := strings.Split(m.Input, "\n")
	if len(lines) > 0 {
		m.HasSelection = true
		m.SelectionStartX = 0
		m.SelectionStartY = 0
		m.SelectionEndX = len([]rune(lines[len(lines)-1]))
		m.SelectionEndY = len(lines) - 1
	}
}

// main entry point for the UI
func NewModel(agent *llmServer.Agent) *Model {
	return &Model{
		agent:   agent,
		CursorX: 0,
		CursorY: 0,
	}
}



