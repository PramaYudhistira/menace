# Menace CLI User Guide

## Getting Started

Menace CLI is an intelligent terminal assistant that helps you with coding, file management, shell commands, and GitHub operations through natural language conversations.

### Quick Start

1. **Install Menace CLI**:
   ```bash
   npm install -g menace-cli
   ```

2. **Set up your API key**:
   ```bash
   export OPENAI_API_KEY="your_openai_api_key_here"
   ```

3. **Launch the application**:
   ```bash
   menace
   ```

That's it! You're ready to start having conversations with your AI assistant.

## Installation Methods

### Method 1: npm (Recommended)

The easiest way to install Menace CLI:

```bash
npm install -g menace-cli
```

This will install the appropriate binary for your platform automatically.

### Method 2: From Source

For developers or those who want the latest features:

```bash
git clone https://github.com/your-org/menace-cli.git
cd menace-cli
npm run build
npm link
```

## Configuration

### API Keys

Menace CLI supports multiple AI providers:

**OpenAI (Default)**:
```bash
export OPENAI_API_KEY="sk-..."
```

**Anthropic**:
```bash
export ANTHROPIC_API_KEY="sk-ant-..."
```

**For Windows users**:
```powershell
setx OPENAI_API_KEY "sk-..."
```

### Model Selection

Once in the application, you can switch between different models:

- Press `Ctrl+M` to open the model selection menu
- Use arrow keys to navigate
- Press `Enter` to select a model
- Press `Escape` to cancel

Available models include:
- **OpenAI**: GPT-4, GPT-3.5-turbo, o4-mini
- **Anthropic**: Claude-3, Claude-3.5-sonnet
- **Ollama**: Local models (requires Ollama installation)

## User Interface

### Main Interface

```
┌─ Menace CLI ─────────────────────────────────────────┐
│ Chat History Area                                    │
│ > user: Hello! Can you help me list files?          │
│ > assistant: I'll help you list files in the        │
│   current directory.                                 │
│                                                      │
│   [COMMAND_SUGGESTION]                               │
│   Command: ls -la                                    │
│   [/COMMAND_SUGGESTION]                              │
│                                                      │
├──────────────────────────────────────────────────────┤
│ > Type your message here...                          │
└──────────────────────────────────────────────────────┘
```

### Key Features

- **Chat History**: Scrollable conversation history
- **Multi-line Input**: Support for complex inputs
- **Command Preview**: See commands before execution
- **File Operations**: Direct file reading and editing
- **GitHub Integration**: Create PRs and manage repositories

## Basic Usage

### Having Conversations

Simply type your questions or requests naturally:

```
> How do I find all Python files in this directory?
```

The assistant will understand your intent and provide helpful responses, including command suggestions when appropriate.

### Command Execution

When the AI suggests a command, you'll see it formatted like this:

```
[COMMAND_SUGGESTION]
Reason: To find all Python files in the current directory
Command: find . -name "*.py" -type f
AwaitingCommandApproval: false
[/COMMAND_SUGGESTION]
```

- If `AwaitingCommandApproval: false`, the command runs automatically
- If `AwaitingCommandApproval: true`, you must approve it first

### File Operations

Ask for file operations naturally:

```
> Can you read the contents of config.json?
> Please edit the README file and add a new section about installation
> Create a new file called example.py with a hello world program
```

The assistant will use function calls to handle file operations safely:

```
[FUNCTION_CALL]
Reason: Read the config.json file to show its contents
AwaitingCommandApproval: false
Payload:
{
    "name": "ReadFileWithLineNumbers", 
    "args": {
        "path": "config.json"
    }
}
[/FUNCTION_CALL]
```

## Advanced Features

### Multi-line Input

Press `Enter` to create new lines in your input. The assistant can handle complex, multi-paragraph requests.

### Text Selection

- Click and drag to select text
- `Ctrl+A` to select all text in the input area
- `Ctrl+C` to copy selected text
- `Ctrl+X` to cut selected text
- `Ctrl+V` to paste from clipboard

### Navigation

- **Scroll Chat**: Mouse wheel or `Page Up`/`Page Down`
- **Input Navigation**: Arrow keys, `Home`, `End`
- **Clear Input**: `Ctrl+L` or `Esc`
- **Exit**: `Ctrl+C` or `Ctrl+Q`

### GitHub Integration

Menace CLI can help with GitHub operations:

```
> Create a pull request for my current branch
> Stage all my changes and commit them with a good message  
> Push my changes to the remote repository
```

The assistant can:
- Stage files (`git add`)
- Create commits with descriptive messages
- Push to remote repositories
- Create pull requests with appropriate titles and descriptions

## Usage Examples

### Development Workflow

**1. Exploring a new codebase:**
```
> Help me understand this project structure
> What does the main.go file do?
> Find all the TODO comments in the code
```

**2. Code editing:**
```  
> Add error handling to the database connection function
> Refactor this function to be more readable
> Add unit tests for the user authentication module
```

**3. Git operations:**
```
> Show me what files have changed
> Create a commit with all my changes
> Create a pull request for this feature
```

