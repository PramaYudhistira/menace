# Extending Menace CLI

This guide explains how to extend Menace CLI with new LLM providers, functions, and capabilities.

## Adding New LLM Providers

Menace CLI uses [LangChainGo](https://github.com/tmc/langchaingo) for LLM integration, making it relatively straightforward to add new providers.

### Step 1: Update Dependencies

Add the new provider to your Go module:

```bash
go get github.com/tmc/langchaingo/llms/newprovider
```

### Step 2: Import the Provider

In `src/llmServer/agent.go`, add the import:

```go
import (
    // ... existing imports
    "github.com/tmc/langchaingo/llms/newprovider"
)
```

### Step 3: Extend SetModel Function

Add your provider to the `SetModel` function switch statement:

```go
func (a *Agent) SetModel(provider string, model string, openSource bool) error {
    a.mu.Lock()
    defer a.mu.Unlock()

    a.provider = provider
    a.Model = model
    a.isOpenSource = openSource

    switch provider {
    // ... existing cases
    case "newprovider":
        llm, err := newprovider.New(
            newprovider.WithToken(os.Getenv("NEWPROVIDER_API_KEY")),
            newprovider.WithModel(model),
            // Add any provider-specific options
        )
        if err != nil {
            return fmt.Errorf("failed to create NewProvider client with model %s: %v", model, err)
        }
        a.llm = llm
        return nil
    default:
        return fmt.Errorf("unknown provider: %s", provider)
    }
}
```

### Step 4: Update UI Model Selection

In `src/ui/model_selection.go`, add your provider to the available options:

```go
var availableProviders = []ProviderOption{
    // ... existing providers
    {
        Name:        "newprovider",
        DisplayName: "New Provider",
        Models: []ModelOption{
            {Name: "model-1", DisplayName: "New Provider Model 1"},
            {Name: "model-2", DisplayName: "New Provider Model 2"},
        },
        RequiresKey: true,
        KeyEnvVar:   "NEWPROVIDER_API_KEY",
    },
}
```

### Step 5: Update Documentation

Update the README and API documentation to include:
- Environment variable requirements
- Supported models
- Configuration examples
- Any provider-specific limitations

### Example: Adding Cohere Support

```go
// In agent.go
case "cohere":
    llm, err := cohere.New(
        cohere.WithToken(os.Getenv("COHERE_API_KEY")),
        cohere.WithModel(model),
    )
    if err != nil {
        return fmt.Errorf("failed to create Cohere client with model %s: %v", model, err)
    }
    a.llm = llm
    return nil
```

## Adding New Function Calls

Function calls enable the LLM to perform structured operations beyond simple text responses.

### Step 1: Define Function Structure

Create a new function definition in the appropriate module. For file operations, add to `src/ui/executor.go`:

```go
type NewFunctionArgs struct {
    Parameter1 string `json:"parameter1"`
    Parameter2 int    `json:"parameter2"`
}

func (m *Model) executeNewFunction(args NewFunctionArgs) (string, error) {
    // Implementation here
    result := fmt.Sprintf("Executed with %s and %d", args.Parameter1, args.Parameter2)
    return result, nil
}
```

### Step 2: Register Function Handler

In the function execution dispatcher (typically in `update.go` or `executor.go`):

```go
func (m *Model) executeFunctionCall(functionName string, args map[string]interface{}) (string, error) {
    switch functionName {
    // ... existing cases
    case "NewFunction":
        var newArgs NewFunctionArgs
        if err := mapToStruct(args, &newArgs); err != nil {
            return "", fmt.Errorf("invalid arguments for NewFunction: %v", err)
        }
        return m.executeNewFunction(newArgs)
    default:
        return "", fmt.Errorf("unknown function: %s", functionName)
    }
}
```

### Step 3: Update System Prompt

Add documentation for the new function in `src/llmServer/system_prompt.go`:

```go
// Add to the system prompt template
`
Example for new function:
[FUNCTION_CALL]
Reason: Explain why this new function is needed
AwaitingCommandApproval: true/false
Payload:
{
    "name": "NewFunction",
    "args": {
        "parameter1": "value1",
        "parameter2": 42
    }
}
[/FUNCTION_CALL]
`
```

### Step 4: Add Error Handling

Ensure proper error handling and user feedback:

```go
func (m *Model) executeNewFunction(args NewFunctionArgs) (string, error) {
    // Validate arguments
    if args.Parameter1 == "" {
        return "", fmt.Errorf("parameter1 cannot be empty")
    }

    // Perform operation with error handling
    result, err := performOperation(args.Parameter1, args.Parameter2)
    if err != nil {
        return "", fmt.Errorf("failed to execute new function: %v", err)
    }

    return fmt.Sprintf("Success: %s", result), nil
}
```

## Adding UI Components

### Step 1: Create Component

Create new UI components using Bubble Tea patterns:

```go
type NewComponent struct {
    width  int
    height int
    // ... component state
}

func NewNewComponent() NewComponent {
    return NewComponent{
        // Initialize state
    }
}

func (c NewComponent) Update(msg tea.Msg) (NewComponent, tea.Cmd) {
    // Handle events
    return c, nil
}

func (c NewComponent) View() string {
    // Render component
    return "New Component View"
}
```

### Step 2: Integrate with Main Model

Add the component to the main UI model:

```go
type Model struct {
    // ... existing fields
    newComponent NewComponent
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    
    // Update component
    m.newComponent, cmd = m.newComponent.Update(msg)
    
    return m, cmd
}

func (m Model) View() string {
    // Include component in view
    return lipgloss.JoinVertical(
        lipgloss.Top,
        m.existingView(),
        m.newComponent.View(),
    )
}
```

## Configuration Extensions

### Adding New Configuration Options

1. **Environment Variables**: Follow the pattern `MENACE_OPTION_NAME`
2. **Config File Support**: Consider adding YAML/JSON configuration
3. **Runtime Configuration**: Add UI toggles for user preferences

### Example: Adding Debug Mode

```go
// In main.go or appropriate config module
func loadConfig() Config {
    return Config{
        Debug:     os.Getenv("MENACE_DEBUG") == "true",
        LogLevel:  os.Getenv("MENACE_LOG_LEVEL"),
        // ... other options
    }
}
```

## Custom Command Parsing

### Adding New Command Formats

Extend the command parsing system to support new structured formats:

```go
func parseCustomCommand(response string) *CustomCommand {
    // Define regex or parsing logic for new format
    pattern := regexp.MustCompile(`\[CUSTOM_COMMAND\](.*?)\[/CUSTOM_COMMAND\]`)
    matches := pattern.FindStringSubmatch(response)
    
    if len(matches) < 2 {
        return nil
    }
    
    // Parse command content
    return &CustomCommand{
        Content: strings.TrimSpace(matches[1]),
    }
}
```

## Platform-Specific Extensions

### Adding Platform-Specific Features

Use Go's build tags for platform-specific functionality:

```go
// file_windows.go
//go:build windows

func platformSpecificFunction() string {
    return "Windows-specific implementation"
}

// file_unix.go  
//go:build !windows

func platformSpecificFunction() string {
    return "Unix-specific implementation"
}
```

## Testing Extensions

### Unit Testing

Create comprehensive tests for new functionality:

```go
func TestNewFunction(t *testing.T) {
    model := &Model{}
    
    args := NewFunctionArgs{
        Parameter1: "test",
        Parameter2: 42,
    }
    
    result, err := model.executeNewFunction(args)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    expected := "Executed with test and 42"
    if result != expected {
        t.Errorf("expected %s, got %s", expected, result)
    }
}
```

### Integration Testing

Test the complete flow from LLM response to execution:

```go
func TestFunctionCallFlow(t *testing.T) {
    agent, err := llmServer.NewAgent("test-key")
    if err != nil {
        t.Fatal(err)
    }
    
    model := ui.NewModel(agent)
    
    // Simulate function call response
    response := `[FUNCTION_CALL]
Reason: Testing new function
AwaitingCommandApproval: false
Payload:
{
    "name": "NewFunction",
    "args": {"parameter1": "test", "parameter2": 42}
}
[/FUNCTION_CALL]`
    
    // Test parsing and execution
    // ... test implementation
}
```

## Build Process Extensions

### Adding New Build Targets

Extend the build scripts to support new platforms or architectures:

```javascript
// In scripts/build-all.js
const buildTargets = [
    // ... existing targets
    { os: "freebsd", arch: "amd64", output: "menace-go-freebsd" },
    { os: "linux", arch: "arm", output: "menace-go-linux-arm" },
];
```

### Custom Build Steps

Add preprocessing or postprocessing steps:

```javascript
// In build scripts
async function customBuildStep() {
    console.log("Running custom build step...");
    // Custom logic here
}

// Include in build process
await customBuildStep();
```

## Security Considerations

### Input Validation

Always validate input for new functions:

```go
func validateInput(input string) error {
    if len(input) > maxInputLength {
        return fmt.Errorf("input too long: %d characters (max: %d)", len(input), maxInputLength)
    }
    
    // Check for dangerous patterns
    dangerous := []string{"rm -rf", "sudo", "; rm"}
    for _, pattern := range dangerous {
        if strings.Contains(input, pattern) {
            return fmt.Errorf("potentially dangerous input detected")
        }
    }
    
    return nil
}
```

### Approval Mechanisms

Implement appropriate approval workflows for new functions:

```go
func (m *Model) requiresApproval(functionName string, args map[string]interface{}) bool {
    dangerousFunctions := []string{"DeleteFile", "SystemCommand", "NetworkRequest"}
    
    for _, fn := range dangerousFunctions {
        if functionName == fn {
            return true
        }
    }
    
    return false
}
```

## Contribution Guidelines

### Code Style

Follow existing patterns and conventions:
- Use Go formatting tools (`gofmt`, `goimports`)
- Add comprehensive documentation
- Include unit tests for new functionality
- Follow the existing error handling patterns

### Documentation Updates

Always update documentation when adding new features:
- API reference for new functions
- System prompt documentation for new capabilities
- README for new configuration options
- Examples and usage guides

### Pull Request Process

1. Create a feature branch
2. Implement the extension with tests
3. Update documentation
4. Submit PR with detailed description
5. Ensure CI passes and address feedback

This extension system allows Menace CLI to grow and adapt to new requirements while maintaining consistency and reliability.