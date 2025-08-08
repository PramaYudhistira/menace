# Menace CLI Architecture

This document provides a comprehensive overview of the Menace CLI architecture, detailing how the Node.js wrapper and Go backend work together to create an intelligent terminal assistant.

## Overview

Menace CLI follows a **hybrid architecture** combining Node.js for cross-platform distribution and Go for high-performance terminal UI and LLM integration. The system provides a sophisticated TUI (Terminal User Interface) that interfaces with multiple Large Language Model providers to offer intelligent command suggestions, code assistance, and terminal automation.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    User Terminal                            │
├─────────────────────────────────────────────────────────────┤
│  Node.js Wrapper (menace.js)                              │
│  ├─ Cross-platform binary detection                        │
│  ├─ Process spawning and management                        │
│  └─ stdio inheritance                                      │
├─────────────────────────────────────────────────────────────┤
│  Go Backend (menace-go)                                    │
│  ├─ Bubble Tea TUI Framework                              │
│  ├─ LLM Agent System                                       │
│  ├─ Command Processing Engine                              │
│  └─ File System Integration                                │
├─────────────────────────────────────────────────────────────┤
│  External Services                                          │
│  ├─ OpenAI API                                            │
│  ├─ Anthropic API                                          │
│  ├─ Ollama (Local LLMs)                                    │
│  └─ GitHub API                                             │
└─────────────────────────────────────────────────────────────┘
```

## Component Architecture

### 1. Node.js Distribution Layer

**Purpose**: Provides cross-platform compatibility and npm distribution.

**Components**:
- **Binary Dispatcher** (`bin/menace.js`): Platform detection and executable selection
- **Build System** (`scripts/`): Cross-compilation management
- **Package Management** (`package.json`): npm integration

**Key Responsibilities**:
- Detect user's platform (Windows, macOS Intel/ARM, Linux)
- Select appropriate Go binary
- Spawn child process with inherited stdio
- Handle process lifecycle and error propagation

```javascript
// Platform detection and binary selection
switch (process.platform) {
    case "win32": binName = "menace-go-win.exe"; break;
    case "darwin": binName = process.arch === "arm64" 
        ? "menace-go-darwin-arm64" 
        : "menace-go-darwin"; break;
    case "linux": binName = "menace-go-linux"; break;
}
```

### 2. Go Backend Core

**Purpose**: High-performance TUI and LLM integration engine.

#### 2.1 Application Entry Point (`main.go`)

**Initialization Sequence**:
1. Environment validation (API keys)
2. LLM agent initialization
3. Bubble Tea UI program creation
4. Event loop startup

```go
func main() {
    apiKey := os.Getenv("OPENAI_API_KEY")
    agent, err := llmServer.NewAgent(apiKey)
    p := tea.NewProgram(
        ui.NewModel(agent),
        tea.WithAltScreen(),
        tea.WithMouseAllMotion(),
    )
    p.Run()
}
```

#### 2.2 LLM Server Package (`llmServer/`)

**Core Components**:

##### Agent System (`agent.go`)
- **Multi-Provider Support**: OpenAI, Anthropic, Ollama
- **Context Management**: Conversation history and state
- **Thread Safety**: Concurrent request handling
- **Dynamic Model Switching**: Runtime provider changes

```go
type Agent struct {
    llm      llms.Model              // LangChain Go interface
    mu       sync.Mutex              // Thread safety
    shell    string                  // Environment detection
    messages []llms.MessageContent   // Conversation context
    ctx      context.Context         // Request context
    provider string                  // Current provider
    Model    string                  // Active model
}
```

##### System Prompt Engineering (`system_prompt.go`)
- **Adaptive Instructions**: Shell-specific prompting
- **Function Calling**: Structured tool usage
- **Safety Measures**: Human approval gates
- **Git Integration**: Version control awareness

##### Command Processing (`command_suggestion.go`)
- **Response Parsing**: Structured command extraction
- **Safety Validation**: Risk assessment
- **Execution Pipeline**: Approval → Execute → Feedback

```go
type CommandSuggestion struct {
    Reason  string
    Command string
    AwaitingCommandApproval bool
}
```

##### GitHub Integration (`github_interface.go`)
- **Pull Request Automation**: Automated PR creation
- **Repository Operations**: Branch management
- **API Integration**: GitHub REST API client

#### 2.3 UI Package (`ui/`)

**Architecture**: Model-View-Update pattern (Elm Architecture)

##### Core Model (`model.go`)
- **State Management**: Centralized application state
- **Text Editing Engine**: Multi-line input with selection
- **Display Properties**: Viewport and layout management
- **Interaction States**: Modal dialogs and user flows

```go
type Model struct {
    // Core state
    Input    string
    Messages []Message
    agent    *llmServer.Agent
    
    // Advanced text editing
    CursorX, CursorY    int
    SelectionStartX, SelectionStartY int
    SelectionEndX, SelectionEndY     int
    
    // UI state management
    IsThinking       bool
    IsConfigOpen     bool
    PendingCommand   *CommandSuggestionMsg
}
```

##### Event Processing (`update.go`)
- **Message Routing**: Event → Handler mapping
- **Asynchronous Operations**: Non-blocking LLM calls
- **State Transitions**: UI flow management
- **Keyboard/Mouse Events**: Comprehensive input handling

**Message Flow**:
```
KeyMsg → ProcessInput → AgentQuery → LLMResponseMsg → CommandParsing → ExecutionApproval
```

##### View Rendering (`view.go`)
- **Responsive Layout**: Terminal size adaptation
- **Text Formatting**: Syntax highlighting and styling
- **Scrolling System**: Efficient content navigation
- **Component Composition**: Modular UI elements

##### Advanced Features
- **Model Selection** (`model_selection.go`): Runtime provider switching
- **Function Execution** (`executor.go`): File operations and system integration
- **Visual Feedback** (`thinking.go`): Loading animations

## Data Flow Architecture

### 1. Request Processing Pipeline

```
User Input → Text Processing → Agent Querying → LLM Response → 
Command Parsing → Safety Check → User Approval → Execution → Feedback Display
```

### 2. Message-Driven Architecture

**Bubble Tea Commands**: Asynchronous operations

```go
// Example: LLM query as async command
func QueryAgent(agent *llmServer.Agent, query string) tea.Cmd {
    return tea.Tick(time.Millisecond, func(t time.Time) tea.Msg {
        response := agent.ProcessQuery(query)
        return LLMResponseMsg{Content: response}
    })
}
```

### 3. State Synchronization

**Thread-Safe Operations**:
- Mutex-protected agent state
- Atomic UI updates
- Consistent message ordering

## Security Architecture

### 1. Command Approval System

**Human-in-the-Loop**: Critical operations require explicit user consent:
- File modifications
- Git operations
- System administration commands
- API calls

### 2. Environment Isolation

**Controlled Execution**:
- Subprocess sandboxing
- Error containment
- Output sanitization
- Resource limits

### 3. API Key Management

**Secure Credential Handling**:
- Environment variable requirements
- No credential persistence
- Runtime validation
- Provider isolation

## Performance Characteristics

### 1. Memory Management

**Efficient Resource Usage**:
- Streaming LLM responses
- Lazy UI rendering
- Bounded conversation history
- Garbage collection friendly

### 2. Concurrent Processing

**Non-Blocking Operations**:
- Asynchronous LLM calls
- Parallel UI updates
- Background file operations
- Responsive user interaction

### 3. Cross-Platform Optimization

**Platform-Specific Builds**:
- Native binary compilation
- OS-specific optimizations
- Minimal runtime dependencies
- Fast startup times

## Extensibility Points

### 1. LLM Provider Integration

**New Provider Addition**:
1. Implement `llms.Model` interface
2. Add provider detection logic
3. Update factory pattern
4. Add configuration options

### 2. Function Calling System

**New Tool Integration**:
1. Define function schema
2. Implement handler logic
3. Add safety checks
4. Update system prompts

### 3. UI Component Extension

**New Interface Elements**:
1. Create Bubble Tea component
2. Add message types
3. Implement update logic
4. Integrate view rendering

## Build and Deployment

### 1. Cross-Platform Compilation

**Multi-Target Build System**:
- Linux (AMD64)
- macOS (Intel/ARM)
- Windows (AMD64)

### 2. Distribution Strategy

**NPM Package Distribution**:
- Pre-compiled binaries
- Platform auto-detection
- Single `npm install` command
- Version management

### 3. Development Workflow

**Local Development**:
1. Go module development
2. Cross-platform testing
3. NPM packaging
4. Integration testing

## Configuration Management

### 1. Environment Variables

**Required Configuration**:
- `OPENAI_API_KEY`: OpenAI models
- `ANTHROPIC_API_KEY`: Claude models (optional)
- `GITHUB_TOKEN`: Git operations (optional)

### 2. Runtime Configuration

**Dynamic Settings**:
- Model selection
- Provider preferences
- UI customization
- Feature toggles

## Error Handling Strategy

### 1. Graceful Degradation

**Failure Modes**:
- API unavailability → Offline mode
- Network issues → Local fallback
- Auth failures → Clear error messages
- Resource limits → Resource cleanup

### 2. User Experience

**Error Communication**:
- Clear error messages
- Recovery suggestions
- Non-fatal error handling
- State preservation

## Future Architecture Considerations

### 1. Plugin System

**Extensibility Framework**:
- WebAssembly plugins
- Language bindings
- Tool integrations
- Custom providers

### 2. Distributed Architecture

**Scaling Considerations**:
- Client-server separation
- Shared context management
- Multi-user support
- Cloud integration

### 3. Advanced AI Features

**Enhanced Capabilities**:
- Multi-modal input
- Context-aware automation
- Learning from usage
- Predictive assistance

This architecture provides a robust foundation for an intelligent terminal assistant while maintaining extensibility, security, and performance across multiple platforms and use cases.