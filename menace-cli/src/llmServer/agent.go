package llmServer

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
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
	provider string
	Model    string
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
		llm:      llm,
		shell:    ModelFactory{}.DetectShell(),
		Model:    "o4-mini-2025-04-16",
		provider: "openai",
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

func (a *Agent) SetModel(provider string, model string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.provider = provider
	a.Model = model

	switch provider {
	case "anthropic":
		llm, err := anthropic.New(
			anthropic.WithToken(os.Getenv("ANTHROPIC_API_KEY")),
			anthropic.WithModel(model),
		)
		if err != nil {
			return fmt.Errorf("failed to create Anthropic client with model %s: %v", model, err)
		}
		a.llm = llm
		return nil
	case "openai":
		llm, err := openai.New(
			openai.WithToken(os.Getenv("OPENAI_API_KEY")),
			openai.WithModel(model),
		)
		if err != nil {
			return fmt.Errorf("failed to create OpenAI client with model %s: %v", model, err)
		}
		a.llm = llm
		return nil
	default:
		return fmt.Errorf("unknown provider: %s", provider)
	}
}

// Returns: System prompt
func getSystemPrompt(shell string) string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "unknown directory"
	}
	return fmt.Sprintf(`You are operating as and within the Menace CLI. You must be safe, precise and helpful.
	Menace-CLI is a lightweight CLI tool that uses large language models to provide intelligent terminal assistance.
	You have access to the local file system and can execute commands in %s. The user will either tell you to execute a command, 
	or you can decide if its best to execute a command.

	You can also edit files, write code, etc. if it is required to finish the task.
	When performing tasks, always ensure that every step is ahieves exactly what the user requests.
	Do not write code, or run commands unless you are certain it is necessary.
	When uncertain, ask clarifying questions.

	ONLY EXECUTE COMMANDS WHICH WORK ON  %s!
	**You are always operating in the current working directory: %s.**
	- Do not use wildcards, recursive, or bare-format flags.

	To execute a command, you must follow this format:
	[COMMAND_SUGGESTION]
	Reason: <explain why this command is needed>
	Command: your_command_here
	[/COMMAND_SUGGESTION]

	For example:
	[COMMAND_SUGGESTION]
	Reason: To list all files in the current directory
	Command: ls
	[/COMMAND_SUGGESTION]

	You can edit files, write code, and search for files and functions using commands.

	If you need to edit a file:
	- First, inform the user and ask for their approval.
	- Only proceed with the edit if the user agrees.
	- If you are unable to make the edit, ask the user to do it manually.

	Do NOT suggest opening files in editors like Notepad, nano, vim, or any GUI or interactive editor.
	If you need to read, write, or modify a file, use shell commands to do so directly (e.g., using echo, type, copy, move, powershell, cat, sed, etc.).
	Never suggest launching an editor; always use direct shell commands for file operations.

	If your command writes to a file (for example, using >, >>, echo, or similar redirection),
	you must clearly explain to the user what was written and to which file. Do not expect output 
	to appear in the terminal for such commands. Always describe the effect of the command in your next 
	response if there is no terminal output. If your next response is a command suggestion, include this 
	explanation in the Reason field.

	You should only run commands one at a time. If you need to run multiple commands, just give 1 command, then wait 
	for the user to respond with the output, and then run that command and so on. All until you have achieved the goal.
	When you have a proposed command, the user might not respond with the output, and instead either say "no", or 
	will give you feedback and direction on what to do next.

	You should respond as if you are part of this real application, not a fictional tool.
	`, shell, shell, cwd)
}
