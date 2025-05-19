package ui

import (
	"menace-go/llmServer"
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
	mf := llmServer.ModelFactory{}
	sidebarContent := HeaderStyle.Render("👹 Menace CLI") +
		"\n" + "Working Directory:" +
		"\n" + dir +
		"\n" + "Running on:" +
		"\n" + mf.DetectShell() +
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
		// measure prefix width and prepare indent for wrapped lines
		prefixWidth := runewidth.StringWidth(prefix)
		prefixIndent := strings.Repeat(" ", prefixWidth)
		firstLine := true
		for _, part := range strings.Split(msg.Content, "\n") {
			// wrap part at content width minus prefix
			lineWidth := wrapWidth - prefixWidth
			if lineWidth < 1 {
				lineWidth = 1
			}
			wrapped := cellbuf.Wrap(part, lineWidth, "")
			for _, line := range strings.Split(wrapped, "\n") {
				if firstLine {
					renderedLines = append(renderedLines, styleFunc(prefix+line))
					firstLine = false
				} else {
					renderedLines = append(renderedLines, styleFunc(prefixIndent+line))
				}
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

	// Render input area with block cursor and proper wrapping/indent
	prefix := "> "
	boxW := termWidth - 24
	prefixW := runewidth.StringWidth(prefix)
	// maxInputW := boxW - 2 - prefixW // no longer needed

	// Multiline input rendering with global horizontal scroll
	lines := strings.Split(m.Input, "\n")
	indent := strings.Repeat(" ", prefixW)
	maxInputW := boxW - 2 - prefixW
	var rendered []string
	for i, line := range lines {
		curPfx := indent
		if i == 0 {
			curPfx = prefix
		}
		runes := []rune(line)
		windowStart := m.WindowStart
		windowEnd := windowStart + maxInputW
		if windowEnd > len(runes) {
			windowEnd = len(runes)
		}
		visible := runes[windowStart:windowEnd]
		if i == m.CursorY {
			// Place block cursor at the correct position within the visible window
			cursorInWindow := m.CursorX - windowStart
			if cursorInWindow >= 0 && cursorInWindow <= len(visible) {
				before := string(visible[:cursorInWindow])
				ch := " "
				if cursorInWindow < len(visible) {
					ch = string(visible[cursorInWindow])
				}
				after := ""
				if cursorInWindow+1 <= len(visible) {
					after = string(visible[cursorInWindow+1:])
				}
				lineStr := before + lipgloss.NewStyle().Reverse(true).Render(ch) + after
				rendered = append(rendered, curPfx+lineStr)
				continue
			}
		}
		rendered = append(rendered, curPfx+string(visible))
	}
	inputContent := strings.Join(rendered, "\n")
	inputPrompt := InputStyle.
		Border(lipgloss.RoundedBorder()).
		Width(boxW).
		Render(inputContent)

	// Combine chat and input
	mainArea := lipgloss.JoinVertical(lipgloss.Top, chatBox, inputPrompt)

	// Layout: sidebar + main area
	screen := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mainArea)

	// Add horizontal margins; place UI flush to top so chat-box top border is visible
	return lipgloss.NewStyle().
		Margin(0, 2). // top/bottom, left/right
		Render(screen)
}
