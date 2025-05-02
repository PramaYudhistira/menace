package ui

import (
	"menace-go/model"

	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	model.Model
}

// Init is called when the program starts.
func (m Model) Init() tea.Cmd {
	_ = model.Model{}
	return nil
}

// Update handles all incoming messages (keypresses, etc.).
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			// Append user message
			m.Messages = append(m.Messages, model.Message{Sender: "user", Content: m.Input})
			// TODO: Replace with real LLM call
			response := "Echo (LLM output here...): " + m.Input
			m.Messages = append(m.Messages, model.Message{Sender: "llm", Content: response})
			m.Input = ""
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
	// Sidebar with header
	dir, err := os.Getwd()
	if err != nil {
		dir = "unknown directory"
	}
	sidebarContent := HeaderStyle.Render("👹 Menace CLI") + "\n" + dir + "\n`Ctrl + C` to exit"
	sidebar := SidebarStyle.Render(sidebarContent)
	sidebar = lipgloss.NewStyle().Align(lipgloss.Top, lipgloss.Left).
		Width(20).
		Height(10).
		Render(sidebar)

	// Render messages with distinct
	var rendered []string
	for _, msg := range m.Messages {
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
	chatBox := ChatStyle.Render(chatBody)

	// Input prompt
	inputPrompt := InputStyle.Render("> " + m.Input)

	// Combine chat and input
	mainArea := lipgloss.JoinVertical(lipgloss.Top, chatBox, inputPrompt)

	mainArea = lipgloss.NewStyle().
		Width(80).
		Render(mainArea)

	// Layout: sidebar + main area
	screen := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mainArea)

	// Add margin
	return lipgloss.NewStyle().Margin(1, 2).Render(screen)
}

// main entry point for the UI
func NewModel() Model {
	return Model{
		Model: model.NewModel(),
	}
}
