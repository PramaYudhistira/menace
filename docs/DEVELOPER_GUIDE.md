# Menace CLI Developer Guide

## Project Overview

Menace CLI is a sophisticated terminal-based AI assistant that combines the power of large language models with a rich text user interface. The project demonstrates modern Go development practices, cross-platform deployment, and effective AI integration patterns.

## Architecture Deep Dive

### Design Philosophy

The application follows a clean separation of concerns with three main layers:

1. **Presentation Layer** (`ui/`): Handles user interface and interaction
2. **Business Logic Layer** (`llmServer/`): Manages AI conversations and command processing  
3. **System Integration Layer**: Handles file operations, shell commands, and GitHub integration

### Component Relationships

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   UI/TUI        │◄──►│   Agent         │◄──►│   LLM APIs      │
│   (Bubble Tea)  │    │   (Core Logic)  │    │   (Providers)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Messages      │    │   Commands      │    │   File System   │
│   (Chat State)  │    │   (Execution)   │    │   (Operations)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Code Organization

### Directory Structure Explained

```
menace-cli/
├── bin/                    # Compiled binaries and Node.js wrapper
│   └── menace.js          # Platform-specific binary launcher
├── scripts/               # Build automation
│   ├── build-all.js      # Cross-platform compilation
│   ├── build-mac.js      # macOS-specific builds
│   └── build-windows.js  # Windows-specific builds
└── src/                   # Go source code
    ├── main.go           # Application entry point
    ├── go.mod            # Go module definition
    ├── llmServer/        # AI and command processing
    │   ├── agent.go      # Core agent logic
    │   ├── command_suggestion.go  # Command parsing
    │   ├── github_interface.go    # GitHub integration
    │   ├── llmFactory.go          # LLM provider management
    │   ├── llmService.go          # Service abstractions
    │   └── system_prompt.go       # Dynamic prompt generation
    └── ui/               # User interface components
        ├── executor.go   # Command execution
        ├── message.go    # Message types and handling
        ├── model.go      # Main UI state management
        ├── model_selection.go  # Model switching UI
        ├── thinking.go   # Loading animations
        ├── ui_utils.go   # Utility functions
        ├── update.go     # Bubble Tea update loop
        └── view.go       # Rendering and layout
```

## Key Design Patterns

### 1. Model-View-Controller (MVC)

- **Model**: `ui/model.go` - Contains all application state
- **View**: `ui/view.go` - Handles rendering and layout
- **Controller**: `ui/update.go` - Processes user input and state transitions

### 2. Command Pattern

Commands are encapsulated as structured data:

```go
type CommandSuggestionMsg struct {
    Command string
    Reason  string
}

type FunctionCallMsg struct {
    Name string
    Args map[string]interface{}
}
```

### 3. Strategy Pattern

LLM providers are interchangeable through a common interface:

```go
// Different providers implement the same interface
switch provider {
case "anthropic":
    llm, err := anthropic.New(...)
case "openai":  
    llm, err := openai.New(...)
case "ollama":
    llm, err := ollama.New(...)
}
```

### 4. Observer Pattern

Bubble Tea's message-passing system acts as an observer pattern implementation, where UI components respond to events.

## State Management

### UI State Flow

```
User Input → Update Function → Model State → View Render → Terminal Display
     ▲                                                         │
     └─────────────── User Interaction ◄──────────────────────┘
```

### Agent State Flow

```
User Message → Agent.SendMessage → LLM API → Response Processing → Command/Function Extraction
      ▲                                                                        │
      └──────────────── Conversation History Update ◄────────────────────────┘
```

## Concurrency Handling

### Thread Safety

The `Agent` struct uses a mutex to ensure thread-safe operations:

```go
func (a *Agent) SendMessage(ctx context.Context, input string) (string, *CommandSuggestion, error) {
    a.mu.Lock()
    defer a.mu.Unlock()
    // ... thread-safe operations
}
```

### Async Operations

- LLM API calls are non-blocking
- File operations are synchronous for safety
- UI updates are handled through Bubble Tea's event system

## Error Handling Strategy

### Layered Error Handling

1. **System Level**: OS-specific errors (file permissions, network)
2. **Application Level**: Business logic errors (invalid commands)
3. **User Level**: User-friendly error messages in the UI

### Error Types

```go
// System errors - logged and handled gracefully
if err != nil {
    return fmt.Errorf("failed to create OpenAI client: %v", err)
}

// User errors - displayed in chat interface
m.AddSystemMessage("Error: " + err.Error())

// Recoverable errors - retry mechanisms
if cmdSuggestion == nil {
    // Continue with text-only response
}
```

