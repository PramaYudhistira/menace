package ui

import (
	"menace-go/llmServer"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/cellbuf"
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
	WindowStart       int // Start of the visible input window
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

// Inserts a character at the cursor position
func (m *Model) InsertCharacter(character string) {
	//TODO: perhaps there is a better way to handle detection of lines?
	lines := strings.Split(m.Input, "\n")

	//Get current line and convert to runes (Unicode handling)
	runes := []rune(lines[m.CursorY])

	//get character that was pressed and convert to rune, to support emojis, accented characters and so on
	ch := []rune(character)[0]

	//actually insert the character at cursor position
	newLine := string(runes[:m.CursorX]) + string(ch) + string(runes[m.CursorX:])

	//update the new line on the current line
	lines[m.CursorY] = newLine
	m.CursorX++
	m.Input = strings.Join(lines, "\n")
}

// Adds a user message to the chat history
func (m *Model) AddUserMessage(message string) {
	m.Messages = append(m.Messages, Message{Sender: "user", Content: message})
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
		runes := []rune(lines[m.CursorY])
		lines[m.CursorY] = string(runes[:m.CursorX-1])
		m.CursorX--
	} else if m.CursorY > 0 {
		// Merge with previous line
		prev := lines[m.CursorY-1]
		m.CursorX = len([]rune(prev))
		lines[m.CursorY-1] += lines[m.CursorY]

		// Spread operator to ensure future lines are also updated in their respective positions
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

// main entry point for the UI
func NewModel(agent *llmServer.Agent) *Model {
	return &Model{
		agent:   agent,
		CursorX: 0,
		CursorY: 0,
	}
}
