# Architecture Documentation

This document provides a comprehensive overview of the Menace CLI architecture, design patterns, and system interactions.

## High-Level Architecture

Menace CLI follows a layered architecture with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                         User Interface                      │
│                      (Bubble Tea TUI)                      │
├─────────────────────────────────────────────────────────────┤
│                      Application Layer                      │
│                    (Event Handling & State)                │
├─────────────────────────────────────────────────────────────┤
│                        LLM Server                          │
│                   (Agent & Conversation)                   │
├─────────────────────────────────────────────────────────────┤
│                      Platform Layer                        │
│                 (Shell, File System, Git)                  │
└─────────────────────────────────────────────────────────────┘
```

## Component Architecture

### 1. Entry Point (`main.go`)

The application bootstrap sequence:

```
main()
├── Validate Environment (API Keys)
├── Initialize LLM Agent
├── Setup Bubble Tea Program
├── Configure UI Model
└── Start Event Loop
```

**Key Responsibilities:**
- Environment validation
- Dependency injection
- Error handling for initialization failures
- Program lifecycle management

### 2. LLM Server Layer (`llmServer/`)

#### Agent Architecture

```
Agent
├── llms.Model (LangChainGo Interface)
│   ├── OpenAI Client
│   ├── Anthropic Client
│   └── Ollama Client
├── Conversation State
│   ├── Message History
│   ├── System Prompt
│   └── Context Management
└── Command Processing
    ├── Response Parsing
    ├── Function Call Extraction
    └── Command Suggestion Generation
```

#### Key Components

**`agent.go`** - Core Intelligence
- Thread-safe conversation management
- Multi-provider LLM abstraction
- Dynamic model switching
- Message chain persistence

**`system_prompt.go`** - Behavior Configuration  
- Context-aware prompt generation
- Platform-specific instructions
- Function call documentation
- Safety constraints

**`command_suggestion.go`** - Response Processing
- Structured response parsing
- Command extraction and validation
- Approval workflow determination

**`llmFactory.go`** - Platform Detection
- Shell environment detection
- OS-specific configurations
- Cross-platform compatibility

### 3. User Interface Layer (`ui/`)

#### Model-View-Update Pattern

Follows Elm/Bubble Tea architecture:

```
Model (State)
├── Agent Reference
├── Message History
├── Input State
│   ├── Text Buffer
│   ├── Cursor Position
│   └── Selection Range
├── UI State
│   ├── Window Dimensions  
│   ├── Viewport Position
│   └── Component Focus
└── Execution State
    ├── Pending Commands
    ├── Approval Queue
    └── Async Operations

Update (Event Handler)
├── Keyboard Events
├── Mouse Events
├── Window Resize
├── Async Command Results
└── Model State Transitions

View (Renderer)
├── Layout Management
├── Component Rendering
├── Styling Application
└── Text Formatting
```

#### Key Components

**`model.go`** - State Management
- Complex UI state coordination
- Text editing capabilities (selection, clipboard)
- Async operation handling
- Window management

**`view.go`** - Rendering Engine
- Three-panel layout system
- Syntax highlighting
- Interactive elements
- Responsive design

**`update.go`** - Event Processing
- Comprehensive input handling
- Command approval workflow
- Async communication with LLM
- State synchronization

**`executor.go`** - Command Execution
- Cross-platform shell execution
- File system operations
- Git integration
- Function call dispatch

## Data Flow Architecture

### Message Processing Pipeline

```
User Input
    ↓
Input Validation
    ↓
UI State Update
    ↓
LLM Agent.SendMessage()
    ↓
LangChainGo Provider
    ↓
LLM Response Processing
    ↓
Command/Function Extraction
    ↓
Approval Check
    ↓
Execution (if approved)
    ↓
Result Display
    ↓
State Persistence
```

### Command Execution Flow

```
Parsed Command/Function
    ↓
Approval Required? ──No──→ Direct Execution
    ↓ Yes
User Approval Prompt
    ↓
User Decision (y/n/e)
    ↓
Execute ──→ Capture Output ──→ Display Result
    ↓           ↓                   ↓
Cancel      Error Handling      Success State
    ↓           ↓                   ↓
No-op       Error Display       History Update
```

### Async Command Pattern

```go
type AsyncCommand struct {
    ID       string
    Command  string
    Callback func(string, error)
}

// In update.go
func (m Model) executeAsync(cmd string) tea.Cmd {
    return func() tea.Msg {
        output, err := executeCommand(cmd)
        return AsyncResult{
            Output: output,
            Error:  err,
        }
    }
}
```

## Cross-Platform Architecture

### Platform Abstraction

```
Platform Interface
├── Shell Detection
│   ├── Windows (CMD/PowerShell)
│   ├── macOS (bash/zsh)  
│   └── Linux (bash)
├── File Operations
│   ├── Path Handling
│   ├── Permission Management
│   └── Cross-platform APIs
└── Clipboard Integration
    ├── OS-specific implementations
    └── Fallback mechanisms
```

### Build System Architecture

```
Build Orchestration (Node.js)
├── Target Definition
│   ├── GOOS/GOARCH Matrix
│   ├── Binary Naming
│   └── Output Organization
├── Go Cross-compilation
│   ├── Environment Setup
│   ├── Dependency Resolution
│   └── Binary Generation
└── Distribution Wrapper
    ├── Platform Detection
    ├── Binary Selection
    └── Process Management
```

## Security Architecture

### Input Validation Pipeline

```
User Input
    ↓
