# UI Package API Documentation

The `ui` package provides the terminal user interface for Menace CLI using the Bubble Tea framework. It includes a sophisticated text editor, chat interface, and file operation capabilities.

## Table of Contents

- [Model](#model)
- [View](#view)
- [Message](#message)
- [Executor](#executor)
- [Key Bindings](#key-bindings)
- [UI Components](#ui-components)

## Model

### Overview
The `Model` struct is the core UI state manager that handles all user interactions, text editing, and communication with the AI agent.

### Structure
```go
type Model struct {
    // Text editing state
    input           []string
    cursor          int
    selectedText    string
    selection       Selection
    clipboard       string
    
    // UI state
    width           int
    height          int
    viewportHeight  int
    scrollOffset    int
    
    // Chat and AI
    messages        []Message
    agent           *llmServer.Agent
    awaitingAI      bool
    pendingCommands []llmServer.CommandSuggestion
    
    // Configuration
    currentModel    string
    availableModels []string
}
```

### Key Methods

#### `NewModel(agent *llmServer.Agent) Model`
Creates a new UI model with the specified AI agent.

**Parameters:**
- `agent` (*llmServer.Agent): The AI agent to use for interactions

**Returns:**
- `Model`: Initialized UI model

#### `Init() tea.Cmd`
Initializes the Bubble Tea model and returns initial commands.

#### `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
Handles all UI events and updates the model state.

**Parameters:**
- `msg` (tea.Msg): The message/event to process

**Returns:**
- `tea.Model`: Updated model state
- `tea.Cmd`: Command to execute

#### `View() string`
Renders the complete UI and returns the display string.

**Returns:**
- `string`: Rendered UI string

### Text Editing Features

#### Cursor Movement
- Arrow keys for character-by-character movement
- Ctrl+Left/Right for word-by-word movement
- Home/End for line beginning/end
- Page Up/Down for viewport scrolling

#### Text Selection
```go
type Selection struct {
    startLine int
    startChar int
    endLine   int
    endChar   int
}
```

- Shift + movement keys for selection
- Ctrl+A for select all
- Mouse selection support

#### Copy/Paste Operations
- Ctrl+C for copy
- Ctrl+V for paste
- Ctrl+X for cut
- Cross-platform clipboard integration

#### Multi-line Editing
- Enter for new lines
- Backspace/Delete for character deletion
- Tab for indentation
- Auto-indentation support

## View

### Overview
The view system renders a two-panel layout with a sidebar and main chat area.

### Layout Components

#### Sidebar
- System information display
- Model selection interface
- Configuration options
- Status indicators

#### Main Chat Area
- Message history display
- Input text editor
- Command suggestion prompts
- Scroll indicators

### Styling

Uses Lipgloss for sophisticated terminal styling:

```go
var (
    sidebarStyle = lipgloss.NewStyle().
        Border(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("240"))
    
    messageStyle = lipgloss.NewStyle().
        Padding(0, 1).
        Margin(1, 0)
    
    inputStyle = lipgloss.NewStyle().
        Border(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("62"))
)
```

### Responsive Design
- Automatic text wrapping
- Dynamic height adjustment
- Terminal resize handling
- Minimum/maximum size constraints

## Message

### Structure
```go
type Message struct {
    Sender  string // "user", "llm", "system"
    Content string
}
```

### Message Types
- **User messages**: Input from the user
- **LLM messages**: AI responses
- **System messages**: Status and information messages

### Message Rendering
- Syntax highlighting for code blocks
- Automatic link detection
- Markdown-like formatting
- Timestamp display (optional)

## Executor

### Overview
The executor handles file operations and shell command execution.

### File Operations

#### `ReadFileWithLineNumbers(filePath string) (string, error)`
Reads a file and returns its contents with line numbers.

**Parameters:**
- `filePath` (string): Path to the file to read

**Returns:**
- `string`: File contents with line numbers
- `error`: Error if file reading fails

**Example Output:**
```
1: package main
2: 
3: import "fmt"
4: 
5: func main() {
```

#### `CreateAndApplyDiffs(filePath string, diffs []DiffOperation) error`
Applies diff operations to a file.

```go
type DiffOperation struct {
    Type     string // "add", "delete", "modify"
    LineNum  int
    Content  string
    Count    int    // For delete operations
}
```

**Parameters:**
- `filePath` (string): Target file path
- `diffs` ([]DiffOperation): List of operations to apply

### Shell Command Execution

#### `ExecuteCommand(command string) (string, error)`
Executes a shell command and returns its output.

**Parameters:**
- `command` (string): Command to execute

**Returns:**
- `string`: Command output
- `error`: Error if execution fails

**Security Features:**
- Command validation
- Approval requirements for dangerous operations
- Sandboxed execution environment

## Key Bindings

### Navigation
- `↑/↓`: Move cursor up/down
- `←/→`: Move cursor left/right
- `Ctrl+←/→`: Move by word
- `Home/End`: Move to line start/end
- `Page Up/Down`: Scroll viewport

### Text Selection
- `Shift + movement`: Select text
- `Ctrl+A`: Select all
- `Mouse drag`: Select with mouse

### Editing
- `Enter`: New line
- `Backspace`: Delete previous character
- `Delete`: Delete next character
- `Tab`: Insert indentation
- `Ctrl+Z`: Undo (limited support)

### Copy/Paste
- `Ctrl+C`: Copy selection
- `Ctrl+V`: Paste
- `Ctrl+X`: Cut selection

### Application Control
- `Ctrl+Q`: Quit application
- `Ctrl+R`: Refresh/redraw
- `Esc`: Cancel current operation

### AI Interaction
- `Ctrl+Enter`: Send message to AI
- `Ctrl+S`: Save current input
- `Tab`: Accept command suggestion

## UI Components

### Input Area
- Multi-line text editor
- Syntax highlighting
- Auto-completion
- Line number display (optional)

### Message Display
- Scrollable message history
- Message type indicators
- Timestamp display
- Code block formatting

### Status Bar
- Current model indicator
- Connection status
- Input mode indicator
- Scroll position

### Command Suggestions
- Interactive approval prompts
- Command preview
- Risk level indicators
- Execution options

## Event Handling

### Mouse Events
- Click to position cursor
- Drag to select text
- Scroll for viewport navigation
- Double-click for word selection

### Keyboard Events
- Character input handling
- Special key combinations
- IME support for international input
- Modifier key detection

### Window Events
- Terminal resize handling
- Focus/blur events
- Suspend/resume support

## Performance Considerations

### Text Rendering
- Efficient diff-based updates
- Virtual scrolling for large files
- Lazy loading of message history
- Optimized syntax highlighting

### Memory Management
- Circular buffer for message history
- Garbage collection of old selections
- Efficient string operations
- Buffer pooling for large operations

## Accessibility

### Features
- Screen reader compatibility
- High contrast mode support
- Keyboard-only navigation
- Configurable font sizes

### Compatibility
- Terminal emulator support
- Unicode and emoji handling
- Cross-platform key binding normalization
- Fallback rendering for limited terminals