package model

// Message represents a chat message with a sender and text.
type Message struct {
	Sender  string // "user", "llm", or "system"
	Content string
}

// Model holds the state for the CLI.
type Model struct {
	Input    string
	Messages []Message
}

// NewModel initializes the CLI model with a welcome message.
func NewModel() Model {
	return Model{
		Input: "",
		Messages: []Message{
			{Sender: "system", Content: "👹 Welcome to Menace CLI!"},
		},
	}
}
