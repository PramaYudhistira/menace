# Menace CLI - API Documentation

## Overview

This document provides comprehensive API documentation for Menace CLI's Go packages and their interfaces.

## Package Structure

```
src/
├── main.go                 # Entry point and initialization
├── llmServer/              # LLM integration and agent logic
│   ├── agent.go           # Core agent implementation
│   ├── command_suggestion.go # Command parsing and execution
│   ├── github_interface.go   # GitHub integration
│   ├── llmFactory.go      # LLM provider factory
│   ├── llmService.go      # LLM service interface
│   └── system_prompt.go   # System prompt management
└── ui/                    # Terminal user interface
    ├── executor.go        # Command execution
    ├── message.go         # Message handling
    ├── model.go          # UI state management
    ├── model_selection.go # Model selection interface
    ├── thinking.go       # Loading animations
    ├── ui_utils.go       # UI utilities
    ├── update.go         # State updates
    └── view.go           # UI rendering
```

## LLM Server Package (`llmServer`)

### Agent

The `Agent` struct is the core component that manages LLM interactions and maintains conversation state.

#### Type Definition

```go
type Agent struct {
    llm      llms.Model        // The underlying LLM implementation
    mu       sync.Mutex        // Thread safety mutex
    shell    string           // Shell type (e.g., "darwin/bash", "windows/CMD")
    messages []llms.MessageContent // Conversation history
    ctx      context.Context  // Request context
    provider string           // LLM provider ("openai", "anthropic", "ollama")
    Model    string           // Model name
    isOpenSource bool         // Whether using open-source model
}
```

#### Constructor

```go
func NewAgent(apiKey string) (*Agent, error)
```

Creates a new Agent instance with OpenAI as the default provider using the `o4-mini-2025-04-16` model.

**Parameters:**
- `apiKey`: OpenAI API key

**Returns:**
- `*Agent`: Configured agent instance
- `error`: Initialization error, if any

**Example:**
```go
agent, err := llmServer.NewAgent("sk-your-openai-key")
if err != nil {
    log.Fatal(err)
}
```

#### Methods

##### SendMessage

```go
func (a *Agent) SendMessage(ctx context.Context, input string) (string, *CommandSuggestion, error)
```

Sends a message to the LLM and processes the response for command suggestions.

**Parameters:**
- `ctx`: Request context
- `input`: User input message

**Returns:**
- `string`: LLM response text
- `*CommandSuggestion`: Parsed command suggestion, if any
- `error`: Processing error, if any

##### ClearHistory

```go
func (a *Agent) ClearHistory()
```

Clears conversation history while preserving the system prompt.

##### SetModel

```go
func (a *Agent) SetModel(provider string, model string, openSource bool) error
```

Changes the LLM provider and model.

**Supported Providers:**
- `"openai"` - OpenAI GPT models
- `"anthropic"` - Anthropic Claude models
- `"ollama"` - Local Ollama models

**Parameters:**
- `provider`: LLM provider name
- `model`: Model identifier
- `openSource`: Whether the model is open source

##### AddToMessageChain

```go
func (a *Agent) AddToMessageChain(new_message string, role llms.ChatMessageType)
```

Adds a message to the conversation history with a specific role.

**Parameters:**
- `new_message`: Message content
- `role`: Message role (`llms.ChatMessageTypeSystem`, `llms.ChatMessageTypeHuman`, `llms.ChatMessageTypeAI`)

### Command Suggestion

The command suggestion system parses LLM responses for executable commands.

#### Type Definition

```go
type CommandSuggestion struct {
    Reason                  string // Explanation of why the command is needed
    Command                 string // The actual command to execute
    AwaitingCommandApproval bool   // Whether human approval is required
}
```

#### Parsing Function

```go
func parseCommandSuggestion(response string) *CommandSuggestion
```

Parses LLM response text for command suggestion blocks in the format:

```
[COMMAND_SUGGESTION]
Reason: <explanation>
Command: <command>
AwaitingCommandApproval: true/false
[/COMMAND_SUGGESTION]
```

### System Prompt Management

#### Functions

```go
func getSystemPrompt(shell string) string
```

Generates the system prompt with shell-specific instructions and current working directory information.

```go
func getSystemPromptNew(shell string) string
```

Alternative system prompt format using JSON response structure.

## UI Package (`ui`)

### Model

The `Model` struct manages the entire UI state and user interactions.

#### Type Definition