## Testing Strategy

### Unit Testing

Focus areas for unit tests:

- `command_suggestion.go`: Command parsing logic
- `ui_utils.go`: Utility functions  
- Message handling and state transitions

### Integration Testing

- LLM provider switching
- File operation workflows
- GitHub integration flows

### Manual Testing

- Cross-platform binary execution
- UI responsiveness and layout
- Command execution safety

## Performance Considerations

### Memory Management

- **Conversation History**: Implement sliding window to prevent unbounded growth
- **Message Buffers**: Efficient string handling for large responses
- **UI Rendering**: Minimize redraws and optimize scroll calculations

### Response Time Optimization

- **Streaming Responses**: Real-time text display as tokens arrive
- **Caching**: Cache frequently used file contents
- **Lazy Loading**: Load model configurations on demand

## Security Best Practices

### Input Validation

```go
// Sanitize user input before processing
func validateCommand(cmd string) error {
    if strings.Contains(cmd, "rm -rf /") {
        return errors.New("dangerous command detected")
    }
    return nil
}
```

### Command Execution Safety

- Always require user approval for destructive commands
- Sandbox command execution to user's current permissions
- Log all executed commands for audit trails

### API Key Management

- Environment variable storage only
- No hardcoded keys in source
- Secure transmission to LLM providers

## Build System Deep Dive

### Multi-Platform Compilation

The build system handles Go's cross-compilation features:

```javascript
// Build for all platforms simultaneously
const buildTargets = [
    { os: "linux", arch: "amd64", output: "menace-go-linux" },
    { os: "darwin", arch: "amd64", output: "menace-go-darwin" },
    { os: "darwin", arch: "arm64", output: "menace-go-darwin-arm64" },
    { os: "windows", arch: "amd64", output: "menace-go-win.exe" }
];

// Cross-compilation environment variables
env: { ...process.env, GOOS, GOARCH }
```

### Node.js Wrapper Benefits

1. **Simplified Distribution**: Single npm package for all platforms
2. **Automatic Binary Selection**: Runtime platform detection
3. **Process Management**: Proper process spawning and cleanup
4. **Error Handling**: Graceful failure modes

## Extension Points

### Adding New LLM Providers

1. Implement the `llms.Model` interface
2. Add provider case in `SetModel()` function
3. Update model selection UI
4. Add configuration documentation

### Custom Function Calls

1. Define function signature in system prompt
2. Add parsing logic in command processing
3. Implement execution logic in `ExecuteFunctionCall()`
4. Add error handling and validation

### UI Components

1. Create new view functions following existing patterns
2. Add state management in `Model` struct
3. Handle keyboard shortcuts in update loop
4. Test across different terminal sizes

## Debugging Techniques

### Development Workflow

1. **Local Development**: Use `go run` for rapid iteration
2. **Debug Logging**: Add debug prints to trace execution flow
3. **UI Debugging**: Test different terminal sizes and color schemes
4. **LLM Debugging**: Log request/response pairs for analysis

### Common Issues

1. **Terminal Compatibility**: Test across different terminal emulators
2. **Unicode Handling**: Ensure proper rune vs byte handling
3. **File Path Issues**: Handle Windows vs Unix path differences
4. **Environment Variables**: Test with and without required variables

## Code Style Guidelines

### Go Best Practices

- Use `gofmt` for consistent formatting
- Follow effective Go idioms (error handling, interfaces)
- Write descriptive variable names
- Keep functions focused and single-purpose

### Documentation Standards

- Document all public functions and types
- Provide usage examples for complex functions
- Explain design decisions in comments
- Keep documentation up-to-date with code changes

## Deployment Considerations

### Distribution Strategy

1. **npm Package**: Primary distribution method
2. **GitHub Releases**: Binary attachments for direct download
3. **Version Management**: Semantic versioning with changelogs
4. **CI/CD Pipeline**: Automated testing and release process

### Runtime Dependencies

- **Go Runtime**: Not required (statically linked binaries)
- **Node.js**: Required for wrapper and npm installation
- **System Tools**: Platform-specific clipboard and shell access

## Future Architecture Considerations

### Scalability

- Plugin system for extensibility
- Configuration file support
- Multiple conversation contexts
- Cloud synchronization options

### Performance

- Response streaming implementation
- Local model support optimization
- Cache layer for frequent operations
- Reduced binary size optimization

### User Experience

- Keyboard shortcut customization
- Theme and color customization
- Enhanced error recovery
- Offline mode capabilities