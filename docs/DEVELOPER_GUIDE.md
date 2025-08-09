# Developer Guide

This guide provides comprehensive information for developers who want to contribute to or understand the Menace CLI codebase.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Development Setup](#development-setup)
- [Code Organization](#code-organization)
- [Key Components](#key-components)
- [Build System](#build-system)
- [Testing](#testing)
- [Debugging](#debugging)
- [Contributing](#contributing)

## Architecture Overview

Menace CLI is built with a modular architecture that separates concerns between AI interaction, user interface, and system operations.

```
┌─────────────────────────────────────────┐
│                Main Entry               │
│              (main.go)                  │
├─────────────────┬───────────────────────┤
│   UI Package    │    LLMServer Package  │
│                 │                       │
│ ┌─────────────┐ │ ┌─────────────────────┐│
│ │   Model     │ │ │      Agent          ││
│ │   View      │◄┼─┤   Command Parser    ││
│ │  Executor   │ │ │   System Prompts    ││
│ │  Messages   │ │ │   GitHub Interface  ││
│ └─────────────┘ │ └─────────────────────┘│
└─────────────────┴───────────────────────┘
```

### Core Principles

1. **Separation of Concerns**: UI logic is separate from AI logic
2. **Modularity**: Each package has a specific responsibility
3. **Extensibility**: Easy to add new LLM providers or UI components
4. **Cross-platform**: Works on Windows, macOS, and Linux

## Development Setup

### Prerequisites

- **Go 1.18+** (recommended 1.21+)
- **Node.js v14+** and npm
- **Git** for version control
- A supported terminal emulator

### Clone and Setup

```bash
# Clone the repository
git clone https://github.com/your-org/menace-cli.git
cd menace-cli/menace-cli

# Install Node.js dependencies (for build scripts)
npm install

# Set up environment variables
export OPENAI_API_KEY="your-openai-key-here"
# Optional: for Anthropic support
export ANTHROPIC_API_KEY="your-anthropic-key-here"
```

### Development Build

```bash
# Build for your current platform
cd src
go build -o ../bin/menace-dev .

# Run directly
go run .

# Or build all platforms
cd ..
npm run build
```

### IDE Setup

#### VS Code
Recommended extensions:
- Go extension by Google
- Error Lens
- GitLens
- Better Comments

#### GoLand
Import the project and ensure Go modules are enabled.

## Code Organization

```
menace-cli/
├── bin/                    # Compiled binaries
├── docs/                   # Documentation (generated)
│   ├── api/               # API documentation
│   └── guides/            # User and developer guides
├── package.json           # npm configuration and build scripts
├── scripts/               # Build scripts
│   ├── build-all.js       # Cross-platform build
│   ├── build-mac.js       # macOS specific build
│   └── build-windows.js   # Windows specific build
└── src/                   # Go source code
    ├── go.mod             # Go module definition
    ├── go.sum             # Go dependency checksums
    ├── main.go            # Application entry point
    ├── llmServer/         # AI and LLM functionality
    │   ├── agent.go       # Main AI agent
    │   ├── command_suggestion.go  # Command parsing
    │   ├── github_interface.go   # GitHub integration
    │   ├── llmFactory.go  # Environment detection
    │   ├── llmService.go  # Legacy LLM service
    │   └── system_prompt.go      # System prompts
    └── ui/                # Terminal user interface
        ├── executor.go    # File operations and command execution
        ├── message.go     # Message structures
        ├── model.go       # UI state management
        ├── model_selection.go  # Model selection UI
        ├── thinking.go    # Loading/thinking indicators
        ├── ui_utils.go    # UI utility functions
        ├── update.go      # Event handling
        └── view.go        # UI rendering
```

## Key Components

### 1. Main Entry Point (`main.go`)

The application bootstrap:
- Validates environment variables
- Initializes the AI agent
- Sets up the Bubble Tea program
- Handles graceful shutdown

### 2. LLMServer Package

#### Agent (`agent.go`)
The primary AI interface supporting multiple providers:
- OpenAI GPT models
- Anthropic Claude models
- Ollama local models

```go
type Agent struct {
    chain    *chains.ConversationChain
    memory   *memory.ConversationBufferMemory
    provider string
}
```

#### Command Suggestions (`command_suggestion.go`)
Parses structured commands from AI responses:
```go
type CommandSuggestion struct {
    Reason             string
    Command            string
    RequiresApproval   bool
    IsFileOperation    bool
}
```

#### System Prompts (`system_prompt.go`)
Defines the AI's behavior and capabilities. The prompt engineering is crucial for:
- Command generation accuracy
- Security and approval requirements
- Cross-platform compatibility
- File operation safety

### 3. UI Package

#### Model (`model.go`)
The core UI state machine handling:
- Text editing with selections
- Message history management
- AI interaction coordination
- Configuration management

#### View (`view.go`)
Renders the terminal interface using Lipgloss:
- Two-panel layout (sidebar + main area)
- Responsive design
- Syntax highlighting
- Mouse interaction support

#### Executor (`executor.go`)
Handles system operations:
- File reading with line numbers
- Diff-based file editing
- Shell command execution
- Safety validations

## Build System

### Build Scripts

The project uses Node.js scripts for cross-compilation:

#### `scripts/build-all.js`
Builds binaries for all supported platforms:
- Windows (amd64)
- macOS (amd64, arm64)
- Linux (amd64, arm64)

```javascript
const platforms = [
    { GOOS: 'windows', GOARCH: 'amd64', ext: '.exe' },
    { GOOS: 'darwin', GOARCH: 'amd64', ext: '' },
    { GOOS: 'darwin', GOARCH: 'arm64', ext: '' },
    { GOOS: 'linux', GOARCH: 'amd64', ext: '' },
    { GOOS: 'linux', GOARCH: 'arm64', ext: '' },
];
```

#### Build Commands
```bash
# Build all platforms
npm run build

# Build specific platforms
npm run build-windows
npm run build-mac

# Development build
cd src && go build .
```

### Go Build Configuration

The project uses standard Go modules with these key dependencies:
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/tmc/langchaingo` - LLM integration
- `github.com/lrstanley/bubblezone` - Mouse support

## Testing

### Unit Testing

Run tests with:
```bash
cd src
go test ./...
```

### Integration Testing

Test with different LLM providers:
```bash
# Test with OpenAI
OPENAI_API_KEY="sk-..." go run .

# Test with Anthropic
ANTHROPIC_API_KEY="sk-ant-..." go run .
```

### Manual Testing Checklist

- [ ] Text editing functionality (cursor, selection, copy/paste)
- [ ] AI message sending and receiving
- [ ] Command suggestion parsing and execution
- [ ] File operations (read, edit, create)
- [ ] Cross-platform compatibility
- [ ] Terminal resize handling
- [ ] Error handling and recovery

## Debugging

### Logging

Add debug logging:
```go
import "log"

// Add debugging information
log.Printf("Debug: %+v", debugData)
```

### Common Issues

#### 1. API Key Problems
```bash
# Verify environment variable
echo $OPENAI_API_KEY

# Check in code
fmt.Printf("API Key length: %d\n", len(os.Getenv("OPENAI_API_KEY")))
```

#### 2. Terminal Rendering Issues
```go
// Force redraw
return m, tea.ClearScreen
```

#### 3. Go Module Issues
```bash
cd src
go mod tidy
go mod download
```

### Debugging Tools

1. **Delve Debugger**
   ```bash
   go install github.com/go-delve/delve/cmd/dlv@latest
   dlv debug
   ```

2. **Printf Debugging**
   Add temporary print statements for state inspection

3. **Terminal Debugging**
   Use a second terminal to monitor logs while running the UI

## Contributing

### Code Style

Follow Go conventions:
- Use `gofmt` for formatting
- Follow effective Go guidelines
- Use meaningful variable names
- Add comments for exported functions

### Commit Messages

Use conventional commits:
```
feat: add support for Anthropic Claude models
fix: resolve cursor position bug in text editor
docs: update API documentation for agent package
refactor: simplify command suggestion parsing
```

### Pull Request Process

1. Create a feature branch from `main`
2. Make your changes with tests
3. Update documentation if needed
4. Ensure all tests pass
5. Submit a pull request with clear description

### Adding New Features

#### New LLM Provider
1. Extend `agent.go` to support the new provider
2. Update system prompts if needed
3. Add configuration options
4. Update documentation

#### New UI Components
1. Add to the `ui` package
2. Integrate with the main model
3. Update view rendering
4. Handle new key bindings

#### New Command Types
1. Extend `CommandSuggestion` structure
2. Update parsing logic in `command_suggestion.go`
3. Add execution logic in `executor.go`
4. Update system prompts

### Performance Guidelines

1. **UI Responsiveness**: Keep UI updates under 16ms
2. **Memory Usage**: Be mindful of message history size
3. **API Calls**: Implement rate limiting and retries
4. **File Operations**: Validate file paths and sizes

### Security Considerations

1. **Command Execution**: Always require approval for dangerous commands
2. **File Access**: Validate file paths and permissions
3. **API Keys**: Never log or expose API keys
4. **Input Validation**: Sanitize all user inputs

## Advanced Topics

### Custom System Prompts

You can customize the AI behavior by modifying `system_prompt.go`:

```go
func GetCustomSystemPrompt() string {
    return `You are a specialized assistant for...`
}
```

### Plugin Architecture

The codebase is designed to support plugins:
- Implement the `Provider` interface for new LLM providers
- Add new command types via the suggestion system
- Extend the UI with custom components

### Performance Optimization

- Use Go's profiling tools to identify bottlenecks
- Optimize frequent UI updates
- Consider caching for expensive operations
- Profile memory usage for large conversations

This guide should help you understand and contribute to the Menace CLI project effectively. For more specific questions, please refer to the API documentation or open an issue on the repository.