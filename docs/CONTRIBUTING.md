# Contributing to Menace CLI

Thank you for your interest in contributing to Menace CLI! This guide will help you get started with development, testing, and contributing to the project.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Development Environment Setup](#development-environment-setup)
3. [Project Structure](#project-structure)
4. [Development Workflow](#development-workflow)
5. [Code Style Guidelines](#code-style-guidelines)
6. [Testing](#testing)
7. [Building and Packaging](#building-and-packaging)
8. [Contributing Guidelines](#contributing-guidelines)
9. [Release Process](#release-process)

## Getting Started

### Prerequisites

Ensure you have the following installed:

- **Go 1.18+**: [Download Go](https://golang.org/dl/)
- **Node.js 14+**: [Download Node.js](https://nodejs.org/)
- **Git**: [Download Git](https://git-scm.com/)

### Quick Setup

```bash
# Clone the repository
git clone https://github.com/your-org/menace-cli.git
cd menace-cli

# Set up environment variables
export OPENAI_API_KEY="sk-your-key-here"
export ANTHROPIC_API_KEY="sk-ant-your-key-here"  # Optional
export GITHUB_TOKEN="ghp_your-token-here"        # Optional

# Install Node.js dependencies (if any)
npm install

# Build the project
npm run build

# Test the build
./bin/menace-go-linux  # On Linux
# or
./bin/menace-go-darwin  # On macOS
# or
./bin/menace-go-win.exe  # On Windows
```

## Development Environment Setup

### Go Environment

1. **Install Go modules**:
   ```bash
   cd menace-cli/src
   go mod download
   go mod tidy
   ```

2. **Verify Go installation**:
   ```bash
   go version  # Should show Go 1.18+
   go env      # Check GOPATH and GOROOT
   ```

3. **Install development tools**:
   ```bash
   # Code formatting
   go install golang.org/x/tools/cmd/goimports@latest
   
   # Linting
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # Testing tools
   go install github.com/onsi/ginkgo/v2/ginkgo@latest
   go install github.com/onsi/gomega/...@latest
   ```

### IDE Setup

#### VS Code (Recommended)

Install the following extensions:
- Go (by Google)
- Go Test Explorer
- Error Lens
- GitLens

**VS Code Settings** (`.vscode/settings.json`):
```json
{
    "go.toolsManagement.checkForUpdates": "local",
    "go.useLanguageServer": true,
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.testFlags": ["-v"],
    "go.buildTags": "integration",
    "editor.codeActionsOnSave": {
        "source.organizeImports": true
    }
}
```

#### GoLand

1. Open the `menace-cli` directory
2. Configure Go SDK to your installed version
3. Enable Go modules integration
4. Set up run configurations for main.go

### Local Testing Setup

1. **Set up test environment**:
   ```bash
   # Create test directory
   mkdir -p test-workspace
   cd test-workspace
   
   # Initialize test Git repo
   git init
   echo "# Test Project" > README.md
   git add README.md
   git commit -m "Initial commit"
   ```

2. **Test different scenarios**:
   ```bash
   # Test in different project types
   mkdir node-project && cd node-project && npm init -y
   mkdir go-project && cd go-project && go mod init test
   mkdir python-project && cd python-project && touch requirements.txt
   ```

## Project Structure

### High-Level Architecture

```
menace-cli/
├── bin/                    # Compiled binaries
│   ├── menace.js          # Node.js wrapper
│   ├── menace-go-linux    # Linux binary
│   ├── menace-go-darwin   # macOS binary
│   └── menace-go-win.exe  # Windows binary
├── scripts/               # Build scripts
│   ├── build-all.js      # Cross-platform build
│   ├── build-mac.js      # macOS-specific build
│   └── build-windows.js  # Windows-specific build  
├── src/                   # Go source code
│   ├── main.go           # Application entry point
│   ├── go.mod            # Go module definition
│   ├── llmServer/        # LLM integration layer
│   │   ├── agent.go      # Core agent logic
│   │   ├── command_suggestion.go
│   │   ├── github_interface.go
│   │   ├── llmFactory.go
│   │   ├── llmService.go
│   │   └── system_prompt.go
│   └── ui/               # Terminal UI components
│       ├── model.go      # Bubble Tea model
│       ├── update.go     # Event handling
│       ├── view.go       # Rendering logic
│       ├── executor.go   # Command execution
│       ├── message.go    # Message types
│       ├── model_selection.go
│       ├── thinking.go   # Loading animations
│       ├── ui_utils.go   # Utility functions
│       └── update.go     # Update logic
├── docs/                  # Documentation
├── tests/                 # Test files
├── package.json          # NPM configuration
├── README.md             # Project overview
└── LICENSE               # MIT License
```

### Key Components

#### LLM Server (`src/llmServer/`)

**agent.go**: Core agent implementation
```go
type Agent struct {
    llm      llms.Model              // LLM interface
    mu       sync.Mutex              // Thread safety
    shell    string                  // Shell detection
    messages []llms.MessageContent   // Context
    ctx      context.Context         // Request context
    provider string                  // Provider name
    Model    string                  // Model name
}
```

**Key responsibilities**:
- LLM provider abstraction
- Context management
- Request/response handling
- Model switching

#### UI Package (`src/ui/`)

**model.go**: UI state management
```go
type Model struct {
    Input    string           // User input
    Messages []Message        // Chat history
    agent    *llmServer.Agent // LLM agent
    // ... UI state fields
}
```

**update.go**: Event processing using Bubble Tea's Elm architecture
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyInput(msg)
    case LLMResponseMsg:
        return m.handleLLMResponse(msg)
    // ... other message types
    }
}
```

## Development Workflow

### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feat/your-feature-name
   ```

2. **Make your changes** following the code style guidelines

3. **Test your changes**:
   ```bash
   # Run tests
   cd src && go test ./...
   
   # Build and test locally
   go build -o ../bin/menace-go-dev .
   ../bin/menace-go-dev
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

### Testing Changes

#### Unit Testing

```bash
cd src
go test ./llmServer -v
go test ./ui -v
```

#### Integration Testing

```bash
# Test with real LLM providers
OPENAI_API_KEY="sk-..." go test ./... -tags integration

# Test cross-platform builds
npm run build
```

#### Manual Testing

1. **Test basic functionality**:
   ```bash
   ./bin/menace-go-linux
   > Hello, how are you?
   > What files are in this directory?
   > Help me create a simple script
   ```

2. **Test model switching**:
   - Press Tab to open model selection
   - Try different providers (OpenAI, Anthropic, Ollama)
   - Verify functionality works across models

3. **Test command execution**:
   ```bash
   > List all files in this directory
   # Should show a suggested command
   # Test approving and rejecting commands
   ```

4. **Test file operations**:
   ```bash
   > Create a new file called test.txt with some content
   > Read the contents of test.txt
   > Add a line to test.txt
   ```

### Debugging

#### Debug Mode

```bash
DEBUG=1 ./bin/menace-go-linux
```

#### Go Debugging with Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the application
cd src
dlv debug . -- 
```

#### Common Debug Points

1. **LLM requests/responses**:
   ```go
   // In agent.go
   fmt.Printf("DEBUG: Sending request: %+v\n", request)
   fmt.Printf("DEBUG: Received response: %+v\n", response)
   ```

2. **UI state changes**:
   ```go
   // In update.go
   fmt.Printf("DEBUG: Processing message: %T %+v\n", msg, msg)
   ```

3. **Command execution**:
   ```go
   // In executor.go
   fmt.Printf("DEBUG: Executing command: %s\n", command)
   ```

## Code Style Guidelines

### Go Code Style

Follow standard Go conventions:

1. **Use gofmt/goimports**:
   ```bash
   goimports -w .
   ```

2. **Follow naming conventions**:
   ```go
   // Good
   type UserMessage struct {
       Content   string
       Timestamp time.Time
   }
   
   // Bad
   type user_message struct {
       content   string
       timestamp time.Time
   }
   ```

3. **Add comments for exported functions**:
   ```go
   // ProcessQuery sends a query to the LLM and returns the response.
   // It manages context, handles errors, and parses structured responses.
   func (a *Agent) ProcessQuery(query string) (string, error) {
       // Implementation
   }
   ```

4. **Handle errors properly**:
   ```go
   // Good
   result, err := someOperation()
   if err != nil {
       return fmt.Errorf("operation failed: %w", err)
   }
   
   // Bad
   result, _ := someOperation()
   ```

5. **Use meaningful variable names**:
   ```go
   // Good
   func (a *Agent) SetModel(providerName, modelName string) error
   
   // Bad  
   func (a *Agent) SetModel(p, m string) error
   ```

### Package Organization

1. **Keep packages focused**:
   - `llmServer`: LLM-related functionality only
   - `ui`: UI components and state management only

2. **Minimize inter-package dependencies**:
   - `ui` can import `llmServer`
   - `llmServer` should not import `ui`
   - Both can import standard library and external packages

3. **Use internal packages for shared utilities**:
   ```go
   // Create src/internal/utils for shared code
   package utils
   
   func FormatMessage(msg string) string {
       // Shared formatting logic
   }
   ```

### Error Handling

1. **Use wrapped errors**:
   ```go
   if err != nil {
       return fmt.Errorf("failed to process query: %w", err)
   }
   ```

2. **Define custom error types for recoverable errors**:
   ```go
   type ModelNotAvailableError struct {
       Provider string
       Model    string
   }
   
   func (e *ModelNotAvailableError) Error() string {
       return fmt.Sprintf("model %s not available from provider %s", e.Model, e.Provider)
   }
   ```

3. **Handle errors gracefully in UI**:
   ```go
   case LLMResponseMsg:
       if msg.Error != nil {
           m.Messages = append(m.Messages, Message{
               Role:    "system",
               Content: fmt.Sprintf("Error: %v", msg.Error),
           })
           return m, nil
       }
   ```

## Testing

### Test Structure

Create tests alongside source files:
```
src/
├── llmServer/
│   ├── agent.go
│   ├── agent_test.go
│   ├── command_suggestion.go
│   └── command_suggestion_test.go
└── ui/
    ├── model.go
    ├── model_test.go
    ├── update.go
    └── update_test.go
```

### Unit Test Examples

#### Testing LLM Agent

```go
// src/llmServer/agent_test.go
package llmServer

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestAgent_ProcessQuery(t *testing.T) {
    tests := []struct {
        name          string
        query         string
        mockResponse  string
        expectedError bool
    }{
        {
            name:         "successful query",
            query:        "Hello",
            mockResponse: "Hi there!",
            expectedError: false,
        },
        {
            name:         "empty query",
            query:        "",
            expectedError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

#### Testing UI Components

```go
// src/ui/model_test.go
package ui

import (
    "testing"
    "github.com/charmbracelet/bubbletea"
    "github.com/stretchr/testify/assert"
)

func TestModel_HandleKeyInput(t *testing.T) {
    model := NewModel(nil)
    
    // Test input handling
    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hello")}
    updatedModel, _ := model.Update(msg)
    
    assert.Equal(t, "hello", updatedModel.(Model).Input)
}
```

### Integration Tests

```go
// tests/integration_test.go
//go:build integration

package tests

import (
    "os"
    "testing"
    "menace-go/llmServer"
)

func TestRealLLMIntegration(t *testing.T) {
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        t.Skip("OPENAI_API_KEY not set")
    }
    
    agent, err := llmServer.NewAgent(apiKey)
    assert.NoError(t, err)
    
    response, err := agent.ProcessQuery("Hello")
    assert.NoError(t, err)
    assert.NotEmpty(t, response)
}
```

### Running Tests

```bash
# Unit tests only
go test ./...

# With integration tests
go test ./... -tags integration

# With coverage
go test ./... -cover

# Verbose output
go test ./... -v

# Specific package
go test ./llmServer -v
```

## Building and Packaging

### Local Development Build

```bash
# Build for current platform
cd src
go build -o ../bin/menace-go-dev .

# Test the build
../bin/menace-go-dev
```

### Cross-Platform Build

```bash
# Build for all platforms
npm run build

# Build for specific platforms
npm run build-windows
npm run build-mac

# Manual cross-compilation
cd src
GOOS=linux GOARCH=amd64 go build -o ../bin/menace-go-linux .
GOOS=darwin GOARCH=amd64 go build -o ../bin/menace-go-darwin .
GOOS=darwin GOARCH=arm64 go build -o ../bin/menace-go-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o ../bin/menace-go-win.exe .
```

### Build Optimization

For release builds, use optimization flags:

```bash
go build -ldflags="-s -w" -o menace-go .
```

Flags explanation:
- `-s`: Remove symbol table and debug info
- `-w`: Remove DWARF debug info
- Results in smaller binary size

### Packaging for Distribution

```bash
# Create distribution archive
npm pack

# Test the package
npm install -g menace-cli-1.0.0.tgz
menace --version
```

## Contributing Guidelines

### Issue Reporting

When reporting issues, please include:

1. **Environment information**:
   - Operating system and version
   - Go version (`go version`)
   - Node.js version (`node --version`)
   - Terminal emulator

2. **Reproduction steps**:
   - Exact commands used
   - Expected vs. actual behavior
   - Error messages (full text)

3. **Context**:
   - What were you trying to accomplish?
   - Is this a regression?
   - Does it happen consistently?

### Pull Request Process

1. **Fork the repository** and create a feature branch
2. **Make your changes** following the style guidelines
3. **Add tests** for new functionality
4. **Update documentation** if needed
5. **Test thoroughly** including cross-platform compatibility
6. **Submit pull request** with clear description

#### Pull Request Template

```markdown
## Description
Brief description of changes and motivation.

## Type of Change
- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to change)
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass  
- [ ] Manual testing completed
- [ ] Cross-platform testing (if applicable)

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] No new warnings introduced
```

### Code Review Guidelines

#### For Reviewers

1. **Check functionality**: Does the code do what it's supposed to do?
2. **Review tests**: Are there adequate tests? Do they pass?
3. **Code quality**: Is the code readable and maintainable?
4. **Performance**: Are there any performance concerns?
5. **Security**: Any security implications?
6. **Documentation**: Is documentation updated?

#### For Contributors

1. **Respond promptly** to review feedback
2. **Be open to suggestions** and different approaches
3. **Test suggested changes** before pushing updates
4. **Update PR description** if scope changes significantly

### Communication

- **GitHub Issues**: Bug reports, feature requests
- **GitHub Discussions**: General questions, ideas
- **Pull Requests**: Code review and discussion
- **Discord/Slack**: Real-time chat (if available)

## Release Process

### Versioning

Menace CLI follows [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

1. **Update version** in `package.json`
2. **Update CHANGELOG.md** with new features/fixes
3. **Test release candidate**:
   ```bash
   npm run build
   # Test all platforms
   ```
4. **Create release tag**:
   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```
5. **Build release artifacts**
6. **Publish to npm**:
   ```bash
   npm publish
   ```
7. **Create GitHub release** with release notes

### Pre-release Testing

Before releasing:

1. **Cross-platform testing** on Windows, macOS, Linux
2. **Integration testing** with all supported LLM providers
3. **Performance testing** with large conversations
4. **Backward compatibility** testing
5. **Documentation review** and updates

---

## Getting Help

- **Documentation**: Check the `docs/` directory
- **GitHub Issues**: Search existing issues or create new ones
- **Code Review**: Request reviews from maintainers
- **Community**: Join discussions with other contributors

Thank you for contributing to Menace CLI! Your contributions help make terminal interaction more intelligent and productive for developers worldwide.