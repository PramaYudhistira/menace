package ui

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/cellbuf"
	"github.com/mattn/go-runewidth"
)

// View renders the entire UI.
func (m Model) View() string {
	termWidth := m.Width
	termHeight := m.Height

	if termWidth == 0 || termHeight == 0 {
		termWidth = 80  // Default width
		termHeight = 20 // Default height
	}

	// Sidebar with header
	dir, err := os.Getwd()
	if err != nil {
		dir = "unknown directory"
	}
	// Sidebar header and controls
	sidebarContent := HeaderStyle.Render("👹 Menace CLI") +
		"\n" + "Working Directory:" +
		"\n" + dir +
		"\n`Ctrl + C` to exit" +
		"\n`↑`/`↓` to scroll"
	sidebar := SidebarStyle.Render(sidebarContent)
	sidebar = lipgloss.NewStyle().
		Align(lipgloss.Left, lipgloss.Top).
		Width(18).              // match updated sidebar width
		Height(termHeight - 2). // Adjust for top/bottom margins
		Render(sidebar)

	// Render messages line by line with wrapping, applying scroll offset
	var renderedLines []string
	// chat content width: total width minus sidebar and margins
	chatWidth := termWidth - 24
	// account for ChatStyle padding (1 left, 1 right)
	wrapWidth := chatWidth - 2
	if wrapWidth < 1 {
		wrapWidth = 1
	}
	for _, msg := range m.Messages {
		var styleFunc func(...string) string
		var prefix string
		switch msg.Sender {
		case "user":
			styleFunc = UserStyle.Render
			prefix = "> "
		case "llm":
			styleFunc = LLMStyle.Render
			prefix = "💭 "
		default:
			styleFunc = func(parts ...string) string { return strings.Join(parts, "") }
			prefix = ""
		}
		// measure prefix width
		prefixWidth := runewidth.StringWidth(prefix)
		for _, part := range strings.Split(msg.Content, "\n") {
			// wrap part at content width minus prefix
			lineWidth := wrapWidth - prefixWidth
			if lineWidth < 1 {
				lineWidth = 1
			}
			wrapped := cellbuf.Wrap(part, lineWidth, "")
			for _, line := range strings.Split(wrapped, "\n") {
				renderedLines = append(renderedLines, styleFunc(prefix+line))
			}
		}
	}
	// Determine how many lines fit in chat box
	visibleLines := termHeight - 7
	if visibleLines < 0 {
		visibleLines = 0
	}
	// Slice lines based on scroll (0 = bottom)
	totalLines := len(renderedLines)
	var linesToRender []string
	if totalLines > visibleLines {
		start := totalLines - visibleLines - m.Scroll
		if start < 0 {
			start = 0
		}
		end := totalLines - m.Scroll
		if end > totalLines {
			end = totalLines
		}
		linesToRender = renderedLines[start:end]
	} else {
		linesToRender = renderedLines
	}
	chatBody := lipgloss.JoinVertical(lipgloss.Top, linesToRender...)
	chatBox := ChatStyle.
		Width(termWidth - 24).  // Adjust width to fit next to the sidebar
		Height(termHeight - 5). // Leave space for input box
		Render(chatBody)

		// Input prompt with its own rectangle
	inputPrompt := InputStyle.
		Border(lipgloss.RoundedBorder()).
		Width(termWidth - 24). // Same width as chat box
		Render("> " + m.Input)

	// Combine chat and input
	mainArea := lipgloss.JoinVertical(lipgloss.Top, chatBox, inputPrompt)

	// Layout: sidebar + main area
	screen := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mainArea)

	// Add horizontal margins; place UI flush to top so chat-box top border is visible
	return lipgloss.NewStyle().
		Margin(0, 2). // top/bottom, left/right
		Render(screen)
}
