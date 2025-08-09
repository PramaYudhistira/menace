# API Reference

This document provides comprehensive documentation for the internal Go APIs in Menace CLI.

## Package Overview

### `llmServer` Package

The `llmServer` package provides the core LLM integration and agent functionality.

#### Types

##### `Agent`

The main agent that processes LLM responses and executes commands.

```go
type Agent struct {
    llm      llms.Model           // LLM implementation (OpenAI, Anthropic, Ollama)
    mu       sync.Mutex          // Thread safety for concurrent access
    shell    string              // Platform shell (e.g., "windows/CMD", "linux/bash")
    messages []llms.MessageContent // Conversation history
    ctx      context.Context     // Request context
    provider string              // Current LLM provider
    Model    string              // Current model name
    isOpenSource bool            // Whether using open-source model
}
```

##### `CommandSuggestion`

Represents a parsed command suggestion from the LLM.

```go
type CommandSuggestion struct {
    Reason                   string // Explanation for the command
    Command                  string // Shell command to execute
    AwaitingCommandApproval bool   // Whether user approval is required
}
```

#### Functions

##### `NewAgent(apiKey string) (*Agent, error)`

Creates a new agent instance with OpenAI as the default provider.

**Parameters:**
- `apiKey`: OpenAI API key

**Returns:**
- `*Agent`: Configured agent instance
- `error`: Error if initialization fails

**Example:**
```go
agent, err := llmServer.NewAgent("sk-...")
if err != nil {
    log.Fatal(err)
}
```

##### `(*Agent) SendMessage(ctx context.Context, input string) (string, *CommandSuggestion, error)`

Sends a message to the LLM and returns the response with any parsed command suggestions.

**Parameters:**
- `ctx`: Request context
- `input`: User input message

**Returns:**
- `string`: LLM response text
- `*CommandSuggestion`: Parsed command suggestion (nil if none)
- `error`: Error if request fails

**Example:**
```go
response, cmd, err := agent.SendMessage(ctx, "list files in current directory")
if err != nil {
    return err
}
if cmd != nil && cmd.Command == "ls" {
    // Execute the suggested command
}
```

##### `(*Agent) ClearHistory()`

Clears the conversation history, keeping only the system message.

##### `(*Agent) SetModel(provider string, model string, openSource bool) error`

Switches the LLM provider and model at runtime.

**Parameters:**
- `provider`: Provider name ("openai", "anthropic", "ollama")
- `model`: Model identifier
- `openSource`: Whether the model is open-source

**Supported Providers:**
- **OpenAI**: Requires `OPENAI_API_KEY`
- **Anthropic**: Requires `ANTHROPIC_API_KEY`  
- **Ollama**: Requires local Ollama server

**Example:**
```go
err := agent.SetModel("anthropic", "claude-3-sonnet", false)
if err != nil {
    log.Printf("Failed to switch model: %v", err)
}
```

##### `(*Agent) AddToMessageChain(new_message string, role llms.ChatMessageType)`

Adds a message directly to the conversation chain.

**Parameters:**
- `new_message`: Message content
- `role`: Message role (System, Human, AI)

### `ui` Package

The `ui` package provides the terminal user interface using Bubble Tea.

#### Types

##### `Model`

Main UI model managing application state.

```go
type Model struct {
    agent           *llmServer.Agent
    messages        []Message
    input           string
    cursor          int
    viewport        viewport.Model
    ready           bool
    windowWidth     int
    windowHeight    int
    // ... additional fields
}
```

##### `Message`

Represents a chat message in the UI.

```go
type Message struct {
    Content   string
    Type      MessageType
    Timestamp time.Time
}
```

##### `MessageType`

Enumeration of message types:

```go
const (
    UserMessage MessageType = iota
    AssistantMessage
    SystemMessage
    CommandOutput
    ErrorMessage
)
```

#### Functions

##### `NewModel(agent *llmServer.Agent) Model`

Creates a new UI model with the given agent.

**Parameters:**
- `agent`: Configured LLM agent

**Returns:**
- `Model`: Initialized UI model

##### `(m Model) Init() tea.Cmd`

Initializes the model and returns initial commands.

##### `(m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)`

Handles UI events and updates the model state.

**Message Types Handled:**
- `tea.KeyMsg`: Keyboard input
- `tea.MouseMsg`: Mouse interactions
- `tea.WindowSizeMsg`: Window resize events
- Custom commands for async operations

##### `(m Model) View() string`

Renders the current UI state as a string.

**Layout:**
- Sidebar with model information and controls
- Chat area with conversation history
- Input area for user messages

#### Key Bindings

- `Ctrl+C`: Quit application
- `Enter`: Send message
- `Ctrl+A`: Select all text
- `Ctrl+C` (with selection): Copy to clipboard
- `Ctrl+V`: Paste from clipboard
- `Ctrl+X`: Cut selected text
- `Ctrl+Z`: Undo last action
- `Ctrl+Y`: Redo last action
- Mouse scrolling: Navigate chat history

## Function Call System

Menace CLI supports structured function calls for file operations and GitHub integration.

### Available Functions

#### `ReadFileWithLineNumbers`

Reads a file and returns content with line numbers.

**Parameters:**
```json
{
    "path": "file/path"
}
```

#### `CreateAndApplyDiffs`

Applies line-level changes to a file using diff operations.

**Parameters:**
```json
{
    "path": "file/path",
    "diffs": [
        {
            "Type": 0,          // 0=Add, 1=Delete, 2=Modify
            "LineIndex": 1,     // 1-based line number
            "OldContent": "",   // Previous content (for Delete/Modify)
            "NewContent": "new" // New content (for Add/Modify)
        }
    ]
}
```

#### `createPullRequest`

Creates a GitHub pull request.

**Parameters:**
```json
{
    "branch_name": "feature-branch",
    "title": "PR Title",
    "summary": "Detailed description"
}
```

## Error Handling

All APIs follow Go conventions for error handling:

- Functions return `error` as the last return value
- Errors are propagated up the call stack
- Critical errors cause program termination with `os.Exit(1)`

## Thread Safety

- The `Agent` type uses `sync.Mutex` for thread-safe access
- UI operations run on the main thread via Bubble Tea's command system
- File operations are executed synchronously within command handlers

## Configuration

### Environment Variables

- `OPENAI_API_KEY`: Required for OpenAI provider
- `ANTHROPIC_API_KEY`: Required for Anthropic provider
- Ollama requires local server configuration

### Shell Detection

The system automatically detects the current shell:
- Windows: `windows/CMD`
- macOS: `darwin/bash` or `darwin/zsh`
- Linux: `linux/bash`

This information is included in the system prompt for context-aware responses.