### System Administration

**1. File management:**
```
> Find all log files older than 7 days
> Create a backup of the config directory
> Check disk space usage
```

**2. Process monitoring:**
```
> Show me running processes using port 3000
> Check memory usage
> Find the process consuming the most CPU
```

### Data Analysis

**1. Text processing:**
```
> Count lines in all Python files
> Find all occurrences of 'TODO' in the codebase
> Extract email addresses from this log file
```

**2. File analysis:**
```
> Analyze the structure of this JSON file
> Compare two configuration files
> Generate a summary of this CSV data
```

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Ctrl+C` or `Ctrl+Q` | Exit application |
| `Ctrl+L` or `Esc` | Clear current input |
| `Ctrl+M` | Open model selection |
| `Ctrl+A` | Select all input text |
| `Ctrl+C` | Copy selected text |
| `Ctrl+X` | Cut selected text |
| `Ctrl+V` | Paste from clipboard |
| `Page Up`/`Page Down` | Scroll chat history |
| `Arrow Keys` | Navigate input cursor |
| `Home`/`End` | Move to beginning/end of line |
| `Enter` | Send message (or new line in multi-line input) |

## Troubleshooting

### Common Issues

**1. "OPENAI_API_KEY environment variable not set"**
- Solution: Set your API key as shown in the configuration section
- Verify with: `echo $OPENAI_API_KEY` (Linux/Mac) or `echo %OPENAI_API_KEY%` (Windows)

**2. Commands not executing**
- Check if `AwaitingCommandApproval: true` - you need to approve manually
- Verify you have permissions to execute the suggested command
- Some destructive commands require explicit approval

**3. File operations failing**
- Ensure you have read/write permissions for the target files
- Check if the file path exists and is accessible
- Verify the file isn't locked by another application

**4. GitHub operations not working**
- Ensure you're in a Git repository
- Verify GitHub authentication (gh CLI or SSH keys)
- Check if you have push permissions to the repository

### Terminal Compatibility

Menace CLI is tested on:
- **Linux**: gnome-terminal, konsole, xterm
- **macOS**: Terminal.app, iTerm2
- **Windows**: Windows Terminal, PowerShell, cmd.exe

For best experience:
- Use a terminal with true color support
- Ensure Unicode/UTF-8 support is enabled
- Consider using a modern terminal emulator

### Performance Tips

**For better response times:**
- Use faster models (e.g., GPT-3.5 over GPT-4)
- Keep conversation history reasonable (clear when needed)
- Use local models with Ollama for offline usage

**For better accuracy:**
- Be specific in your requests
- Provide context about your project or goals
- Break complex tasks into smaller steps

## Best Practices

### Effective Communication

**1. Be specific:**
❌ "Fix this code"
✅ "Add error handling to the database connection in main.go"

**2. Provide context:**
❌ "Create a function"
✅ "Create a function that validates email addresses for our user registration system"

**3. Ask for explanations:**
✅ "Explain what this command does before running it"
✅ "Why did you choose this approach?"

### Safety Guidelines

**1. Review commands before execution:**
- Always understand what a command does
- Be cautious with file deletion or system modification commands
- Use `AwaitingCommandApproval: true` for destructive operations

**2. File operation safety:**
- Back up important files before major edits
- Review diffs before applying changes
- Test changes in a development environment first

**3. GitHub operation safety:**
- Review commits before pushing
- Use descriptive commit messages
- Create pull requests rather than direct pushes to main branches

### Productivity Tips

**1. Use conversation history:**
- Previous context helps with follow-up questions
- Reference earlier solutions for similar problems

**2. Combine operations:**
- "Read config.json, then update the database URL and commit the changes"
- Chain related tasks in a single request

**3. Learn from responses:**
- Ask "How does this work?" to understand solutions
- Request explanations of complex commands or code

## Advanced Configuration

### Custom Model Settings

While in the application, you can experiment with different models to find what works best for your use case:

- **Fast responses**: Use GPT-3.5 or o4-mini
- **Complex reasoning**: Use GPT-4 or Claude-3
- **Privacy/offline**: Use Ollama with local models

### Environment Customization

You can customize the experience by setting additional environment variables:

```bash
# Example: Set working directory context
export MENACE_PROJECT_TYPE="web_development"

# Example: Set preferred shell explicitly
export SHELL="/bin/bash"
```

### Integration with Development Tools

Menace CLI works well alongside:
- **Git**: Enhanced commit messages and PR creation
- **VS Code**: Use as an external terminal
- **tmux/screen**: Run in persistent sessions
- **CI/CD**: Get help with pipeline configurations

## Getting Help

### In-Application Help

- Ask the AI directly: "How do I use this application?"
- Request examples: "Show me examples of file operations"
- Get explanations: "Explain what this command does"

### Community and Support

- **GitHub Issues**: Report bugs or request features
- **Documentation**: Check the latest docs for updates
- **Examples**: Browse the examples directory for common use cases

Remember: Menace CLI is designed to be conversational and helpful. Don't hesitate to ask questions, request clarifications, or ask for alternative approaches!