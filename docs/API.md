# Menace CLI API Documentation

## Overview

Menace CLI is a cross-platform terminal-based AI assistant built with Go and Node.js. It provides an interactive Text User Interface (TUI) powered by Bubble Tea and integrates with various Language Model APIs including OpenAI, Anthropic, and Ollama.

## Architecture

### Core Components

#### 1. Agent (`llmServer/agent.go`)

The `Agent` struct is the core component that manages LLM interactions and maintains conversation state.

```go
type Agent struct {
    llm          llms.Model           // The underlying LLM client
    mu           sync.Mutex           // Thread safety mutex
    shell        string               // Shell environment info (e.g., "linux/bash")
    messages     []llms.MessageContent // Conversation history
    ctx          context.Context      // Request context
    provider     string               // LLM provider ("openai", "anthropic", "ollama")
    Model        string               // Model name (e.g., "o4-mini-2025-04-16")
    isOpenSource bool                 // Flag for open-source models
}
```

**Key Methods:**

- `NewAgent(apiKey string) (*Agent, error)` - Creates a new agent instance with OpenAI as default provider
- `SendMessage(ctx context.Context, input string) (string, *CommandSuggestion, error)` - Sends user input to LLM and returns response with optional command suggestions
- `SetModel(provider, model string, openSource bool) error` - Switches between different LLM providers and models
- `ClearHistory()` - Resets conversation history while keeping system prompt
- `AddToMessageChain(message string, role llms.ChatMessageType)` - Adds messages to conversation history

#### 2. UI Model (`ui/model.go`)

The `Model` struct manages the terminal user interface state and user interactions.

```go
type Model struct {
    Input                   string             // Current user input
    Messages               []Message          // Chat history displayed to user
    agent                  *llmServer.Agent   // Reference to the LLM agent
    Width, Height          int                // Terminal dimensions
    Scroll                 int                // Chat scroll position
    CursorX, CursorY       int                // Input cursor position
    WindowStart            int                // Horizontal scrolling offset
    IsThinking             bool               // Loading state indicator
    IsConfigOpen           bool               // Configuration panel state
    PendingCommand         *CommandSuggestionMsg // Command awaiting user approval
    PendingFunctionCall    *FunctionCallMsg   // Function awaiting execution
    AwaitingCommandApproval bool              // Approval workflow state
}
```

**Key Methods:**

- `NewModel(agent *llmServer.Agent) *Model` - Creates new UI model with agent reference
- `InsertCharacter(character string)` - Handles text input with Unicode support
- `HandleBackSpace()` - Manages text deletion and selection
- `HandleScroll(direction int)` - Manages chat history scrolling
- `AddUserMessage(message string)` - Adds user message to chat history
- `AddAgentMessage(message string)` - Adds AI response to chat history
- `AIChat(msg string) tea.Msg` - Initiates AI conversation and returns Bubble Tea message

### 3. Command Processing

#### Command Suggestions

Commands are suggested by the AI using structured markup:

```
[COMMAND_SUGGESTION]
Reason: <explanation>
Command: <shell_command>
AwaitingCommandApproval: true/false
[/COMMAND_SUGGESTION]
```

#### Function Calls

File operations and GitHub actions use function call syntax:

```
[FUNCTION_CALL]
Reason: <explanation>
AwaitingCommandApproval: true/false
Payload:
{
    "name": "ReadFileWithLineNumbers",
    "args": {
        "path": "example.py"
    }
}
[/FUNCTION_CALL]
```

**Available Functions:**

1. `ReadFileWithLineNumbers` - Read file with line numbers
2. `CreateAndApplyDiffs` - Apply line-level changes to files
3. `createPullRequest` - Create GitHub pull requests

### 4. LLM Factory (`llmServer/llmFactory.go`)

Manages different LLM providers and model configurations.

**Supported Providers:**
- **OpenAI**: GPT-4, GPT-3.5, o4-mini models
- **Anthropic**: Claude models (requires ANTHROPIC_API_KEY)
- **Ollama**: Local open-source models

### 5. System Prompts (`llmServer/system_prompt.go`)

Dynamic system prompts that adapt to the user's shell environment and working directory.

**Features:**
- Shell-aware command suggestions
- Current working directory context
- File operation capabilities
- GitHub integration instructions
- Safety guidelines and limitations

## Message Types

### Bubble Tea Messages

The application uses Bubble Tea's message-passing architecture:

```go
// Core message types
type LLMResponseMsg struct {
    Content string
}

type CommandSuggestionMsg struct {
    Command string
    Reason  string
}

type FunctionCallMsg struct {
    Name string
    Args map[string]interface{}
}

type SystemMessage struct {
    Content string
}
```

## File Operations

### Diff System

File modifications use a structured diff system:

```go
type DiffType int

const (
    Add    DiffType = 0  // Add new line
    Delete DiffType = 1  // Remove existing line  
    Modify DiffType = 2  // Change existing line
)

type LineDiff struct {
    Type       DiffType
    LineIndex  int    // 1-based line number
    OldContent string // Previous content (Delete/Modify)
    NewContent string // New content (Add/Modify)
}
```

## GitHub Integration

### Pull Request Creation

```go
// Function signature for PR creation
func CreatePullRequest(branchName, title, summary string) error
```

**Required Environment:**
- Configured Git repository
- GitHub authentication (via gh CLI or tokens)
- Proper branch setup

## Build System

### Cross-Platform Compilation

The Node.js build scripts support multiple architectures:

```javascript
const buildTargets = [
    { os: "linux", arch: "amd64", output: "menace-go-linux" },
    { os: "darwin", arch: "amd64", output: "menace-go-darwin" }, 
    { os: "darwin", arch: "arm64", output: "menace-go-darwin-arm64" },
    { os: "windows", arch: "amd64", output: "menace-go-win.exe" }
];
```

### Runtime Binary Selection

The Node.js wrapper (`bin/menace.js`) automatically selects the correct binary based on platform:

```javascript
switch (process.platform) {
    case "win32": binName = "menace-go-win.exe"; break;
    case "darwin": binName = process.arch === "arm64" ? 
        "menace-go-darwin-arm64" : "menace-go-darwin"; break;
    case "linux": binName = "menace-go-linux"; break;
}
```

## Configuration

### Environment Variables

- `OPENAI_API_KEY` - Required for OpenAI models
- `ANTHROPIC_API_KEY` - Required for Anthropic models  
- Additional environment variables are detected automatically (shell, working directory)

### Model Selection

Users can switch between models at runtime through the configuration interface accessible via keyboard shortcuts in the TUI.

## Error Handling

The system implements comprehensive error handling:

- **LLM API Errors**: Network timeouts, authentication failures, rate limits
- **File System Errors**: Permission denied, file not found, disk space
- **Command Execution Errors**: Invalid commands, shell execution failures
- **UI State Errors**: Terminal resize handling, input validation

## Security Considerations

- Commands require explicit user approval for destructive operations
- File modifications are logged and reversible through diffs
- GitHub operations require proper authentication
- Shell command execution is limited to user's current permissions
- No arbitrary code execution without user consent

## Performance

- **Streaming Support**: Real-time response rendering
- **Efficient Scrolling**: Optimized chat history navigation
- **Memory Management**: Conversation history limits to prevent memory leaks
- **Concurrent Safety**: Thread-safe agent operations with mutex protection

## Extensibility

The modular architecture supports:

- **New LLM Providers**: Add providers by implementing the `llms.Model` interface
- **Custom Functions**: Extend function call system with new capabilities
- **UI Themes**: Customize appearance through Lip Gloss styling
- **Shell Support**: Add support for additional shell environments