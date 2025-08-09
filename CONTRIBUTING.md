# Contributing to Menace CLI

Thank you for your interest in contributing to Menace CLI! This document provides guidelines and information for contributors.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Code Standards](#code-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Release Process](#release-process)

## Getting Started

### Prerequisites

Before contributing, ensure you have:

- **Go 1.18+** - For building the core application
- **Node.js 14+** - For build scripts and npm packaging
- **Git** - For version control
- **Make** (optional) - For convenience commands

### Development Environment Setup

1. **Fork the repository**:
   ```bash
   # Fork via GitHub UI, then clone your fork
   git clone https://github.com/your-username/menace-cli.git
   cd menace-cli
   ```

2. **Install dependencies**:
   ```bash
   cd menace-cli/src
   go mod download
   
   # For Node.js dependencies (build scripts)
   cd ..
   npm install --only=dev
   ```

3. **Set up environment variables**:
   ```bash
   export OPENAI_API_KEY="your_test_api_key"
   export ANTHROPIC_API_KEY="your_anthropic_key"  # Optional
   ```

4. **Build and test**:
   ```bash
   npm run build
   npm test  # When tests are available
   ```

## Development Setup

### Local Development Workflow

1. **Run from source** (fastest for development):
   ```bash
   cd menace-cli/src
   go run . 
   ```

2. **Build specific platform**:
   ```bash
   npm run build-linux    # Linux binary
   npm run build-mac      # macOS binaries  
   npm run build-windows  # Windows binary
   ```

3. **Build all platforms**:
   ```bash
   npm run build
   ```

### Project Structure Understanding

```
menace-cli/
├── docs/                  # Documentation (API, User Guide, etc.)
├── src/                   # Go source code
│   ├── main.go           # Entry point
│   ├── llmServer/        # AI and command processing
│   └── ui/               # User interface components
├── scripts/              # Build automation
├── bin/                  # Compiled binaries
└── package.json          # npm configuration
```

### Key Files to Understand

- `src/main.go` - Application bootstrap and initialization
- `src/ui/model.go` - Core UI state management
- `src/ui/update.go` - Event handling and state transitions
- `src/ui/view.go` - Rendering and layout
- `src/llmServer/agent.go` - LLM interaction logic
- `src/llmServer/system_prompt.go` - AI behavior configuration

## Code Standards

### Go Code Style

Follow standard Go conventions:

```go
// ✅ Good: Descriptive function names with documentation
// SendMessage sends a user message to the LLM and returns the response
// along with any parsed command suggestions.
func (a *Agent) SendMessage(ctx context.Context, input string) (string, *CommandSuggestion, error) {
    // Implementation
}

// ✅ Good: Proper error handling
if err != nil {
    return "", nil, fmt.Errorf("failed to get response from LLM: %v", err)
}

// ✅ Good: Consistent variable naming
var (
    responseText    string
    cmdSuggestion  *CommandSuggestion
    messageHistory []llms.MessageContent
)
```

### Documentation Standards

- **Public functions**: Must have descriptive comments
- **Complex logic**: Explain the "why" not just the "what"
- **Examples**: Provide usage examples for non-trivial functions

```go
// ✅ Good documentation
// ParseCommandSuggestion extracts command suggestions from LLM responses.
// It looks for [COMMAND_SUGGESTION] blocks and parses the structured content.
//
// Example input:
//   "[COMMAND_SUGGESTION]\nCommand: ls -la\nReason: List files\n[/COMMAND_SUGGESTION]"
//
// Returns nil if no command suggestion is found.
func parseCommandSuggestion(response string) *CommandSuggestion {
    // Implementation
}
```

### UI Code Guidelines

- **State management**: Keep state in the Model struct
- **Pure functions**: View functions should be pure (no side effects)
- **Responsive design**: Handle different terminal sizes
- **Accessibility**: Support keyboard navigation

```go
// ✅ Good: Pure view function
func (m Model) renderChatHistory(width, height int) string {
    // Returns rendered string without modifying model state
}

// ✅ Good: State update pattern
func (m *Model) handleUserInput(char string) {
    m.Input += char
    m.CursorX++
}
```

### JavaScript/Node.js Standards

For build scripts, follow these conventions:

```javascript
// ✅ Good: Clear build configuration
const buildTargets = [
    { os: "linux", arch: "amd64", output: "menace-go-linux" },
    { os: "darwin", arch: "amd64", output: "menace-go-darwin" },
    // ...
];

// ✅ Good: Error handling
try {
    execSync(`go build -o "${outPath}" .`, { stdio: "inherit" });
    console.log(`✅ Built ${output}`);
} catch (err) {
    console.error(`❌ Build failed for ${GOOS}/${GOARCH}`);
    process.exit(1);
}
```

## Testing

### Unit Testing Strategy

Create tests for:

1. **Pure functions**: Command parsing, utility functions
2. **State management**: Model state transitions
3. **Business logic**: Agent message handling

```go
// Example test structure
func TestParseCommandSuggestion(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected *CommandSuggestion
    }{
        {
            name:  "valid command suggestion",
            input: "[COMMAND_SUGGESTION]\nCommand: ls\nReason: List files\n[/COMMAND_SUGGESTION]",
            expected: &CommandSuggestion{Command: "ls", Reason: "List files"},
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := parseCommandSuggestion(tt.input)
            // Assert expectations
        })
    }
}
```

### Integration Testing

Test the following integration points:

- LLM provider switching
- File operations (read/write/edit)
- Command execution workflows
- GitHub integration

### Manual Testing Checklist

Before submitting PRs, test:

- [ ] Application starts without errors
- [ ] Basic chat functionality works
- [ ] Model switching works
- [ ] File operations complete successfully
- [ ] Command suggestions parse correctly
- [ ] Cross-platform compatibility (if possible)
- [ ] Terminal resizing handles gracefully
- [ ] Keyboard shortcuts function properly

## Pull Request Process

### Before Creating a PR

1. **Update documentation**: If you changed APIs or added features
2. **Add tests**: For new functionality
3. **Run quality checks**:
   ```bash
   go fmt ./...
   go vet ./...
   go test ./...
   ```
4. **Test the build**:
   ```bash
   npm run build
   ./bin/menace  # Test the built binary
   ```

### PR Description Template

```markdown
## Summary
Brief description of the changes and motivation.

## Changes
- [ ] Feature A: Description
- [ ] Bug fix B: Description  
- [ ] Documentation updates

## Testing
- [ ] Added unit tests
- [ ] Manual testing completed
- [ ] Cross-platform testing (specify platforms)

## Breaking Changes
Describe any breaking changes and migration path.

## Screenshots/Demo
For UI changes, include screenshots or GIFs.
```

### Review Process

1. **Automated checks**: Must pass CI/CD pipeline
2. **Code review**: At least one maintainer approval
3. **Manual testing**: Core functionality verification
4. **Documentation**: Updated if needed

## Issue Guidelines

### Bug Reports

Use this template for bug reports:

```markdown
## Bug Description
Clear description of the bug.

## Steps to Reproduce
1. Step one
2. Step two
3. Step three

## Expected Behavior
What should happen.

## Actual Behavior
What actually happens.

## Environment
- OS: [e.g., macOS 12.0, Ubuntu 20.04, Windows 11]
- Terminal: [e.g., iTerm2, Windows Terminal]
- Menace CLI version: [e.g., 1.0.0]
- Go version: [e.g., 1.19.0]
- Node.js version: [e.g., 16.14.0]

## Additional Context
Logs, screenshots, or other relevant information.
```

### Feature Requests

Use this template for feature requests:

```markdown
## Feature Summary
Brief description of the requested feature.

## Motivation
Why is this feature needed? What problem does it solve?

## Detailed Description
Detailed explanation of the feature and how it should work.

## Possible Implementation
If you have ideas about implementation approach.

## Alternatives Considered
Other approaches you've considered and why they weren't suitable.
```

### Labels and Prioritization

We use these labels:
- `bug` - Something isn't working
- `enhancement` - New feature or request  
- `documentation` - Improvements to docs
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention needed
- `priority: high/medium/low` - Priority levels

## Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features, backward compatible
- **PATCH**: Bug fixes, backward compatible

### Release Checklist

For maintainers releasing new versions:

1. **Update version numbers**:
   - `package.json`
   - Any Go version constants
   - Documentation references

2. **Update changelog**:
   - Document all changes since last release
   - Credit contributors

3. **Create release tag**:
   ```bash
   git tag -a v1.2.3 -m "Release version 1.2.3"
   git push origin v1.2.3
   ```

4. **Build and publish**:
   - Automated via CI/CD
   - Manual verification of published artifacts

## Development Tips

### Debugging Techniques

1. **Add debug logging**:
   ```go
   fmt.Printf("[DEBUG] Current state: %+v\n", m)
   ```

2. **Use Go's built-in tools**:
   ```bash
   go run -race .          # Race condition detection
   go run . -v             # Verbose output
   dlv debug               # Delve debugger
   ```

3. **UI debugging**:
   - Test with different terminal sizes
   - Verify keyboard shortcut handling
   - Check text rendering edge cases

### Common Pitfalls

1. **Concurrency issues**: Always use mutexes when accessing shared state
2. **Unicode handling**: Use `[]rune` for text manipulation, not `[]byte`
3. **Cross-platform paths**: Use `filepath` package, not string manipulation
4. **Terminal compatibility**: Test on multiple terminal emulators

### Useful Resources

- [Effective Go](https://golang.org/doc/effective_go)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [LangChain Go Documentation](https://github.com/tmc/langchaingo)

## Community Guidelines

### Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help newcomers get started
- Maintain professional communication

### Getting Help

- **GitHub Discussions**: For questions and general discussions
- **Issues**: For bug reports and feature requests
- **Discord/Slack**: Community chat (if available)

### Recognition

Contributors are recognized in:
- CHANGELOG.md for significant contributions
- GitHub contributors page
- Special recognition for first-time contributors

---

Thank you for contributing to Menace CLI! Your efforts help make this tool better for everyone.