package main

import (
	"menace-go/llmServer"
)

func main() {
	// Get API key
	llmServer.PushToGitHub()
	// apiKey := os.Getenv("OPENAI_API_KEY")
	// if apiKey == "" {
	// 	fmt.Println("Error: OPENAI_API_KEY environment variable not set")
	// 	os.Exit(1)
	// }

	// // Initialize agent with langchaingo
	// agent, err := llmServer.NewAgent(apiKey)
	// if err != nil {
	// 	fmt.Printf("Error initializing agent: %v\n", err)
	// 	os.Exit(1)
	// }

	// // Initialize UI with the agent
	// zone.NewGlobal()
	// p := tea.NewProgram(
	// 	ui.NewModel(agent), // Pass the agent to the UI
	// 	tea.WithAltScreen(),
	// 	tea.WithMouseAllMotion(),
	// )

	// if _, err := p.Run(); err != nil {
	// 	fmt.Println("Error running Menace CLI:", err)
	// 	os.Exit(1)
	// }
}
