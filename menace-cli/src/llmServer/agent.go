package llmServer

import (
	"context"
	"fmt"
	"sync"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// Agent represents the main agent that processes LLM responses and executes commands.
//
// Does not include System messages
type Agent struct {
	llm      llms.Model
	mu       sync.Mutex
	shell    string // typically in the form "windows/CMD", "linux/bash", "darwin/bash" etc
	messages []llms.MessageContent
	ctx      context.Context
}

// NewAgent creates a new agent instance
//
// Returns: Agent, error
func NewAgent(apiKey string) (*Agent, error) {

	llm, err := openai.New(
		openai.WithToken(apiKey),
		openai.WithModel("o4-mini-2025-04-16"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %v", err)
	}

	ctx := context.Background()

	return &Agent{
		llm:   llm,
		shell: ModelFactory{}.DetectShell(),
		messages: []llms.MessageContent{
			{
				Role:  llms.ChatMessageTypeSystem,
				Parts: []llms.ContentPart{llms.TextContent{Text: getSystemPrompt(ModelFactory{}.DetectShell())}},
			},
		},
		ctx: ctx,
	}, nil
}

// SendMessage sends a message to the LLM and returns the response
//
// Sends a message to the LLM, returns a response.
// Parses the response, and returns a CommandSuggestion if found
//
// Does not interact with UI model.Messages at all.
// Returns: response, commandSuggestion, error
func (a *Agent) SendMessage(ctx context.Context, input string) (string, *CommandSuggestion, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Add user message to history
	a.messages = append(a.messages, llms.MessageContent{
		Role:  llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{llms.TextContent{Text: input}},
	})

	// Get response from LLM
	response, err := a.llm.GenerateContent(a.ctx, a.messages, llms.WithTemperature(1))
	if err != nil {
		return "", nil, fmt.Errorf("failed to get response from LLM: %v", err)
	}

	// Extract the response text
	var responseText string
	if len(response.Choices) > 0 {
		responseText = response.Choices[0].Content
	}

	// Parse for command suggestion
	cmdSuggestion := parseCommandSuggestion(responseText)

	// Add assistant's response to history
	a.messages = append(a.messages, llms.MessageContent{
		Role:  llms.ChatMessageTypeAI,
		Parts: []llms.ContentPart{llms.TextContent{Text: responseText}},
	})

	return responseText, cmdSuggestion, nil
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
