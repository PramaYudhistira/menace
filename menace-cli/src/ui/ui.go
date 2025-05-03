package ui

import (
	"menace-go/model"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	model.Model
	Width  int
	Height int
	// Scroll offset (0 = bottom of chat, increase to scroll up)
	Scroll int
}

// Init is called when the program starts.
func (m Model) Init() tea.Cmd {
	_ = model.Model{}
	return nil
}

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
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			// Scroll up the chat history
			// Calculate how many messages fit (adjusted for chat frame and input height)
			visible := m.Height - 7
			if visible < 0 {
				visible = 0
			}
			maxOff := len(m.Messages) - visible
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
		case "enter":
			// Append user message
			//TODO: Call LLM, send request via m.Input
			m.Messages = append(m.Messages, model.Message{Sender: "user", Content: m.Input})
			response := "Echo (LLM output here...): " + m.Input
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
		"\n" + dir +
		"\n`Ctrl + C` to exit" +
		"\n`↑`/`↓` to scroll"
	sidebar := SidebarStyle.Render(sidebarContent)
	sidebar = lipgloss.NewStyle().
		Align(lipgloss.Left, lipgloss.Top).
		Width(18).              // match updated sidebar width
		Height(termHeight - 2). // Adjust for top/bottom margins
		Render(sidebar)

	// Render messages with distinct styles, applying scroll offset
	// Determine how many messages fit in chat box
	visibleCount := termHeight - 7
	if visibleCount < 0 {
		visibleCount = 0
	}
	// Select slice of messages based on scroll (0 = bottom)
	total := len(m.Messages)
	var slice []model.Message
	if total > visibleCount {
		start := total - visibleCount - m.Scroll
		if start < 0 {
			start = 0
		}
		end := total - m.Scroll
		if end > total {
			end = total
		}
		slice = m.Messages[start:end]
	} else {
		slice = m.Messages
	}
	var rendered []string
	for _, msg := range slice {
		switch msg.Sender {
		case "user":
			rendered = append(rendered, UserStyle.Render("> "+msg.Content))
		case "llm":
			rendered = append(rendered, LLMStyle.Render("💭 "+msg.Content))
		default:
			rendered = append(rendered, msg.Content)
		}
	}
	chatBody := lipgloss.JoinVertical(lipgloss.Top, rendered...)
	chatBox := ChatStyle.
		Width(termWidth - 24).  // Adjust width to fit next to the sidebar
		Height(termHeight - 5). // Leave space for input box
		Render(chatBody)

		// Input prompt with its own rectangle
	inputPrompt := InputStyle.
		Border(lipgloss.RoundedBorder()).
		Width(termWidth - 22). // Same width as chat box
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

// main entry point for the UI
func NewModel() Model {
	return Model{
		Model: model.NewModel(),
	}
}
