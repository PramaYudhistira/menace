# Configuration Guide

This guide covers all configuration options for Menace CLI, including environment variables, command-line options, and runtime settings.

## Table of Contents

- [Environment Variables](#environment-variables)
- [LLM Provider Configuration](#llm-provider-configuration)
- [Terminal Settings](#terminal-settings)
- [Security Settings](#security-settings)
- [Advanced Configuration](#advanced-configuration)
- [Troubleshooting](#troubleshooting)

## Environment Variables

Menace CLI uses environment variables for configuration. Set these in your shell profile or before running the application.

### Required Variables

#### `OPENAI_API_KEY`
OpenAI API key for GPT model access.

**Format:** `sk-...` (starts with sk-)
**Example:**
```bash
export OPENAI_API_KEY="sk-1234567890abcdef1234567890abcdef12345678"
```

**Getting an API key:**
1. Visit [OpenAI API Keys](https://platform.openai.com/api-keys)
2. Create a new secret key
3. Copy the key (it won't be shown again)

### Optional Variables

#### `ANTHROPIC_API_KEY`
Anthropic API key for Claude model access.

**Format:** `sk-ant-...`
**Example:**
```bash
export ANTHROPIC_API_KEY="sk-ant-1234567890abcdef1234567890abcdef12345678"
```

#### `OLLAMA_BASE_URL`
Base URL for Ollama server (local LLM hosting).

**Default:** `http://localhost:11434`
**Example:**
```bash
export OLLAMA_BASE_URL="http://192.168.1.100:11434"
```

#### `MENACE_DEFAULT_MODEL`
Default LLM model to use on startup.

**Options:**
- `gpt-4` (OpenAI)
- `gpt-3.5-turbo` (OpenAI)
- `claude-3-opus` (Anthropic)
- `claude-3-sonnet` (Anthropic)
- `llama2` (Ollama - requires local setup)

**Example:**
```bash
export MENACE_DEFAULT_MODEL="gpt-4"
```

#### `MENACE_MAX_TOKENS`
Maximum tokens per API request.

**Default:** `4096`
**Range:** `100-8192`
**Example:**
```bash
export MENACE_MAX_TOKENS="2048"
```

#### `MENACE_TEMPERATURE`
AI response creativity (0.0 = deterministic, 1.0 = creative).

**Default:** `0.7`
**Range:** `0.0-1.0`
**Example:**
```bash
export MENACE_TEMPERATURE="0.5"
```

## LLM Provider Configuration

### OpenAI Configuration

OpenAI is the primary supported provider with the most features.

**Supported Models:**
- `gpt-4` - Most capable, higher cost
- `gpt-4-turbo` - Fast and capable
- `gpt-3.5-turbo` - Faster, lower cost

**Rate Limits:**
- Free tier: 3 requests/minute
- Paid tier: 500+ requests/minute (varies by model)

**Configuration:**
```bash
export OPENAI_API_KEY="your-key-here"
export MENACE_DEFAULT_MODEL="gpt-4"
export MENACE_MAX_TOKENS="4096"
```

### Anthropic Configuration

Anthropic Claude models offer strong performance with different capabilities.

**Supported Models:**
- `claude-3-opus` - Most capable
- `claude-3-sonnet` - Balanced performance
- `claude-3-haiku` - Fastest

**Configuration:**
```bash
export ANTHROPIC_API_KEY="your-key-here"
export MENACE_DEFAULT_MODEL="claude-3-sonnet"
```

### Ollama Configuration (Local LLMs)

Run models locally with Ollama for privacy and offline use.

**Setup Ollama:**
```bash
# Install Ollama (macOS/Linux)
curl -fsSL https://ollama.ai/install.sh | sh

# Pull a model
ollama pull llama2
ollama pull codellama
```

**Configuration:**
```bash
export OLLAMA_BASE_URL="http://localhost:11434"
export MENACE_DEFAULT_MODEL="llama2"
```

**Available Models:**
- `llama2` - General purpose
- `codellama` - Code-focused
- `mistral` - Fast and efficient
- `dolphin-mixtral` - Advanced reasoning

## Terminal Settings

### Terminal Compatibility

Menace CLI works best with modern terminal emulators that support:
- 256 colors or true color
- Mouse input
- Unicode/UTF-8
- Terminal resize events

**Recommended Terminals:**
- **macOS:** Terminal.app, iTerm2
- **Windows:** Windows Terminal, ConEmu
- **Linux:** gnome-terminal, konsole, alacritty

### Display Settings

#### `MENACE_UI_THEME`
UI color theme selection.

**Options:** `dark`, `light`, `auto`
**Default:** `auto`
**Example:**
```bash
export MENACE_UI_THEME="dark"
```

#### `MENACE_FONT_SIZE`
Relative font size adjustment (terminal dependent).

**Options:** `small`, `medium`, `large`
**Default:** `medium`
**Example:**
```bash
export MENACE_FONT_SIZE="large"
```

#### `MENACE_SIDEBAR_WIDTH`
Sidebar width in characters.

**Default:** `30`
**Range:** `20-50`
**Example:**
```bash
export MENACE_SIDEBAR_WIDTH="35"
```

### Mouse Support

Enable/disable mouse support in the terminal.

#### `MENACE_MOUSE_SUPPORT`
**Options:** `true`, `false`
**Default:** `true`
**Example:**
```bash
export MENACE_MOUSE_SUPPORT="false"
```

## Security Settings

### Command Approval

#### `MENACE_AUTO_APPROVE_SAFE`
Automatically approve commands marked as safe.

**Options:** `true`, `false`
**Default:** `false`
**Example:**
```bash
export MENACE_AUTO_APPROVE_SAFE="true"
```

#### `MENACE_REQUIRE_APPROVAL_ALL`
Require approval for all commands (maximum security).

**Options:** `true`, `false`
**Default:** `false`
**Example:**
```bash
export MENACE_REQUIRE_APPROVAL_ALL="true"
```

### File Access

#### `MENACE_ALLOWED_PATHS`
Restrict file access to specific paths (colon-separated).

**Example:**
```bash
export MENACE_ALLOWED_PATHS="/home/user/projects:/tmp"
```

#### `MENACE_DENIED_PATHS`
Explicitly deny access to paths (colon-separated).

**Example:**
```bash
export MENACE_DENIED_PATHS="/etc:/root:/home/user/.ssh"
```

### Network Access

#### `MENACE_DISABLE_GITHUB`
Disable GitHub integration features.

**Options:** `true`, `false`
**Default:** `false`
**Example:**
```bash
export MENACE_DISABLE_GITHUB="true"
```

## Advanced Configuration

### Logging

#### `MENACE_LOG_LEVEL`
Set logging verbosity.

**Options:** `debug`, `info`, `warn`, `error`
**Default:** `info`
**Example:**
```bash
export MENACE_LOG_LEVEL="debug"
```

#### `MENACE_LOG_FILE`
Write logs to file instead of stderr.

**Example:**
```bash
export MENACE_LOG_FILE="/tmp/menace.log"
```

### Performance

#### `MENACE_MESSAGE_HISTORY`
Maximum number of messages to keep in memory.

**Default:** `100`
**Range:** `10-1000`
**Example:**
```bash
export MENACE_MESSAGE_HISTORY="200"
```

#### `MENACE_API_TIMEOUT`
API request timeout in seconds.

**Default:** `30`
**Range:** `5-120`
**Example:**
```bash
export MENACE_API_TIMEOUT="60"
```

### Development

#### `MENACE_DEBUG_UI`
Enable UI debugging information.

**Options:** `true`, `false`
**Default:** `false`
**Example:**
```bash
export MENACE_DEBUG_UI="true"
```

#### `MENACE_CUSTOM_PROMPT`
Path to custom system prompt file.

**Example:**
```bash
export MENACE_CUSTOM_PROMPT="/path/to/custom_prompt.txt"
```

## Platform-Specific Configuration

### Windows

Set environment variables in PowerShell:
```powershell
# Permanent setting
setx OPENAI_API_KEY "your-key-here"

# Session-only setting
$env:OPENAI_API_KEY = "your-key-here"
```

In Command Prompt:
```cmd
set OPENAI_API_KEY=your-key-here
```

### macOS/Linux

Add to shell profile (`~/.bashrc`, `~/.zshrc`, etc.):
```bash
export OPENAI_API_KEY="your-key-here"
export MENACE_DEFAULT_MODEL="gpt-4"
```

For session-only:
```bash
OPENAI_API_KEY="your-key-here" menace
```

## Configuration Files

### Future Support

Configuration file support is planned for future versions:

```yaml
# ~/.config/menace/config.yaml
llm:
  default_provider: "openai"
  default_model: "gpt-4"
  max_tokens: 4096
  temperature: 0.7

ui:
  theme: "dark"
  sidebar_width: 30
  mouse_support: true

security:
  auto_approve_safe: false
  require_approval_all: false
  allowed_paths: ["/home/user/projects", "/tmp"]
```

## Troubleshooting

### Common Issues

#### API Key Not Working
```bash
# Check if set
echo $OPENAI_API_KEY

# Test manually
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
  https://api.openai.com/v1/models
```

#### Model Not Found
- Verify the model name is correct
- Check your API plan supports the model
- Ensure sufficient credits/quota

#### Terminal Display Issues
- Update your terminal emulator
- Check terminal size: `echo $COLUMNS x $LINES`
- Try different terminal settings

#### Connection Issues
```bash
# Test internet connectivity
curl -I https://api.openai.com/

# Check proxy settings
echo $HTTP_PROXY
echo $HTTPS_PROXY
```

### Debug Mode

Enable maximum debugging:
```bash
export MENACE_LOG_LEVEL="debug"
export MENACE_DEBUG_UI="true"
menace
```

### Getting Help

1. Check the logs for error messages
2. Verify all required environment variables are set
3. Test with minimal configuration
4. Check the [troubleshooting section](../README.md#troubleshooting) in README
5. Open an issue on GitHub with:
   - Your operating system
   - Terminal emulator
   - Environment variables (redacted)
   - Error messages

## Configuration Examples

### Basic Setup
```bash
export OPENAI_API_KEY="sk-1234..."
menace
```

### Multi-Provider Setup
```bash
export OPENAI_API_KEY="sk-1234..."
export ANTHROPIC_API_KEY="sk-ant-5678..."
export MENACE_DEFAULT_MODEL="claude-3-sonnet"
menace
```

### High-Security Setup
```bash
export OPENAI_API_KEY="sk-1234..."
export MENACE_REQUIRE_APPROVAL_ALL="true"
export MENACE_ALLOWED_PATHS="/home/user/safe-projects"
export MENACE_DISABLE_GITHUB="true"
menace
```

### Development Setup
```bash
export OPENAI_API_KEY="sk-1234..."
export MENACE_LOG_LEVEL="debug"
export MENACE_DEBUG_UI="true"
export MENACE_MESSAGE_HISTORY="500"
menace
```

This configuration guide should help you customize Menace CLI to your specific needs and environment.