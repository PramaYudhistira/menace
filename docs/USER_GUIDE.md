# Menace CLI User Guide

Welcome to Menace CLI - your intelligent terminal assistant powered by AI. This guide will help you get started and make the most of Menace's capabilities.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Installation](#installation)
3. [Configuration](#configuration)
4. [Basic Usage](#basic-usage)
5. [Advanced Features](#advanced-features)
6. [Model Management](#model-management)
7. [Command Execution](#command-execution)
8. [File Operations](#file-operations)
9. [Git Integration](#git-integration)
10. [Troubleshooting](#troubleshooting)
11. [Tips & Best Practices](#tips--best-practices)

## Quick Start

Get up and running with Menace CLI in 3 simple steps:

1. **Install**: `npm install -g menace-cli`
2. **Configure**: `export OPENAI_API_KEY="your-api-key-here"`
3. **Run**: `menace`

That's it! You now have an AI-powered terminal assistant ready to help.

## Installation

### Option 1: NPM Installation (Recommended)

```bash
npm install -g menace-cli
```

This installs pre-built binaries for your platform (Windows, macOS, Linux).

### Option 2: From Source

If you prefer to build from source or need the latest development version:

```bash
# Clone the repository
git clone https://github.com/your-org/menace-cli.git
cd menace-cli

# Set up the environment
export OPENAI_API_KEY="sk-your-key-here"

# Build all platform binaries
npm run build

# Link for local development
npm link
```

### System Requirements

- **Node.js**: v14.0 or higher
- **Operating System**: Windows 10+, macOS 10.15+, or Linux
- **Memory**: 50MB RAM minimum
- **Disk Space**: 25MB for installation

## Configuration

### API Keys

Menace CLI supports multiple AI providers. Configure the ones you want to use:

#### OpenAI (Required)

```bash
# Unix/Linux/macOS
export OPENAI_API_KEY="sk-your-key-here"

# Windows PowerShell
$env:OPENAI_API_KEY="sk-your-key-here"

# Windows Command Prompt
setx OPENAI_API_KEY "sk-your-key-here"
```

#### Anthropic (Optional)

```bash
export ANTHROPIC_API_KEY="your-anthropic-key-here"
```

#### GitHub Integration (Optional)

For Git operations and pull request creation:

```bash
export GITHUB_TOKEN="ghp_your-token-here"
```

### Persistent Configuration

To make environment variables persistent:

**Linux/macOS**: Add to `~/.bashrc`, `~/.zshrc`, or `~/.profile`:
```bash
echo 'export OPENAI_API_KEY="sk-your-key-here"' >> ~/.bashrc
source ~/.bashrc
```

**Windows**: Use the System Properties dialog or PowerShell:
```powershell
[Environment]::SetEnvironmentVariable("OPENAI_API_KEY", "sk-your-key-here", "User")
```

## Basic Usage

### Starting Menace

Simply run:
```bash
menace
```

You'll see the Menace CLI interface with:
- A configuration panel showing your current model
- A chat area for conversations
- An input area at the bottom

### Your First Conversation

Try these example queries:

```
> How do I list all files in the current directory?
```

```
> Help me create a Node.js project with Express
```

```
> What's the difference between git merge and git rebase?
```

### Interface Navigation

**Keyboard Shortcuts**:
- `Enter`: Send message or execute suggested command
- `Escape`: Cancel current operation or close modals
- `Ctrl+L`: Clear the chat history
- `Ctrl+C`: Copy selected text
- `Ctrl+V`: Paste from clipboard
- `Ctrl+A`: Select all input text
- `Page Up/Down`: Scroll through chat history
- `Tab`: Open model selection (when available)

**Mouse Support**:
- Click to position cursor
- Drag to select text
- Scroll wheel to navigate chat history
- Click buttons to execute commands

## Advanced Features

### Text Selection and Clipboard

Menace CLI supports advanced text editing:

1. **Text Selection**:
   - Click and drag to select text
   - `Shift + Arrow Keys` for keyboard selection
   - `Ctrl+A` to select all input

2. **Clipboard Operations**:
   - `Ctrl+C`: Copy selected text
   - `Ctrl+X`: Cut selected text  
   - `Ctrl+V`: Paste from clipboard

### Multi-line Input

For longer queries or code snippets:

1. Type your first line
2. Continue typing - Menace will automatically handle line breaks
3. Use arrow keys to navigate within your input
4. Press `Enter` to send when ready

Example:
```
> I need help with this Python code:
def fibonacci(n):
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

How can I optimize this?
```

### Context Awareness

Menace maintains conversation context:

```
> Create a new Python file called main.py
AI: I'll create a Python file for you...

> Add a function to calculate factorial
AI: I'll add a factorial function to the main.py file we just created...
```

## Model Management

### Viewing Available Models

Access the model selection interface by pressing `Tab` or clicking the configuration panel.

Available providers and models:

| Provider | Models | Description |
|----------|--------|-------------|
| **OpenAI** | GPT-4, GPT-4 Turbo, GPT-3.5 Turbo | Cloud-based models |
| **Anthropic** | Claude-3 Opus, Sonnet, Haiku | Advanced reasoning models |
| **Ollama** | Llama 2, Code Llama, Mistral | Local models (if installed) |

### Switching Models

1. Press `Tab` to open model selection
2. Use arrow keys or mouse to select
3. Press `Enter` to confirm
4. Model switches immediately for new conversations

### Model Characteristics

**GPT-4**: Best for complex reasoning and analysis
- Excellent for code review and architecture
- Great for detailed explanations
- Higher cost per request

**GPT-3.5 Turbo**: Fast and efficient for most tasks
- Good balance of speed and capability
- Cost-effective for routine queries
- Suitable for command suggestions

**Claude-3 Opus**: Exceptional for creative and analytical tasks
- Superior reasoning capabilities
- Excellent for code generation
- Great for complex problem-solving

**Local Models (Ollama)**: Privacy-focused, no internet required
- No API costs
- Complete privacy
- Requires local installation

### Setting Up Local Models

To use Ollama models:

1. **Install Ollama**:
   ```bash
   curl -fsSL https://ollama.ai/install.sh | sh
   ```

2. **Download Models**:
   ```bash
   ollama pull llama2
   ollama pull codellama
   ```

3. **Start Ollama Service**:
   ```bash
   ollama serve
   ```

4. **Use in Menace**: Select Ollama models from the model menu

## Command Execution

### Understanding Command Suggestions

When Menace suggests shell commands, you'll see:

```
AI: To install Express.js, you can use npm:

[Suggested Command]
npm install express
[Execute] [Cancel]
```

### Safety Levels

Commands are categorized by risk level:

**🟢 Safe Commands** (Auto-approved):
- File reading: `cat`, `less`, `head`, `tail`
- Directory listing: `ls`, `dir`, `find`
- Git status: `git status`, `git log`, `git diff`
- Package information: `npm list`, `pip show`

**🟡 Moderate Risk** (User approval required):
- File creation/editing: `touch`, `mkdir`, `echo >`
- Package installation: `npm install`, `pip install`
- Git operations: `git commit`, `git push`

**🔴 High Risk** (Always requires confirmation):
- File deletion: `rm`, `del`
- System commands: `sudo`, administrative tasks
- Network operations: `curl`, `wget` with modifications

### Approving Commands

When prompted for approval:

1. **Review the command** carefully
2. **Understand what it does** - ask Menace for clarification if needed
3. **Check the directory** you're in
4. **Consider the impact** on your system
5. Press `Enter` to approve or `Escape` to cancel

### Command History

Menace remembers executed commands within the session:

```
> Show me the last command we ran
AI: The last command we executed was: npm install express
```

## File Operations

### Reading Files

```
> Show me the contents of package.json
```

```
> What's in the README.md file?
```

```
> Read the first 10 lines of server.js
```

### Writing and Modifying Files

```
> Create a new file called config.js with basic configuration
```

```
> Add error handling to the main function in app.py
```

```
> Update the version number in package.json to 2.0.0
```

### File Safety Features

- **Automatic backups**: Original files are backed up before modification
- **Diff display**: See exactly what changes will be made
- **Confirmation prompts**: Approve file modifications before they execute
- **Rollback capability**: Undo changes if something goes wrong

### Working with Projects

Menace understands project context:

```
> This is a React project. Add a new component for user authentication
AI: I'll create a React component for user authentication. Based on your project structure, I'll add it to the components directory...
```

```
> Help me set up testing for this Node.js API
AI: I can see this is an Express.js project. I'll help you set up Jest for testing...
```

## Git Integration

### Basic Git Operations

```
> What changes have I made since the last commit?
```

```
> Help me commit these changes with a good commit message
```

```
> Create a new branch for the user authentication feature
```

### Advanced Git Features

#### Pull Request Creation

Menace can automatically create pull requests:

```
> Create a pull request for my authentication feature
AI: I'll create a pull request for your changes:

Title: Add user authentication system
Description: 
- Implemented JWT-based authentication
- Added login/logout endpoints
- Created user registration flow
- Added middleware for protected routes

[Create PR] [Edit] [Cancel]
```

Requirements for PR creation:
- `GITHUB_TOKEN` environment variable set
- Git repository with GitHub remote
- Uncommitted changes or new branch

#### Branch Management

```
> Switch to the main branch and pull latest changes
```

```
> Merge my feature branch and clean up
```

```
> Show me the differences between my branch and main
```

### Git Best Practices

Menace helps you follow Git best practices:

- **Meaningful commit messages**: Suggests clear, descriptive commit messages
- **Branch naming**: Recommends conventional branch naming patterns
- **Conflict resolution**: Guides you through merge conflicts
- **Code review**: Helps prepare code for review

## Troubleshooting

### Common Issues

#### "API Key Not Found" Error

**Problem**: Missing or invalid API key
**Solution**: 
```bash
# Verify key is set
echo $OPENAI_API_KEY

# Set the key (replace with your actual key)
export OPENAI_API_KEY="sk-your-key-here"
```

#### "Model Not Available" Error

**Problem**: Selected model is not accessible
**Solutions**:
1. Check your API key permissions
2. Verify your account has access to the model
3. Try switching to a different model
4. Check your account balance/credits

#### Slow Response Times

**Causes & Solutions**:
- **Network latency**: Check internet connection
- **Model selection**: GPT-4 is slower but more capable than GPT-3.5
- **Query complexity**: Simpler questions get faster responses
- **API limits**: You may be hitting rate limits

#### Terminal Display Issues

**Problem**: Text rendering problems
**Solutions**:
1. Ensure terminal supports Unicode
2. Try a different terminal emulator
3. Adjust terminal window size
4. Check color/theme compatibility

#### Permission Errors

**Problem**: Cannot create files or run commands
**Solutions**:
1. Check directory permissions
2. Run from a directory you own
3. For system operations, consider if `sudo` is needed
4. Verify file paths are correct

### Getting Help

Within Menace, you can ask for help:

```
> How do I use this feature?
```

```
> What keyboard shortcuts are available?
```

```
> I'm getting an error, can you help?
```

### Debug Mode

For persistent issues, you can run with debug information:

```bash
DEBUG=1 menace
```

This provides additional logging information to help diagnose problems.

## Tips & Best Practices

### Effective Querying

#### Be Specific
```
❌ "Help me with Python"
✅ "Help me fix this Python function that's giving a TypeError"
```

#### Provide Context
```
❌ "This doesn't work"  
✅ "This React component isn't rendering. Here's my code: [paste code]"
```

#### Ask Follow-up Questions
```
"Can you explain what this command does before I run it?"
"What are the risks of this approach?"
"Are there alternative solutions?"
```

### Security Best Practices

#### Review Commands Carefully
- Always understand what a command does before running it
- Be especially careful with `rm`, `sudo`, and file modifications
- Ask for explanation if you're unsure

#### Protect Sensitive Information
- Don't paste API keys or passwords in chat
- Be cautious when showing configuration files
- Use environment variables for secrets

#### File Safety
- Keep important files backed up
- Test destructive operations in a safe directory first
- Use version control for important projects

### Workflow Integration

#### Development Workflow
1. Start Menace in your project directory
2. Use it for research and problem-solving
3. Let it suggest and execute safe commands
4. Review and approve file modifications
5. Use Git integration for version control

#### Learning and Documentation
```
> Explain this error message I'm seeing
> What's the best practice for this situation?
> Can you suggest some resources to learn more about this?
```

#### Project Setup
```
> Set up a new React project with TypeScript
> Initialize this directory as a Git repository
> Create a basic project structure for a REST API
```

### Keyboard Shortcuts Summary

| Shortcut | Action |
|----------|--------|
| `Enter` | Send message / Execute command |
| `Escape` | Cancel operation |
| `Ctrl+L` | Clear chat |
| `Ctrl+C/X/V` | Copy/Cut/Paste |
| `Ctrl+A` | Select all |
| `Page Up/Down` | Scroll chat |
| `Tab` | Model selection |
| `Arrow Keys` | Navigate input/selection |

### Model Selection Strategy

**For Code Tasks**:
- **Complex debugging**: GPT-4 or Claude-3 Opus
- **Quick syntax help**: GPT-3.5 Turbo
- **Code review**: GPT-4 or Claude-3 Sonnet
- **Learning concepts**: Any model works well

**For System Administration**:
- **Complex configurations**: GPT-4
- **Quick command help**: GPT-3.5 Turbo
- **Security-related tasks**: GPT-4 or Claude-3 Opus

**For Privacy-Sensitive Work**:
- Use Ollama local models
- Avoid sending sensitive code to cloud models
- Consider using generic examples instead of real data

### Performance Optimization

#### Speed Up Responses
- Use GPT-3.5 Turbo for simpler tasks
- Keep queries concise and focused
- Ask one question at a time
- Use local models for offline work

#### Reduce Costs
- Start with GPT-3.5 Turbo before upgrading to GPT-4
- Use Claude-3 Haiku for cost-effective advanced features
- Consider Ollama for tasks that don't require cloud models
- Clear context periodically with `Ctrl+L`

### Troubleshooting Checklist

When something isn't working:

1. ✅ **Check API keys** are set and valid
2. ✅ **Verify internet connection** for cloud models
3. ✅ **Try a different model** to isolate issues
4. ✅ **Restart Menace** to reset state
5. ✅ **Check terminal compatibility** and size
6. ✅ **Review error messages** carefully
7. ✅ **Ask Menace for help** with specific issues

---

## Getting More Help

- **Documentation**: Check the docs folder for detailed technical documentation
- **Issues**: Report bugs or feature requests on GitHub
- **Community**: Join discussions with other Menace users
- **Updates**: Keep Menace updated with `npm update -g menace-cli`

Welcome to the future of terminal assistance! Menace CLI is designed to make your command-line experience more productive, educational, and enjoyable. The more you use it, the better you'll become at leveraging AI to enhance your development workflow.