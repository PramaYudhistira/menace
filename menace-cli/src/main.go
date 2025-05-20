package main

import (
	"fmt"
	"menace-go/llmServer"
	"menace-go/ui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

func main() {
	// Get API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY environment variable not set")
		os.Exit(1)
	}

	// Initialize agent with langchaingo
	agent, err := llmServer.NewAgent(apiKey)
	if err != nil {
		fmt.Printf("Error initializing agent: %v\n", err)
		os.Exit(1)
	}

	// Initialize UI with the agent
	zone.NewGlobal()
	p := tea.NewProgram(
		ui.NewModel(agent), // Pass the agent to the UI
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running Menace CLI:", err)
		os.Exit(1)
	}
}
