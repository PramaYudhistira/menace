package llmServer

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// Agent represents the main agent that processes LLM responses and executes commands
// Does not include System messages
type Agent struct {
	llm      llms.Model
	mu       sync.Mutex
	shell    string // typically in the form "windows/CMD", "linux/bash", "darwin/bash" etc
	messages []llms.MessageContent
}

// NewAgent creates a new agent instance
func NewAgent(apiKey string) (*Agent, error) {

	// Right now, o3-2025-04-16 cannot be used with the openai package...
	// Perhaps its time for OSS contribution again?
	// it wont work with o3-2025-04-16 since we can't set temperature to 0 because openai package is broken
	llm, err := openai.New(
		openai.WithToken(apiKey),
		openai.WithModel("gpt-4.1-2025-04-14"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %v", err)
	}

	return &Agent{
		llm:   llm,
		shell: ModelFactory{}.DetectShell(),
		messages: []llms.MessageContent{
			{
				Role:  llms.ChatMessageTypeSystem,
				Parts: []llms.ContentPart{llms.TextContent{Text: getSystemPrompt(ModelFactory{}.DetectShell())}},
			},
		},
	}, nil
}

// SendMessage sends a message to the LLM and returns the response
func (a *Agent) SendMessage(ctx context.Context, input string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Add user message to history
	a.messages = append(a.messages, llms.MessageContent{
		Role:  llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{llms.TextContent{Text: input}},
	})

	// Get response from LLM
	response, err := a.llm.GenerateContent(ctx, a.messages)
	if err != nil {
		return "", fmt.Errorf("failed to get response from LLM: %v", err)
	}

	// Extract the response text
	var responseText string
	if len(response.Choices) > 0 {
		responseText = response.Choices[0].Content
	}

	// Add assistant's response to history
	a.messages = append(a.messages, llms.MessageContent{
		Role:  llms.ChatMessageTypeAI,
		Parts: []llms.ContentPart{llms.TextContent{Text: responseText}},
	})

	return responseText, nil
}

// ClearHistory clears the conversation history
//
// Only persistent in the backend
func (a *Agent) ClearHistory() {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Keep only the system message
	a.messages = []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{Text: getSystemPrompt(a.shell)}},
		},
	}
}

// getSystemPrompt returns the system prompt for the agent
func getSystemPrompt(shell string) string {
	shellType := strings.Split(shell, "/")[1]
	return fmt.Sprintf(`You are operating as and within the Menace CLI. You must be safe, precise and helpful.
	Menace-CLI is a Go-based CLI tool that uses large language models to provide intelligent terminal assistance.
	You have access to the local file system and can execute commands in %s.

	Before executing any command, explain your intent and reasoning.
	Execute commands sequentially - one at a time.
	After each command, you'll receive its output (especially for file operations like ls).

	When you need to execute a command, format it like this:
	`+"```"+shellType+"\n"+`your_command_here
	`+"```\n"+`

	For example, to list files:
	`+"```"+shellType+"\n"+`ls
	`+"```\n"+`

	You should respond as if you are part of this real application, not a fictional tool.
	`, shell)
}

// Command represents a parsed command from the LLM response
type Command struct {
	Type    string // "shell", "file", etc.
	Content string // the command to execute
	Error   error
}
