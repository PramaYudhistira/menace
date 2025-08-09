# Contributing to Menace CLI

Thank you for your interest in contributing to Menace CLI! This guide will help you get started with contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Contribution Types](#contribution-types)
- [Code Style](#code-style)
- [Testing](#testing)
- [Documentation](#documentation)
- [Submitting Changes](#submitting-changes)
- [Review Process](#review-process)

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:

- **Be respectful**: Treat all contributors with respect and kindness
- **Be inclusive**: Welcome newcomers and help them get started
- **Be constructive**: Provide helpful feedback and suggestions
- **Be collaborative**: Work together to improve the project

## Getting Started

### Prerequisites

- **Go 1.18+** installed and configured
- **Node.js v14+** and npm
- **Git** for version control
- A supported terminal emulator
- API key for testing (OpenAI/Anthropic/Ollama)

### Development Setup

1. **Fork the Repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/menace-cli.git
   cd menace-cli/menace-cli
   ```

2. **Set Up Development Environment**
   ```bash
   # Install Node.js dependencies
   npm install
   
   # Set up Go modules
   cd src
   go mod download
   ```

3. **Configure Environment**
   ```bash
   # Set up API keys for testing
   export OPENAI_API_KEY="your-test-key-here"
   export MENACE_LOG_LEVEL="debug"  # Enable debug logging
   ```

4. **Verify Setup**
   ```bash
   # Test build
   go build .
   
   # Test run (should start the UI)
   go run .
   ```

### Project Structure

Understanding the codebase structure will help you contribute effectively:

```
menace-cli/
├── src/
│   ├── main.go            # Application entry point
│   ├── llmServer/         # AI and LLM functionality
│   │   ├── agent.go       # Main AI agent
│   │   ├── command_suggestion.go  # Command parsing
│   │   ├── github_interface.go    # GitHub operations
│   │   ├── llmFactory.go  # Environment detection
│   │   ├── llmService.go  # Legacy LLM service
│   │   └── system_prompt.go       # AI prompts
│   └── ui/                # Terminal user interface
│       ├── model.go       # UI state management
│       ├── view.go        # UI rendering
│       ├── executor.go    # Command execution
│       ├── message.go     # Message structures
│       └── ...
├── scripts/               # Build scripts
├── docs/                  # Documentation
└── package.json           # Build configuration
```

## Development Workflow

### Branching Strategy

We use a feature branch workflow:

1. **Main Branch**: `main` - stable, production-ready code
2. **Feature Branches**: `feature/description` - new features
3. **Bug Fix Branches**: `fix/description` - bug fixes
4. **Documentation Branches**: `docs/description` - documentation updates

### Creating a Feature Branch

```bash
# Start from the latest main
git checkout main
git pull upstream main

# Create and switch to feature branch
git checkout -b feature/your-feature-name

# Make your changes...
git add .
git commit -m "feat: add your feature description"

# Push to your fork
git push origin feature/your-feature-name
```

### Keeping Your Fork Updated

```bash
# Add upstream remote (one time setup)
git remote add upstream https://github.com/original-org/menace-cli.git

# Keep your main branch updated
git checkout main
git pull upstream main
git push origin main

# Rebase your feature branch
git checkout feature/your-feature-name
git rebase main
```

## Contribution Types

### 🐛 Bug Fixes

1. **Search Existing Issues**: Check if the bug is already reported
2. **Create Issue**: If not found, create a detailed bug report
3. **Fix the Bug**: Create a fix with tests
4. **Test Thoroughly**: Ensure the fix works and doesn't break anything else

**Bug Fix Checklist:**
- [ ] Issue clearly describes the problem
- [ ] Fix addresses the root cause
- [ ] Tests added to prevent regression
- [ ] Documentation updated if needed

### ✨ New Features

1. **Discuss First**: Open an issue to discuss the feature
2. **Design Review**: Get feedback on the approach
3. **Implement**: Code the feature with tests
4. **Document**: Update relevant documentation

**Feature Checklist:**
- [ ] Feature discussed and approved
- [ ] Implementation follows project patterns
- [ ] Tests cover new functionality
- [ ] Documentation updated
- [ ] Breaking changes noted

### 📚 Documentation

Documentation improvements are always welcome:
- Fix typos and grammar
- Add examples and use cases
- Improve clarity and organization
- Add missing documentation

### 🚀 Performance Improvements

Performance contributions should:
- Include benchmarks showing improvement
- Not sacrifice code readability
- Be thoroughly tested
- Consider memory usage impact

## Code Style

### Go Code Style

We follow standard Go conventions:

```go
// Good: Clear function names with proper comments
// ProcessUserMessage handles user input and returns AI response
func ProcessUserMessage(input string) (string, error) {
    if input == "" {
        return "", ErrEmptyInput
    }
    
    // Process the input...
    return response, nil
}

// Good: Proper error handling
result, err := someFunction()
if err != nil {
    return fmt.Errorf("failed to process: %w", err)
}

// Good: Consistent struct formatting
type Config struct {
    APIKey      string `json:"api_key"`
    Model       string `json:"model"`
    MaxTokens   int    `json:"max_tokens"`
}
```

### Code Formatting

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# Vet code
go vet ./...
```

### Naming Conventions

- **Functions**: `CamelCase` for exported, `camelCase` for private
- **Variables**: `camelCase` (short names for short scopes)
- **Constants**: `CamelCase` or `UPPER_CASE` for package-level
- **Interfaces**: `-er` suffix (e.g., `Reader`, `Writer`)

### Comments

```go
// Package comment describes the package purpose
package llmServer

// Exported function comments start with the function name
// NewAgent creates a new AI agent with the specified configuration.
// It returns an error if the configuration is invalid.
func NewAgent(config Config) (*Agent, error) {
    // Internal comments explain complex logic
    // Initialize the LLM chain with memory management
    return &Agent{}, nil
}
```

## Testing

### Unit Tests

Write tests for all new functionality:

```go
func TestNewAgent(t *testing.T) {
    tests := []struct {
        name    string
        config  Config
        wantErr bool
    }{
        {
            name:   "valid config",
            config: Config{APIKey: "test-key"},
            wantErr: false,
        },
        {
            name:   "empty API key",
            config: Config{},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewAgent(tt.config)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewAgent() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestNewAgent ./llmServer

# Run tests with verbose output
go test -v ./...
```

### Integration Tests

For features that interact with external services:

```go
func TestIntegrationWithOpenAI(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        t.Skip("OPENAI_API_KEY not set")
    }
    
    // Test actual API integration...
}
```

### Manual Testing

Before submitting:
- [ ] Test on your target platform
- [ ] Test with different terminal emulators
- [ ] Test edge cases and error conditions
- [ ] Test UI responsiveness and interactions

## Documentation

### Code Documentation

- Document all exported functions and types
- Include usage examples for complex functions
- Document expected behavior and edge cases
- Use proper Go doc format

### User Documentation

When adding features that affect users:
- Update README.md if needed
- Add configuration examples
- Update the relevant documentation files
- Include troubleshooting information

### API Documentation

For changes to the API:
- Update `docs/api/` documentation
- Include examples of usage
- Document breaking changes
- Update version information

## Submitting Changes

### Commit Messages

Use conventional commits format:

```
type(scope): description

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```bash
feat(ui): add syntax highlighting to code blocks

fix(llmServer): resolve memory leak in message processing

docs(api): update configuration examples

refactor(ui): simplify text selection logic
```

### Pull Request Process

1. **Prepare Your PR**
   ```bash
   # Ensure your branch is up to date
   git checkout main
   git pull upstream main
   git checkout your-feature-branch
   git rebase main
   
   # Run tests and formatting
   go test ./...
   go fmt ./...
   ```

2. **Create Pull Request**
   - Use the PR template
   - Provide clear description
   - Link related issues
   - Include screenshots for UI changes

3. **PR Template**
   ```markdown
   ## Description
   Brief description of changes
   
   ## Type of Change
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Documentation update
   - [ ] Performance improvement
   
   ## Testing
   - [ ] Unit tests added/updated
   - [ ] Integration tests pass
   - [ ] Manual testing completed
   
   ## Checklist
   - [ ] Code follows project style
   - [ ] Self-review completed
   - [ ] Documentation updated
   - [ ] No breaking changes (or documented)
   ```

### PR Requirements

- [ ] All tests pass
- [ ] Code is properly formatted
- [ ] Documentation is updated
- [ ] Commit messages follow conventions
- [ ] No merge conflicts
- [ ] PR description is complete

## Review Process

### What Reviewers Look For

- **Correctness**: Does the code do what it's supposed to do?
- **Style**: Does it follow project conventions?
- **Performance**: Are there any performance concerns?
- **Security**: Are there any security implications?
- **Testing**: Is the code adequately tested?
- **Documentation**: Is it properly documented?

### Responding to Reviews

- **Be Responsive**: Reply to comments promptly
- **Be Open**: Consider feedback objectively
- **Ask Questions**: If something isn't clear, ask
- **Make Changes**: Address reviewer concerns
- **Test Again**: After making changes, test thoroughly

### Review Timeline

- **Initial Review**: Within 2-3 business days
- **Follow-up Reviews**: Within 1-2 business days
- **Merge**: After approval and passing CI

## Getting Help

### Where to Ask Questions

1. **GitHub Discussions**: General questions and ideas
2. **GitHub Issues**: Bug reports and feature requests
3. **PR Comments**: Questions about specific code
4. **Discord/Slack**: Real-time community chat (if available)

### Common Questions

**Q: How do I test with different LLM providers?**
A: Set the appropriate environment variables and create test configurations for each provider.

**Q: My tests are failing in CI but passing locally.**
A: Check for platform-specific issues and ensure you're using the same Go version.

**Q: How do I add a new command type?**
A: Extend the `CommandSuggestion` type and update the parsing logic in `command_suggestion.go`.

## Recognition

Contributors who make significant contributions may be:
- Added to the project's contributor list
- Mentioned in release notes
- Given commit access (for long-term contributors)

Thank you for contributing to Menace CLI! 🚀