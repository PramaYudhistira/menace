# Menace CLI API Documentation

This document provides comprehensive API documentation for developers working with or extending Menace CLI.

## Table of Contents

1. [LLM Server Package](#llm-server-package)
2. [UI Package](#ui-package)
3. [Message Types](#message-types)
4. [Function Calling System](#function-calling-system)
5. [Integration APIs](#integration-apis)
6. [Extension Points](#extension-points)

## LLM Server Package

The `llmServer` package provides the core intelligence layer, managing LLM interactions and agent behavior.

### Agent

**Primary Interface**: The `Agent` struct serves as the main interface for LLM operations.

#### Type Definition

```go
type Agent struct {
    llm      llms.Model              // LangChain Go LLM interface
    mu       sync.Mutex              // Thread-safety mutex  
    shell    string                  // Detected shell environment
    messages []llms.MessageContent   // Conversation history
    ctx      context.Context         // Request context
    provider string                  // Current LLM provider name
    Model    string                  // Current model identifier
    isOpenSource bool               // Local vs. cloud model flag
}
```

#### Constructor

```go
func NewAgent(apiKey string) (*Agent, error)
```

**Parameters**:
- `apiKey`: OpenAI API key for authentication

**Returns**:
- `*Agent`: Initialized agent instance
- `error`: Initialization error (invalid API key, network issues)

**Example**:
```go
agent, err := llmServer.NewAgent("sk-...")
if err != nil {
    log.Fatal("Failed to create agent:", err)
}
```

#### Core Methods

##### ProcessQuery

```go
func (a *Agent) ProcessQuery(query string) (string, error)
```

**Purpose**: Process user input and generate LLM response.

**Parameters**:
- `query`: User's natural language input

**Returns**:
- `string`: LLM-generated response
- `error`: Processing error

**Features**:
- Automatic context management
- Response parsing for commands/function calls
- Thread-safe operation

##### SetModel

```go
func (a *Agent) SetModel(provider, modelName string) error
```

**Purpose**: Switch LLM provider and model at runtime.

**Parameters**:
- `provider`: Provider name ("openai", "anthropic", "ollama")
- `modelName`: Specific model identifier

**Supported Providers & Models**:

| Provider | Models |
|----------|--------|
| OpenAI | gpt-4, gpt-4-turbo, gpt-3.5-turbo |
| Anthropic | claude-3-opus, claude-3-sonnet, claude-3-haiku |
| Ollama | llama2, codellama, mistral, (local models) |

##### GetAvailableModels

```go
func (a *Agent) GetAvailableModels(provider string) ([]string, error)
```

**Purpose**: Retrieve list of available models for a provider.

**Parameters**:
- `provider`: Target provider name

**Returns**:
- `[]string`: List of model identifiers
- `error`: Query error

### Command Suggestion System

#### Type Definition

```go
type CommandSuggestion struct {
    Reason  string  // Explanation of why command is suggested
    Command string  // Shell command to execute
    AwaitingCommandApproval bool // Requires user confirmation
}
```

#### Parsing Protocol

**Tag Format**: Commands are extracted from LLM responses using structured tags:

```
[COMMAND_SUGGESTION]
reason: Install missing package
command: npm install express
[/COMMAND_SUGGESTION]
```

#### Parser Function

```go
func ParseCommandSuggestion(response string) (*CommandSuggestion, error)
```

**Purpose**: Extract command suggestions from LLM responses.

**Parameters**:
- `response`: Raw LLM response text

**Returns**:
- `*CommandSuggestion`: Parsed command data
- `error`: Parsing error

### System Prompt Management

#### GetSystemPrompt

```go
func GetSystemPrompt(shell string) string
```

**Purpose**: Generate context-aware system prompts for LLMs.

**Parameters**:
- `shell`: Detected shell environment ("bash", "powershell", "cmd", etc.)

**Features**:
- Shell-specific command syntax
- Function calling definitions
- Safety instructions
- Git integration guidance

**Example Output**:
```
You are an intelligent terminal assistant for bash on Linux.
Available functions:
- read_file(path): Read file contents
- write_file(path, content): Write to file
- create_pull_request(title, description): Create GitHub PR
...
```

### GitHub Integration

#### CreatePullRequest

```go
func CreatePullRequest(title, description, branchName string) error
```

**Purpose**: Automated GitHub pull request creation.

**Parameters**:
- `title`: PR title
- `description`: PR description/body
- `branchName`: Source branch name

**Requirements**:
- `GITHUB_TOKEN` environment variable
- Git repository with GitHub origin
- Uncommitted changes or new branch

## UI Package

The `ui` package implements the terminal user interface using the Bubble Tea framework.

### Model

**Core Structure**: The `Model` struct manages all UI state and user interactions.

#### Type Definition

```go
type Model struct {
    // Core application state
    Input    string           // Current user input
    Messages []Message        // Conversation history
    agent    *llmServer.Agent // LLM agent reference
    
    // Display properties
    Width, Height int         // Terminal dimensions
    Scroll       int         // Vertical scroll position
    
    // Text editing state
    CursorX, CursorY    int  // Cursor position
    WindowStart         int  // Horizontal scroll offset
    
    // Selection system
    SelectionStartX, SelectionStartY int
    SelectionEndX, SelectionEndY     int
    IsHighlighting                   bool
    
    // UI interaction states
    IsThinking       bool                    // LLM processing indicator
    ThinkingDots     int                    // Animation frame
    IsConfigOpen     bool                   // Configuration modal
    ConfigCursor     int                    // Config selection
    
    // Command execution flow
    PendingCommand              *CommandSuggestionMsg
    PendingFunctionCall         *FunctionCallMsg  
    AwaitingCommandApproval     bool
}
```

#### Constructor

```go
func NewModel(agent *llmServer.Agent) Model
```

**Parameters**:
- `agent`: Initialized LLM agent instance

**Returns**:
- `Model`: UI model ready for Bubble Tea program

### Message System

#### Core Message Types

##### User Input Processing

```go
type tea.KeyMsg struct {
    Type  tea.KeyType
    Runes []rune
    Alt   bool
}
```

**Key Bindings**:
- `Enter`: Submit query or approve command
- `Ctrl+C/X/V`: Clipboard operations
- `Ctrl+A`: Select all text
- `Ctrl+L`: Clear screen
- `Escape`: Cancel operations
- `Arrow Keys`: Navigation
- `Page Up/Down`: Scrolling

##### LLM Communication

```go
type LLMResponseMsg struct {
    Content string
    Error   error
}
```

**Usage**: Carries LLM responses back to UI for display and processing.

##### Command Execution

```go
type CommandSuggestionMsg struct {
    Suggestion *llmServer.CommandSuggestion
    Error      error
}
```

**Flow**: 
1. LLM suggests command
2. UI prompts for approval
3. User accepts/rejects
4. Command executes if approved

##### Function Calls

```go
type FunctionCallMsg struct {
    FunctionName string
    Arguments    map[string]interface{}
    Error        error
}
```

**Supported Functions**:
- `read_file(path)`: File reading
- `write_file(path, content)`: File writing  
- `create_pull_request(title, description)`: GitHub PR creation

#### Custom Commands

**Async Operations**: UI uses Bubble Tea commands for non-blocking operations:

```go
func QueryAgentCmd(agent *llmServer.Agent, query string) tea.Cmd {
    return func() tea.Msg {
        response, err := agent.ProcessQuery(query)
        return LLMResponseMsg{Content: response, Error: err}
    }
}
```

### View System

#### Layout Components

##### Chat Interface

**Structure**:
```
┌─ Configuration ─┐  ┌─────── Chat History ───────┐
│ Model: GPT-4    │  │ User: How do I install X?  │
│ Provider: OpenAI│  │ AI: You can use npm...     │  
│ Status: Ready   │  │ [Suggested Command]        │
└─────────────────┘  │ $ npm install X            │
                     │ [Execute] [Cancel]         │
                     └────────────────────────────┘
                     ┌─────── Input Area ─────────┐
                     │ > Your message here...     │
                     └────────────────────────────┘
```

##### Text Rendering

```go
func (m Model) renderMessage(msg Message) string
```

**Features**:
- Syntax highlighting for different message types
- Word wrapping for long content
- Unicode support
- Color coding for users vs. system vs. AI

##### Responsive Design

**Viewport Adaptation**:
- Automatic layout adjustment for terminal size changes
- Efficient scrolling for long conversations  
- Content-aware line breaking
- Minimum size constraints

### Text Editing Engine

#### Advanced Input Handling

**Capabilities**:
- Multi-line text input
- Text selection with mouse/keyboard
- Copy/paste with system clipboard
- Unicode-aware cursor positioning
- Horizontal scrolling for long lines

#### Selection System

```go
func (m *Model) updateSelection(x, y int) {
    if m.IsHighlighting {
        m.SelectionEndX, m.SelectionEndY = x, y
    }
}
```

**Features**:
- Click and drag selection
- Keyboard selection with Shift+arrows
- Select all with Ctrl+A
- Visual selection highlighting

## Function Calling System

### Function Definition

**Schema Format**: Functions are defined using structured schemas for LLM understanding:

```go
type FunctionDefinition struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`  
    Parameters  map[string]interface{} `json:"parameters"`
}
```

### Built-in Functions

#### File Operations

##### read_file

```go
Parameters: {
    "type": "object",
    "properties": {
        "path": {
            "type": "string",
            "description": "File path to read"
        }
    },
    "required": ["path"]
}
```

**Implementation**:
```go
func executeReadFile(path string) (string, error) {
    content, err := ioutil.ReadFile(path)
    return string(content), err
}
```

##### write_file

```go
Parameters: {
    "type": "object", 
    "properties": {
        "path": {"type": "string", "description": "Target file path"},
        "content": {"type": "string", "description": "File content"},
        "line_number": {"type": "integer", "description": "Insert at line (optional)"}
    },
    "required": ["path", "content"]
}
```

**Features**:
- Full file replacement or line insertion
- Backup creation for safety
- Permission validation
- Atomic write operations

#### Git Operations

##### create_pull_request

```go
Parameters: {
    "type": "object",
    "properties": {
        "title": {"type": "string", "description": "PR title"},
        "description": {"type": "string", "description": "PR body"},
        "branch": {"type": "string", "description": "Source branch (optional)"}
    },
    "required": ["title", "description"]  
}
```

**Process**:
1. Validate git repository
2. Check for uncommitted changes
3. Create/switch to branch
4. Commit changes
5. Push to GitHub
6. Create PR via API

### Extension API

#### Adding New Functions

**Steps**:

1. **Define Function Schema**:
```go
newFunction := FunctionDefinition{
    Name: "my_function",
    Description: "Description of what it does",
    Parameters: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "param1": map[string]interface{}{
                "type": "string",
                "description": "Parameter description",
            },
        },
        "required": []string{"param1"},
    },
}
```

2. **Implement Handler**:
```go
func executeMyFunction(args map[string]interface{}) (interface{}, error) {
    param1 := args["param1"].(string)
    // Implementation logic
    return result, nil
}
```

3. **Register Function**:
```go
func init() {
    RegisterFunction("my_function", executeMyFunction)
}
```

## Integration APIs

### LLM Provider Integration

#### Interface Requirements

**LangChain Go Compatibility**: New providers must implement the `llms.Model` interface:

```go
type Model interface {
    GenerateContent(ctx context.Context, messages []MessageContent, options ...CallOption) (*ContentResponse, error)
    GenerateContentStream(ctx context.Context, messages []MessageContent, options ...CallOption) (*ContentResponseStream, error)
}
```

#### Adding New Providers

**Example: Custom Provider**:

```go
type CustomProvider struct {
    apiKey   string
    baseURL  string
    client   *http.Client
}

func (p *CustomProvider) GenerateContent(ctx context.Context, messages []MessageContent, options ...CallOption) (*ContentResponse, error) {
    // Implementation for your provider's API
    response, err := p.client.Post(p.baseURL, buildRequest(messages))
    return parseResponse(response), err
}
```

**Integration Steps**:
1. Implement `llms.Model` interface
2. Add provider detection in `llmFactory.go`
3. Update model selection UI
4. Add configuration options

### External Tool Integration

#### Shell Command Integration

**Safe Execution Framework**:

```go
func ExecuteCommand(cmd string, requireApproval bool) error {
    if requireApproval {
        approved := requestUserApproval(cmd)
        if !approved {
            return errors.New("command rejected by user")
        }
    }
    
    return exec.Command("sh", "-c", cmd).Run()
}
```

**Safety Levels**:
- **Safe**: File operations, git status, ls, etc.
- **Moderate**: Package installs, git commits
- **Dangerous**: rm, sudo commands, system modifications

#### API Integration

**HTTP Client Framework**:

```go
func CallExternalAPI(endpoint string, payload interface{}) (*http.Response, error) {
    client := &http.Client{Timeout: 30 * time.Second}
    data, _ := json.Marshal(payload)
    return client.Post(endpoint, "application/json", bytes.NewBuffer(data))
}
```

## Extension Points

### Plugin Architecture

#### Plugin Interface

```go
type Plugin interface {
    Name() string
    Description() string  
    Initialize(config map[string]interface{}) error
    Execute(command string, args []string) (interface{}, error)
    Cleanup() error
}
```

#### Registration System

```go
var pluginRegistry = make(map[string]Plugin)

func RegisterPlugin(name string, plugin Plugin) {
    pluginRegistry[name] = plugin
}

func GetPlugin(name string) (Plugin, bool) {
    plugin, exists := pluginRegistry[name]
    return plugin, exists
}
```

### Configuration System

#### Settings Management

```go
type Configuration struct {
    DefaultProvider string            `json:"default_provider"`
    DefaultModel    string            `json:"default_model"`  
    APIKeys         map[string]string `json:"api_keys"`
    Preferences     UserPreferences   `json:"preferences"`
}
```

**Persistence**: Configuration stored in user's home directory as JSON.

### Event System

#### Event Types

```go
type EventType int

const (
    EventUserInput EventType = iota
    EventLLMResponse
    EventCommandExecution
    EventError
    EventConfigChange
)
```

#### Event Handlers

```go
type EventHandler func(event Event) error

func RegisterEventHandler(eventType EventType, handler EventHandler) {
    eventHandlers[eventType] = append(eventHandlers[eventType], handler)
}
```

**Use Cases**:
- Logging user interactions
- Analytics collection
- Custom integrations
- Debugging and monitoring

This API documentation provides comprehensive coverage of Menace CLI's internal APIs and extension points, enabling developers to understand, modify, and extend the system effectively.