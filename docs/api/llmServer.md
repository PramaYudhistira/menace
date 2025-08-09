# LLMServer Package API Documentation

The `llmServer` package provides the core AI functionality for Menace CLI, including multi-provider LLM support, command parsing, and GitHub integration.

## Table of Contents

- [Agent](#agent)
- [LLM Service](#llm-service)
- [Command Suggestions](#command-suggestions)
- [System Prompts](#system-prompts)
- [Model Factory](#model-factory)
- [GitHub Interface](#github-interface)

## Agent

### Overview
The `Agent` struct is the primary interface for interacting with Large Language Models. It supports multiple providers including OpenAI, Anthropic, and Ollama through the LangChain Go library.

### Structure
```go
type Agent struct {
    chain    *chains.ConversationChain
    memory   *memory.ConversationBufferMemory
    provider string
}
```

### Functions

#### `NewAgent(apiKey string) (*Agent, error)`
Creates a new AI agent with the specified API key.

**Parameters:**
- `apiKey` (string): The API key for the LLM provider

**Returns:**
- `*Agent`: A configured agent instance
- `error`: Error if agent creation fails

**Example:**
```go
agent, err := llmServer.NewAgent("sk-your-api-key-here")
if err != nil {
    log.Fatal(err)
}
```

#### `SendMessage(message string) (string, []CommandSuggestion, error)`
Processes a user message and returns the AI response along with any extracted command suggestions.

**Parameters:**
- `message` (string): The user's input message

**Returns:**
- `string`: The AI's response text
- `[]CommandSuggestion`: Parsed command suggestions from the response
- `error`: Error if message processing fails

## LLM Service

### Overview
The `LLMService` provides a singleton pattern for managing OpenAI GPT model interactions. This is legacy code that's been superseded by the Agent implementation.

### Structure
```go
type LLMService struct {
    client     *openai.Client
    model      string
    maxTokens  int
    context    []openai.ChatCompletionMessage
}
```

### Functions

#### `GetLLMService() *LLMService`
Returns the singleton instance of LLMService.

#### `SendMessage(message string) (string, error)`
Sends a message to the OpenAI API and returns the response.

## Command Suggestions

### Overview
The command suggestion system parses structured commands from natural language LLM responses using special markup tags.

### Structure
```go
type CommandSuggestion struct {
    Reason             string
    Command            string
    RequiresApproval   bool
    IsFileOperation    bool
}
```

### Functions

#### `ParseCommandSuggestions(response string) []CommandSuggestion`
Extracts command suggestions from an LLM response using markup patterns.

**Parameters:**
- `response` (string): The raw LLM response text

**Returns:**
- `[]CommandSuggestion`: List of parsed command suggestions

**Markup Format:**
Commands are extracted using this pattern:
```
<command reason="explanation" approval="true|false" type="shell|file">
actual command here
</command>
```

## System Prompts

### Overview
Contains the system prompts that define the AI's behavior and capabilities.

### Functions

#### `GetSystemPrompt() string`
Returns the current system prompt that defines the AI assistant's role and capabilities.

#### `GetNewSystemPrompt() string`
Returns an updated JSON-based system prompt format with enhanced capabilities.

**Capabilities defined in prompts:**
- Shell command assistance
- File reading and editing operations
- Git and GitHub integration
- Cross-platform compatibility
- Code analysis and suggestions

## Model Factory

### Overview
The `ModelFactory` detects the user's operating system and shell environment to provide appropriate command suggestions.

### Functions

#### `DetectShellEnvironment() string`
Detects the current shell environment and returns a description.

**Returns:**
- `string`: Description of the detected environment (e.g., "Windows PowerShell", "Linux bash", "macOS zsh")

**Detection Logic:**
- Windows: Distinguishes between CMD and PowerShell
- Unix-like: Detects bash, zsh, fish, and other common shells
- Fallback: Provides generic Unix shell detection

## GitHub Interface

### Overview
Provides GitHub integration capabilities for repository operations.

### Functions

#### `CreatePullRequest(title, body, baseBranch, headBranch string) error`
Creates a new pull request on GitHub.

**Parameters:**
- `title` (string): PR title
- `body` (string): PR description
- `baseBranch` (string): Target branch for the PR
- `headBranch` (string): Source branch for the PR

#### `PushToGitHub(branch string) error`
Pushes the current branch to GitHub.

**Parameters:**
- `branch` (string): Branch name to push

#### `CheckForUnpushedChanges() (bool, error)`
Checks if there are local changes that haven't been pushed to the remote repository.

**Returns:**
- `bool`: True if there are unpushed changes
- `error`: Error if check fails

#### `IsolatedAIMessage(message string) (string, error)`
Sends a message to the AI without affecting the main conversation context.

**Parameters:**
- `message` (string): Message to send

**Returns:**
- `string`: AI response
- `error`: Error if message fails

## Error Handling

All functions in the llmServer package follow Go's idiomatic error handling patterns:
- Functions return errors as their last return value
- Errors are wrapped with contextual information
- Critical errors may cause the application to exit

## Configuration

The package uses environment variables for configuration:
- `OPENAI_API_KEY`: Required for OpenAI provider
- `ANTHROPIC_API_KEY`: Required for Anthropic provider
- `OLLAMA_BASE_URL`: Base URL for Ollama server (optional)

## Thread Safety

The Agent implementation is designed to be thread-safe through the use of LangChain Go's memory management. However, concurrent access to the same agent instance should be handled carefully.