# Contributing to Menace CLI

Thank you for your interest in contributing to Menace CLI! This guide will help you get started with development and contribution.

## Table of Contents

1. [Development Setup](#development-setup)
2. [Project Architecture](#project-architecture)
3. [Development Workflow](#development-workflow)
4. [Code Style Guidelines](#code-style-guidelines)
5. [Testing](#testing)
6. [Adding Features](#adding-features)
7. [Debugging](#debugging)
8. [Pull Request Process](#pull-request-process)

## Development Setup

### Prerequisites

- **Go 1.18+** with modules enabled
- **Node.js v14+** and npm
- **Git** for version control
- **OpenAI API key** for testing
- **GitHub CLI** (optional, for PR testing)

### Clone and Build

1. **Fork the repository** on GitHub

2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/terminal-chat.git
   cd terminal-chat/menace-cli
   ```

3. **Set up environment variables**:
   ```bash
   export OPENAI_API_KEY="your-test-api-key"
   export ANTHROPIC_API_KEY="your-anthropic-key"  # Optional
   ```

4. **Install dependencies**:
   ```bash
   cd src && go mod tidy
   ```

5. **Build the project**:
   ```bash
   npm run build
   ```

6. **Test the build**:
   ```bash
   ./bin/menace-go-linux  # On Linux
   ./bin/menace-go-darwin # On macOS
   ```

### Development Environment

For development, you can run the Go binary directly without the npm wrapper:

```bash
cd src
go run . 
```

This allows for faster iteration during development.

## Project Architecture

### High-Level Overview

```
menace-cli/
├── package.json           # npm metadata and build scripts
├── bin/                   # Compiled binaries and npm wrapper
│   └── menace.js         # Node.js wrapper script
├── scripts/              # Build automation
│   ├── build-all.js      # Cross-platform build script
│   ├── build-mac.js      # macOS-specific build
│   └── build-windows.js  # Windows-specific build
└── src/                  # Go source code
    ├── main.go           # Entry point
    ├── llmServer/        # LLM integration layer
    └── ui/               # Terminal UI using Bubble Tea
```

### Core Components

#### 1. LLM Server (`src/llmServer/`)

**Purpose**: Handles all AI model interactions and command parsing.

- **`agent.go`**: Core agent with conversation management
- **`command_suggestion.go`**: Parses AI responses for commands
- **`github_interface.go`**: GitHub API integration
- **`llmFactory.go`**: Factory for different LLM providers
- **`system_prompt.go`**: Dynamic system prompt generation

#### 2. UI Layer (`src/ui/`)

**Purpose**: Terminal interface built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

- **`model.go`**: Main UI state management
- **`view.go`**: Rendering logic
- **`update.go`**: Event handling and state updates
- **`executor.go`**: Command execution
- **`message.go`**: Message types and handling

#### 3. Build System

**Purpose**: Cross-platform binary generation and npm packaging.

- **Node.js wrapper**: Platform detection and binary spawning
- **Build scripts**: Automated cross-compilation for Windows, macOS, Linux

### Data Flow

```
User Input → UI Model → Agent → LLM Provider → Response → Command Parser → UI Update
```

1. **User types** in the terminal UI
2. **UI model** captures input and sends to agent
3. **Agent** formats message and sends to LLM
4. **LLM response** is parsed for commands or function calls
5. **UI updates** to show response and handle command execution

## Development Workflow

### Setting Up a Development Branch

```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Make your changes
# ...

# Test your changes
go run . # or npm run build && ./bin/menace

# Commit your changes
git add .
git commit -m "feat: add your feature description"

# Push to your fork
git push origin feature/your-feature-name
```

### Live Development

For rapid development, use:

```bash
cd src
# Run with live reload (you'll need to restart manually)
go run .
```

For UI development specifically, you might want to add debug logging:

```bash
DEBUG=1 go run .
```

### Building and Testing

```bash
# Build all platform binaries
npm run build

# Build specific platforms
npm run build-mac
npm run build-windows

# Test specific platform
./bin/menace-go-linux
```

## Code Style Guidelines

### Go Code Style

Follow standard Go conventions:

1. **Use `gofmt`** to format code:
   ```bash
   gofmt -w src/
   ```

2. **Run `go vet`** for static analysis:
   ```bash
   cd src && go vet ./...
   ```

3. **Use meaningful names**:
   ```go
   // Good
   func (a *Agent) SendMessage(ctx context.Context, input string) (string, *CommandSuggestion, error)
   
   // Avoid
   func (a *Agent) Send(c context.Context, i string) (string, *CS, error)
   ```

4. **Document public functions**:
   ```go
   // SendMessage sends a message to the LLM and returns the response.
   // It parses the response for command suggestions and updates conversation history.
   func (a *Agent) SendMessage(ctx context.Context, input string) (string, *CommandSuggestion, error) {
       // implementation
   }
   ```

### JavaScript Code Style

For build scripts, use:

1. **Node.js best practices**
2. **Consistent formatting**
3. **Error handling**

### UI Code Guidelines

When working with Bubble Tea UI:

1. **Separate concerns**: Keep model, view, and update logic separate
2. **Handle errors gracefully**: Always provide user feedback for errors
3. **Maintain state consistency**: Use mutex for concurrent access
4. **Responsive design**: Handle terminal resizing

## Testing

### Manual Testing

1. **Basic functionality**:
   ```bash
   # Test basic chat
   go run . 
   # Type: "Hello, list files in current directory"
   ```

2. **File operations**:
   ```bash
   # Test file reading
   # Type: "Show me the contents of main.go"
   
   # Test file editing
   # Type: "Add a comment to line 10 in main.go"
   ```

3. **Command execution**:
   ```bash
   # Test command suggestions
   # Type: "Create a new directory called 'test'"
   ```

4. **Model switching**:
   - Press Ctrl+M to open model selection
   - Test with different providers (OpenAI, Anthropic, Ollama)

### Cross-Platform Testing

Test on different platforms:

```bash
# Build for all platforms
npm run build

# Test Windows binary (in WSL or Windows)
./bin/menace-go-win.exe

# Test macOS binaries
./bin/menace-go-darwin        # Intel
./bin/menace-go-darwin-arm64  # Apple Silicon

# Test Linux binary
./bin/menace-go-linux
```

### Integration Testing

Test GitHub integration:

```bash
# Test Git operations
# Type: "What's the current git status?"
# Type: "Create a commit with message 'test commit'"

# Test PR creation (requires GitHub CLI)
# Type: "Create a pull request for this branch"
```

## Adding Features

### Adding New LLM Providers

1. **Update `llmFactory.go`**:
   ```go
   func (mf ModelFactory) CreateLLM(provider string, apiKey string, model string) (llms.Model, error) {
       switch provider {
       case "your_provider":
           return your_provider.New(
               your_provider.WithToken(apiKey),
               your_provider.WithModel(model),
           )
       // ... existing cases
       }
   }
   ```

2. **Update `agent.go` SetModel method**:
   ```go
   case "your_provider":
       llm, err := your_provider.New(
           your_provider.WithToken(os.Getenv("YOUR_PROVIDER_API_KEY")),
           your_provider.WithModel(model),
       )
   ```

3. **Update model selection UI** in `model_selection.go`

### Adding New Function Calls

1. **Define the function** in appropriate package:
   ```go
   func YourNewFunction(arg1 string, arg2 int) error {
       // implementation
   }
   ```

2. **Add to system prompt** in `system_prompt.go`:
   ```go
   // Add documentation for your function call format
   ```

3. **Handle in UI executor** in `model.go`:
   ```go
   case "yourNewFunction":
       arg1, _ := m.PendingFunctionCall.Args["arg1"].(string)
       arg2, _ := m.PendingFunctionCall.Args["arg2"].(float64)
       err = YourNewFunction(arg1, int(arg2))
   ```

### Adding UI Features

1. **Update Model struct** in `model.go`:
   ```go
   type Model struct {
       // existing fields
       YourNewField bool
   }
   ```

2. **Handle input events** in `update.go`:
   ```go
   case tea.KeyMsg:
       switch msg.String() {
       case "your_key":
           m.YourNewField = !m.YourNewField
       }
   ```

3. **Update rendering** in `view.go`:
   ```go
   func (m Model) View() string {
       // Add rendering for your feature
   }
   ```

## Debugging

### Debug Logging

Add debug logging throughout the code:

```go
import "log"

func (a *Agent) SendMessage(ctx context.Context, input string) (string, *CommandSuggestion, error) {
    if os.Getenv("DEBUG") != "" {
        log.Printf("Sending message: %s", input)
    }
    // rest of implementation
}
```

Run with debug logging:

```bash
DEBUG=1 go run .
```

### Common Debugging Scenarios

1. **LLM responses not parsing correctly**:
   - Check system prompt format
   - Verify parsing logic in `command_suggestion.go`
   - Test with different models

2. **UI state issues**:
   - Add logging to update functions
   - Check thread safety with mutexes
   - Verify Bubble Tea message handling

3. **File operations failing**:
   - Check file permissions
   - Verify paths are correct
   - Test with different file types

### Debugging Tools

1. **Go debugging** with Delve:
   ```bash
   go install github.com/go-delve/delve/cmd/dlv@latest
   dlv debug .
   ```

2. **Terminal UI debugging**:
   - Use log files instead of stdout
   - Add debug panels to the UI
   - Test in different terminal types

## Pull Request Process

### Before Submitting

1. **Test thoroughly**:
   - Test on your target platform
   - Test basic functionality
   - Test edge cases

2. **Check code quality**:
   ```bash
   cd src
   gofmt -w .
   go vet ./...
   go mod tidy
   ```

3. **Build all platforms**:
   ```bash
   npm run build
   ```

4. **Update documentation** if needed

### PR Guidelines

1. **Clear title**: Describe what the PR does
   - ✅ "feat: add support for Anthropic Claude models"  
   - ❌ "fix stuff"

2. **Detailed description**:
   - What does this PR do?
   - Why is it needed?
   - How was it tested?
   - Any breaking changes?

3. **Small, focused changes**: One feature per PR when possible

4. **Link issues**: If fixing an issue, reference it: "Fixes #123"

### Example PR Description

```markdown
## Summary

Adds support for Anthropic Claude models as an alternative to OpenAI.

## Changes

- Added Anthropic provider to `llmFactory.go`
- Updated `agent.go` with Claude model support
- Added Claude models to model selection UI
- Updated documentation

## Testing

- Tested with Claude 3 Sonnet and Claude Instant
- Verified model switching works correctly
- Tested on macOS and Linux

## Breaking Changes

None - this is backward compatible.

Fixes #42
```

### Review Process

1. **Automated checks**: Ensure CI passes
2. **Code review**: Address feedback from maintainers
3. **Testing**: Maintainers may test on different platforms
4. **Merge**: After approval, the PR will be merged

---

## Getting Help

- **GitHub Issues**: For bugs and feature requests
- **Discussions**: For questions and general discussion
- **Code Review**: Don't hesitate to ask for feedback early

Thank you for contributing to Menace CLI! 🚀