package ui

type Message struct {
	Sender  string // "user", "llm", or "system"
	Content string
}
