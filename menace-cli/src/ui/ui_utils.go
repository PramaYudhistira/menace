package ui

import (
	"menace-go/llmServer"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LoadingMsg represents a loading animation frame
type LoadingMsg struct {
	Frame int
}

// LLMResponseMsg represents a message from the LLM
type LLMResponseMsg struct {
	Content string
}

// Styles for Menace CLI UI
// Includes "boxes" for chat and input area
var (
	SidebarStyle = lipgloss.NewStyle().
			Width(20).
			PaddingRight(1).
			Foreground(lipgloss.Color("15"))

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10")).
			MarginBottom(1)

	ChatStyle = lipgloss.NewStyle().
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("12"))

	InputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6")).
			Bold(true)

	UserStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true)

	LLMStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("4")).
			Italic(true)
	// Style for the block cursor in the input
	CursorStyle = lipgloss.NewStyle().Reverse(true)
	shellType   = strings.Split(llmServer.ModelFactory{}.DetectShell(), "/")[1]
)

// loadingAnimation returns a command that sends loading animation frames
func loadingAnimation() tea.Cmd {
	return tea.Tick(time.Millisecond*300, func(t time.Time) tea.Msg {
		return LoadingMsg{Frame: int((t.UnixNano() / int64(time.Millisecond*300)) % 4)}
	})
}
