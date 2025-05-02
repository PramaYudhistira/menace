package ui

import "github.com/charmbracelet/lipgloss"

// Styles for Menace CLI UI
var (
	SidebarStyle = lipgloss.NewStyle().
			Width(20).
			PaddingRight(1).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("8"))

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
)