```go
type Model struct {
    Input    string    // Current user input
    Messages []Message // Chat history
    agent    *llmServer.Agent // LLM agent instance
    Width    int       // Terminal width
    Height   int       // Terminal height
    Scroll   int       // Chat scroll position
    
    // Cursor management
    CursorX     int // Horizontal cursor position
    CursorY     int // Vertical cursor position (line number)
    WindowStart int // Horizontal scroll offset
    
    // Text selection
    SelectionStartX int
    SelectionStartY int
    SelectionEndX   int
    SelectionEndY   int
    IsHighlighting  bool
    
    // UI state
    IsThinking   bool // Loading animation state
    ThinkingDots int  // Animation frame counter
    
    // Configuration
    IsConfigOpen bool // Model selection UI
    ConfigCursor int  // Selected model index
    
    // Command execution
    PendingCommand          *CommandSuggestionMsg
    PendingFunctionCall     *FunctionCallMsg
    AwaitingCommandApproval bool
}
```

#### Constructor

```go
func NewModel(agent *llmServer.Agent) *Model
```

Creates a new UI model with the provided agent.

#### Key Methods

##### Input Management

```go
func (m *Model) InsertCharacter(character string)
func (m *Model) InsertNewLine()
func (m *Model) HandleBackSpace()
func (m *Model) HandleDelete()
```

Text input and editing operations with Unicode support.

##### Cursor Navigation

```go
func (m *Model) HandleHorizontalCursorMovement(direction string)
func (m *Model) HandleVerticalCursorMovement(direction string)
```

Cursor movement with line-aware navigation.

##### Text Selection

```go
func (m *Model) SelectAll()
func (m *Model) GetSelectedText() string
func (m *Model) CutSelectedText()
```

Text selection and clipboard operations.

##### Chat Management

```go
func (m *Model) AddUserMessage(message string)
func (m *Model) AddAgentMessage(message string)
func (m *Model) AddSystemMessage(message string)
```

Message history management.

##### Window Management

```go
func (m *Model) ResizeWindow(msg tea.Msg)
func (m *Model) HandleScroll(direction int)
func (m *Model) UpdateWindowStart(maxInputW int)
```

Terminal resize handling and scrolling.

##### AI Integration

```go
func (m *Model) AIChat(msg string) tea.Msg
```

Sends message to the AI agent and returns appropriate response message.

### Message Types

The UI uses various message types for the Bubble Tea event system:

```go
type LLMResponseMsg struct {
    Content string
}

type SystemMessage struct {
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
```

## Function Calls

Menace CLI supports several built-in function calls for file operations and GitHub integration:

### File Operations

#### ReadFileWithLineNumbers

Reads a file and returns content with line numbers.

```json
{
    "name": "ReadFileWithLineNumbers",
    "args": {
        "path": "example.py"
    }
}
```

#### CreateAndApplyDiffs

Applies line-level changes to files.

```json
{
    "name": "CreateAndApplyDiffs",
    "args": {
        "path": "example.txt",
        "diffs": [
            {
                "Type": 0,        // 0=Add, 1=Delete, 2=Modify
                "LineIndex": 2,   // 1-based line number
                "OldContent": "", // Previous content (Delete/Modify only)
                "NewContent": "new line" // New content (Add/Modify only)
            }
        ]
    }
}
```

### GitHub Integration

#### createPullRequest

Creates a GitHub pull request.

```json
{
    "name": "createPullRequest",
    "args": {
        "branch_name": "feature-branch",
        "title": "Add new feature",
        "summary": "Detailed description of changes"
    }
}
```

## Error Handling

All public functions return appropriate error types:

- **Agent initialization errors**: Invalid API keys, network issues
- **LLM communication errors**: API rate limits, invalid responses
- **File operation errors**: Permission issues, file not found
- **Command execution errors**: Invalid commands, execution failures

## Thread Safety

The Agent struct uses mutex locking to ensure thread-safe access to the conversation history and LLM client. The UI runs on a single thread using the Bubble Tea framework.

## Configuration

### Environment Variables

- `OPENAI_API_KEY`: OpenAI API key (required for OpenAI provider)
- `ANTHROPIC_API_KEY`: Anthropic API key (required for Anthropic provider)

### Shell Detection

The system automatically detects the shell environment:
- **macOS/Linux**: `darwin/bash`, `linux/bash`
- **Windows**: `windows/CMD`

This information is used to generate appropriate system prompts and command suggestions.