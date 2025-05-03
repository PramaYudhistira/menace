package main

import (
	"fmt"
	"menace-go/ui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(
		ui.NewModel(),
		tea.WithAltScreen()) // alternate screen
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running Menace CLI:", err)
		os.Exit(1)
	}
}
