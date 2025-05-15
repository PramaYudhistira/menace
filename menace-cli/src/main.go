package main

import (
	"fmt"
	"menace-go/llmServer"
	"menace-go/ui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize LLM service
	llm := llmServer.GetInstance()
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Warning: OPENAI_API_KEY environment variable not set")
	} else {
		llm.Configure(apiKey)
	}

	// Initialize agent
	agent := llmServer.NewAgent(llm)

	p := tea.NewProgram(
		ui.NewModel(agent),
		tea.WithAltScreen()) // alternate screen
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running Menace CLI:", err)
		os.Exit(1)
	}
}
