# Menace CLI - User Guide

## Table of Contents

1. [Getting Started](#getting-started)
2. [Basic Usage](#basic-usage)
3. [Advanced Features](#advanced-features)
4. [Configuration](#configuration)
5. [Keyboard Shortcuts](#keyboard-shortcuts)
6. [Tips and Best Practices](#tips-and-best-practices)
7. [Troubleshooting](#troubleshooting)

## Getting Started

### Prerequisites

Before installing Menace CLI, ensure you have:

- **Go 1.18+** (for building from source)
- **Node.js v14+** and npm
- An **OpenAI API key** (get one from [OpenAI Platform](https://platform.openai.com/))
- Optional: **Anthropic API key** for Claude models
- Optional: **Ollama** installed for local models

### Installation

#### Option 1: From npm (Recommended)

```bash
npm install -g menace-cli
```

#### Option 2: From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/PramaYudhistira/terminal-chat.git
   cd terminal-chat/menace-cli
   ```

2. Set your API keys:
   ```bash
   export OPENAI_API_KEY="your_key_here"
   export ANTHROPIC_API_KEY="your_anthropic_key_here"  # Optional
   ```

3. Build all platform binaries:
   ```bash
   npm run build
   ```

4. Link the CLI globally:
   ```bash
   npm link
   ```

### First Launch

After installation, launch Menace CLI:

```bash
menace
```

You'll see a terminal interface with:
- **Chat area**: Displays conversation history
- **Input area**: Where you type messages
- **Status line**: Shows current model and connection status

## Basic Usage

### Starting a Conversation

Simply type your message and press **Enter** to send it to the AI assistant:

```
> Hello, can you help me list all Python files in this directory?
```

The AI will respond and may suggest commands:

```
AI: I'll help you list all Python files. Let me run the find command:

[Command Suggestion] find . -name "*.py" -type f
```

### Command Execution

When the AI suggests a command, you have several options:

- **Press Enter**: Execute the command automatically
- **Type "no"**: Reject the command
- **Type feedback**: Modify the approach
- **Press Escape**: Cancel and continue chatting

### Multi-line Input

Press **Shift+Enter** to add new lines without sending:

```
> I need help with a Python script that:
  - Reads a CSV file
  - Processes the data
  - Saves results to JSON
```

## Advanced Features

### File Operations

Menace CLI can read, edit, and create files through AI assistance:

#### Reading Files

```
> Show me the contents of config.py
```

The AI will read the file and display it with line numbers:

```
1: import os
2: import json
3: 
4: CONFIG_PATH = "settings.json"
```

#### Editing Files

```
> Add error handling to line 5 in config.py
```

The AI will propose specific changes and apply them after your confirmation.

#### Creating Files

```
> Create a new Python script for data validation with proper error handling
```

### Model Selection

Press **Ctrl+M** to open the model selection menu:

- **OpenAI Models**: GPT-4, GPT-3.5, and latest models
- **Anthropic Models**: Claude 3, Claude Instant
- **Local Models**: Ollama-supported models

Use arrow keys to navigate and Enter to select.

### GitHub Integration

Menace CLI integrates with Git and GitHub:

#### Commit Changes

```
> Stage all changes and commit them with a descriptive message
```

#### Create Pull Requests

```
> Create a pull request for this feature branch
```

The AI will:
1. Analyze your changes
2. Generate an appropriate title and description
3. Create the PR using GitHub CLI

### Text Selection and Clipboard

#### Selecting Text

- **Mouse**: Click and drag to select text
- **Keyboard**: Hold Shift and use arrow keys
- **Select All**: Ctrl+A (Cmd+A on macOS)

#### Clipboard Operations

- **Copy**: Ctrl+C (Cmd+C on macOS)
- **Cut**: Ctrl+X (Cmd+X on macOS)
- **Paste**: Ctrl+V (Cmd+V on macOS)

### Scrolling and Navigation

- **Mouse Wheel**: Scroll through chat history
- **Page Up/Down**: Navigate quickly through long conversations
- **Home/End**: Jump to beginning/end of current line
- **Ctrl+Home/End**: Jump to beginning/end of input area

## Configuration

### Environment Variables

Set these environment variables for full functionality:

```bash
# Required for OpenAI models
export OPENAI_API_KEY="sk-your-key-here"

# Optional for Anthropic models
export ANTHROPIC_API_KEY="your-anthropic-key"

# Optional for custom endpoints
export OPENAI_BASE_URL="https://your-custom-endpoint.com"
```

### Model Configuration

The AI assistant can work with different models optimized for various tasks:

#### Code-Focused Models
- **GPT-4**: Best for complex coding tasks
- **Claude 3**: Excellent for code review and analysis

#### General Purpose Models
- **GPT-3.5 Turbo**: Fast and cost-effective
- **Claude Instant**: Quick responses for simple tasks

#### Local Models (Ollama)
- **CodeLlama**: Local code assistance
- **Mistral**: General purpose local model

### Custom System Prompts

The system prompt automatically adapts to:
- Your operating system and shell
- Current working directory
- Available tools and integrations

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| **Enter** | Send message / Execute suggested command |
| **Shift+Enter** | New line in input |
| **Escape** | Cancel pending operations |
| **Ctrl+C** (Cmd+C) | Copy selected text |
| **Ctrl+V** (Cmd+V) | Paste from clipboard |
| **Ctrl+X** (Cmd+X) | Cut selected text |
| **Ctrl+A** (Cmd+A) | Select all text in input |
| **Ctrl+M** | Open model selection |
| **Ctrl+L** | Clear conversation history |
| **Ctrl+D** | Exit application |
| **Page Up/Down** | Scroll chat history |
| **Arrow Keys** | Navigate cursor |
| **Shift+Arrows** | Select text |

## Tips and Best Practices

### Effective Communication

1. **Be Specific**: Instead of "fix my code," say "fix the syntax error in line 15 of main.py"

2. **Provide Context**: Share error messages, file contents, or project details

3. **Break Down Complex Tasks**: 
   ```
   > I need to build a REST API. Let's start with:
     1. Setting up the project structure
     2. Installing dependencies
     3. Creating the main server file
   ```

### Code Assistance

1. **Show Files First**: Let the AI see your code before asking for changes
   ```
   > Show me users.py, then help me add input validation
   ```

2. **Explain Requirements**: Describe what you want to achieve
   ```
   > Add logging to this function that captures both success and error cases
   ```

3. **Review Suggestions**: Always review generated code before accepting

### File Management

1. **Use Relative Paths**: Reference files relative to your current directory
2. **Backup Important Files**: The AI can modify files, so keep backups
3. **Version Control**: Use git to track changes made by the AI

### Performance Optimization

1. **Clear History**: Use Ctrl+L to clear long conversations and improve performance
2. **Use Appropriate Models**: Choose faster models for simple tasks
3. **Batch Related Questions**: Group related questions to maintain context

## Troubleshooting

### Common Issues

#### API Key Problems

**Issue**: "Error: OPENAI_API_KEY environment variable not set"

**Solution**: 
```bash
export OPENAI_API_KEY="your-actual-api-key"
```

Make it permanent by adding to your shell configuration:
```bash
echo 'export OPENAI_API_KEY="your-key"' >> ~/.bashrc  # or ~/.zshrc
source ~/.bashrc  # or ~/.zshrc
```

#### Connection Issues

**Issue**: "Failed to get response from LLM"

**Solutions**:
1. Check your internet connection
2. Verify API key is valid and has credits
3. Try switching to a different model
4. Check for API rate limits

#### File Permission Errors

**Issue**: "Permission denied when reading/writing files"

**Solutions**:
1. Check file permissions: `ls -la filename`
2. Use `sudo` if necessary (with caution)
3. Ensure you're in the correct directory

#### Binary Not Found

**Issue**: Platform-specific binary missing

**Solution**:
```bash
cd menace-cli
npm run build  # Rebuild all binaries
```

### Performance Issues

#### Slow Responses

1. **Check Model**: Larger models are slower but more capable
2. **Clear History**: Long conversations can slow down responses
3. **Reduce Context**: Avoid very long file contents or outputs

#### High Memory Usage

1. **Restart Application**: Ctrl+D then restart
2. **Clear Chat History**: Ctrl+L
3. **Use Smaller Models**: Switch to GPT-3.5 or Claude Instant

### Getting Help

If you encounter issues:

1. **Check Logs**: Look for error messages in the terminal
2. **GitHub Issues**: Report bugs at [GitHub repository](https://github.com/PramaYudhistira/terminal-chat/issues)
3. **Clear State**: Try restarting with a fresh conversation
4. **Update**: Ensure you're using the latest version

### Debug Mode

For troubleshooting, you can enable verbose logging:

```bash
DEBUG=1 menace
```

This will show additional information about API calls, file operations, and internal state changes.

---

## Next Steps

- Explore the [API Documentation](API.md) for technical details
- Check the [Contributing Guide](CONTRIBUTING.md) to help improve Menace CLI
- Read the [Deployment Guide](DEPLOYMENT.md) for advanced configurations