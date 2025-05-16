package llmServer

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// Agent represents the main agent that processes LLM responses and executes commands
type Agent struct {
	llmService *LLMService
	mu         sync.Mutex
	shell      string // typically in the form "windows/CMD", "linux/bash", "darwin/bash" etc
}

// NewAgent creates a new agent instance
func NewAgent(llmService *LLMService) *Agent {
	return &Agent{
		llmService: llmService,
		shell:      ModelFactory{}.DetectShell(),
	}
}

// Command represents a parsed command from the LLM response
type Command struct {
	Type    string // "shell", "file", etc.
	Content string // the command to execute
	Error   error
}

// parseResponse attempts to extract commands from the LLM response
// note, currently unused...
func (a *Agent) parseResponse(response string) ([]Command, error) {
	commands := []Command{}

	// Get shell type from the detected shell
	// since its in the form "windows/CMD", "linux/bash", "darwin/bash" etc, we split into 2 parts
	parts := strings.Split(a.shell, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid shell format: %s", a.shell)
	}
	shellType := parts[1]

	// Look for commands in the response
	// same format as instructed in the system prompt
	pattern := "```" + shellType
	if strings.Contains(response, pattern) {
		parts := strings.Split(response, pattern)
		for _, part := range parts[1:] {
			endIdx := strings.Index(part, "```")
			if endIdx == -1 {
				continue
			}

			cmd := strings.TrimSpace(part[:endIdx])
			if cmd != "" {
				commands = append(commands, Command{
					Type:    "shell",
					Content: cmd,
				})
			}
		}
	}

	return commands, nil
}

// parseCommand extracts the command found in the LLM response
// Returns nil if no command is found
func ParseCommand(response string) *Command {
	// Look for the first command block
	pattern := "```" + shellType
	if !strings.Contains(response, pattern) {
		return nil
	}

	// Split on the pattern and take the first command block
	parts := strings.Split(response, pattern)
	if len(parts) < 2 {
		return nil
	}

	// Find the end of the command block
	endIdx := strings.Index(parts[1], "```")
	if endIdx == -1 {
		return &Command{
			Error: fmt.Errorf("malformed command block: missing closing ```"),
		}
	}

	// Extract and trim the command
	cmd := strings.TrimSpace(parts[1][:endIdx])
	if cmd == "" {
		return nil
	}

	return &Command{
		Type:    "shell",
		Content: cmd,
	}
}

// executeCommand runs a single command and returns the result
func (a *Agent) ExecuteCommand(cmd Command) (string, error) {
	if cmd.Type != "shell" {
		return "", fmt.Errorf("unsupported command type: %s", cmd.Type)
	}

	// Split shell into OS and shell type
	parts := strings.Split(a.shell, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid shell format: %s", a.shell)
	}

	osType, shellType := parts[0], parts[1]

	var execCmd *exec.Cmd
	if osType == "windows" {
		if shellType == "PowerShell" {
			execCmd = exec.Command("powershell", "-Command", cmd.Content)
		} else {
			// CMD
			execCmd = exec.Command("cmd", "/C", cmd.Content)
		}
	} else {
		// Unix-like systems
		execCmd = exec.Command(shellType, "-c", cmd.Content)
	}

	// Capture both stdout and stderr
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command execution failed: %v", err)
	}

	return string(output), nil
}

// Run is the main agent loop that processes user input and executes commands
// this main loop has multiple commands at the same time...
func (a *Agent) Run(input string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Get response from LLM
	response, err := a.llmService.SendMessage(input)
	if err != nil {
		return fmt.Errorf("failed to get LLM response: %v", err)
	}

	// Parse commands from response
	commands, err := a.parseResponse(response.Choices[0].Message.Content)
	if err != nil {
		return fmt.Errorf("failed to parse commands: %v", err)
	}

	// Execute each command
	for _, cmd := range commands {
		result, err := a.ExecuteCommand(cmd)
		if err != nil {
			// Log error but continue with other commands
			fmt.Printf("Error executing command: %v\n", err)
			continue
		}

		// Send result back to LLM for context
		contextMsg := fmt.Sprintf("Command executed: %s\nResult: %s", cmd.Content, result)
		_, err = a.llmService.SendMessage(contextMsg)
		if err != nil {
			fmt.Printf("Failed to send command result to LLM: %v\n", err)
		}
	}

	return nil
}
