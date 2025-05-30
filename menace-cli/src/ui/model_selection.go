package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type ModelInfo struct {
	Provider string
	ID       string
}

var ClosedSourceModels = map[string]ModelInfo{
	"GPT 4.1": {Provider: "openai", ID: "gpt-4-0125-preview"},
	"GPT 3.5": {Provider: "openai", ID: "gpt-3.5-turbo"},
	"o4-mini": {Provider: "openai", ID: "o4-mini-2025-04-16"},
	"Claude":  {Provider: "anthropic", ID: "claude-3-opus-20240229"},
}

var ModelKeys = []string{}

var AvailableModels = map[string]ModelInfo{}

func (m *Model) ConfigView(termHeight, termWidth int) string {
	var configContent strings.Builder
	configContent.WriteString(HeaderStyle.Render("Select Model") + "\n\n")

	for i, model := range ModelKeys {
		style := lipgloss.NewStyle()
		if i == m.ConfigCursor {
			style = style.
				Foreground(lipgloss.Color("#8be9fd")).
				Bold(true)
		}
		configContent.WriteString(style.Render("> "+model) + "\n")
	}

	configContent.WriteString("\n" + HeaderStyle.Render("Controls:"))
	configContent.WriteString("\n↑/↓: Navigate")
	configContent.WriteString("\nEnter: Select")
	configContent.WriteString("\nEsc: Back")

	configBox := lipgloss.NewStyle().
		Width(termWidth - 24).
		Height(termHeight - 5).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Render(configContent.String())

	return zone.Scan(lipgloss.NewStyle().
		Margin(0, 2).
		Render(configBox))
}

// OpenConfig opens the config page
func (m *Model) OpenConfig() {
	m.IsConfigOpen = true
	m.ConfigCursor = 0

	// run shell script called "ollama list" and get the output
	AvailableModels = make(map[string]ModelInfo)
	for model := range ClosedSourceModels {
		AvailableModels[model] = ClosedSourceModels[model]
	}
	ollamas, ollamaErr := runShellCommand("ollama list")
	if ollamaErr == nil {
		for _, ollama := range strings.Split(ollamas, "\n") {
			if !strings.Contains(ollama, "NAME") && ollama != "" {
				modelName := strings.Split(ollama, "  ")[0]
				AvailableModels[strings.Split(modelName, ":")[0]] = ModelInfo{Provider: "ollama", ID: modelName}
			}
		}
	}
	ModelKeys = make([]string, 0, len(AvailableModels))
	for model := range AvailableModels {
		ModelKeys = append(ModelKeys, model)
	}
}

// CloseConfig closes the config page
func (m *Model) CloseConfig() {
	m.IsConfigOpen = false
	m.ConfigCursor = 0
}

// HandleConfigNavigation handles up/down navigation in config page
func (m *Model) HandleConfigNavigation(direction string) {
	if !m.IsConfigOpen {
		return
	}

	if direction == tea.KeyUp.String() {
		if m.ConfigCursor > 0 {
			m.ConfigCursor--
		}
	} else if direction == tea.KeyDown.String() {
		if m.ConfigCursor < len(AvailableModels)-1 {
			m.ConfigCursor++
		}
	}
}

// SelectModel selects the current model in config
func (m *Model) SelectModel() {
	if !m.IsConfigOpen {
		return
	}

	selectedModel := ModelKeys[m.ConfigCursor]

	modelInfo, exists := AvailableModels[selectedModel]
	if !exists {
		m.AddSystemMessage("Error: Invalid model selected")
		m.CloseConfig()
		return
	}

	// Update the agent's model
	err := m.agent.SetModel(
		modelInfo.Provider,
		modelInfo.ID,
		modelInfo.Provider == "ollama",
	)
	if err != nil {
		m.AddSystemMessage("Error switching model: " + err.Error())
		m.CloseConfig()
		return
	}
	m.AddSystemMessage("Switched to model: " + selectedModel)
	m.CloseConfig()
}