Length Validation
    ↓
Pattern Filtering (dangerous commands)
    ↓
Shell Escape Validation
    ↓
Context-aware Sanitization
    ↓
Approved for Processing
```

### Approval Workflow

```go
type ApprovalLevel int

const (
    AutoApprove ApprovalLevel = iota  // Safe operations
    UserApproval                      // Potentially destructive
    ExplicitApproval                  // High-risk operations
)

func determineApprovalLevel(operation Operation) ApprovalLevel {
    switch operation.Type {
    case "ReadFile":
        return AutoApprove
    case "WriteFile", "DeleteFile":
        return UserApproval
    case "SystemCommand":
        return ExplicitApproval
    default:
        return UserApproval
    }
}
```

## Error Handling Architecture

### Error Propagation Pattern

```
Error Source
    ↓
Component Error Handler
    ↓
Bubble Tea Command System
    ↓
UI Error Display
    ↓
User Notification
    ↓
Graceful Recovery
```

### Error Types

```go
type ErrorType int

const (
    NetworkError ErrorType = iota
    FileSystemError
    LLMError
    ValidationError
    SystemError
)

type MenaceError struct {
    Type    ErrorType
    Message string
    Context map[string]interface{}
    Wrapped error
}
```

## Performance Architecture

### Async Operations

```
Main UI Thread
├── Event Processing (non-blocking)
├── State Updates (synchronous)
└── View Rendering (fast)

Background Threads  
├── LLM Communication
├── File Operations
├── Shell Command Execution
└── Network Requests
```

### Memory Management

```
Message History
├── Sliding Window (configurable size)
├── Compression for old messages
└── Garbage Collection triggers

UI State
├── Viewport Virtualization
├── Lazy Rendering
└── Component Pooling
```

### Caching Strategy

```
LLM Responses
├── Local conversation cache
├── Command suggestion cache
└── Function call results cache

File System
├── File content cache
├── Directory listing cache
└── Git status cache
```

## Extensibility Architecture

### Plugin System Design

```
Extension Interface
├── Provider Registration
├── Function Registration
├── UI Component Registration
└── Configuration Integration

Extension Discovery
├── Built-in Extensions
├── External Plugin Loading
└── Configuration-based Enablement
```

### Provider Architecture

```go
type LLMProvider interface {
    Initialize(config Config) error
    SendMessage(context.Context, []Message) (Response, error)
    GetModels() []Model
    ValidateConfig(Config) error
}

type ProviderRegistry struct {
    providers map[string]LLMProvider
    mutex     sync.RWMutex
}
```

## Configuration Architecture

### Configuration Sources

```
Configuration Hierarchy
├── Environment Variables (highest priority)
├── Configuration File (~/.menace/config.yaml)
├── Runtime Settings (in-memory)
└── Default Values (lowest priority)
```

### Configuration Schema

```yaml
llm:
  default_provider: "openai"
  providers:
    openai:
      api_key: "${OPENAI_API_KEY}"
      model: "gpt-4"
    anthropic:
      api_key: "${ANTHROPIC_API_KEY}"
      model: "claude-3-sonnet"

ui:
  theme: "dark"
  mouse_support: true
  word_wrap: true

execution:
  auto_approve_safe: true
  command_timeout: 30s
  max_history: 1000
```

## Monitoring and Observability

### Logging Architecture

```
Log Levels
├── DEBUG: Internal state changes
├── INFO: User actions and responses
├── WARN: Recoverable errors
├── ERROR: System failures
└── FATAL: Unrecoverable errors

Log Destinations
├── Console (development)
├── File (~/.menace/logs/)
└── External systems (configurable)
```

### Metrics Collection

```
Performance Metrics
├── Response Times (LLM, File ops)
├── Memory Usage
├── Command Success Rates
└── User Interaction Patterns

Usage Metrics
├── Feature Usage Frequency
├── Provider Performance
├── Error Rates
└── Session Duration
```

## Deployment Architecture

### Binary Distribution

```
Distribution Strategy
├── Platform-specific Binaries
│   ├── Windows (x64)
│   ├── macOS (Intel/ARM)
│   └── Linux (x64/ARM)
├── Node.js Wrapper
│   ├── Platform Detection
│   ├── Binary Selection
│   └── Process Management
└── Package Management
    ├── npm Registry
    ├── GitHub Releases
    └── Package Managers
```

### Update Mechanism

```
Update Flow
├── Version Check (GitHub API)
├── Download New Binary
├── Verification (checksums)
├── Atomic Replacement
└── Cleanup Old Version
```

## Future Architecture Considerations

### Scalability Improvements

1. **Streaming Responses**: Real-time LLM output rendering
2. **Parallel Execution**: Multiple command execution
3. **Distributed Mode**: Remote LLM providers
4. **Plugin Ecosystem**: Third-party extensions

### Performance Enhancements

1. **Response Caching**: Intelligent cache invalidation
2. **Predictive Loading**: Anticipate user needs
3. **Background Processing**: Non-blocking operations
4. **Resource Optimization**: Memory and CPU efficiency

### Security Enhancements

1. **Sandboxed Execution**: Isolated command execution
2. **Permission System**: Granular access control
3. **Audit Logging**: Complete operation tracking
4. **Secure Storage**: Encrypted configuration/cache

This architecture enables Menace CLI to provide a responsive, secure, and extensible terminal assistant while maintaining clear separation of concerns and robust error handling.