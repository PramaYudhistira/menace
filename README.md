# Menace CLI

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.18%2B-blue.svg)
![Node Version](https://img.shields.io/badge/node-14%2B-green.svg)

Menace CLI is a powerful, cross-platform, terminal-based AI assistant that brings intelligent automation to your command line. Built with Go and powered by advanced language models, Menace provides an interactive text user interface (TUI) that can understand natural language, execute commands, manipulate files, and assist with development workflows.

## ✨ Key Features

- **🤖 Multi-Provider LLM Support**: OpenAI GPT, Anthropic Claude, and local Ollama models
- **💻 Cross-Platform**: Native binaries for Windows, macOS (Intel & Apple Silicon), and Linux
- **🖱️ Interactive TUI**: Mouse support, syntax highlighting, and responsive design
- **📝 Advanced Text Editor**: Full-featured editing with selections, copy/paste, and multi-line input
- **🔧 Smart Command Suggestions**: Context-aware shell command recommendations
- **📁 File Operations**: Read, edit, and create files with diff-based modifications
- **🐙 GitHub Integration**: Built-in git operations and pull request creation
- **🛡️ Security First**: Command approval system and sandboxed execution
- **⚡ Fast & Lightweight**: Efficient Go backend with minimal resource usage

## 🚀 Quick Start

### Prerequisites

- **Go 1.18+** (for building from source)
- **Node.js v14+** and npm (for build scripts)
- **OpenAI API key** or other supported LLM provider

### Installation

#### Option 1: From npm (Coming Soon!)

```bash
npm install --global menace-cli
```

#### Option 2: From Source

```bash
# Clone the repository
git clone https://github.com/your-org/menace-cli.git
cd menace-cli/menace-cli

# Set up your API key
export OPENAI_API_KEY="your-openai-key-here"

# Build for all platforms
npm run build

# Link for local use
npm link
```

### First Run

```bash
# Start Menace CLI
menace

# Or set the API key inline
OPENAI_API_KEY="your-key" menace
```

## 🎮 Usage Examples

### Basic Interactions

```bash
# Ask for help with commands
> How do I find files larger than 100MB?

# Get assistance with git operations
> Create a new branch and switch to it

# File manipulation
> Read the package.json file and show me the dependencies

# Code assistance
> Explain what this function does: [select code]
```

### Advanced Features

- **Multi-line Input**: Shift+Enter for new lines
- **Text Selection**: Mouse drag or Shift+arrows
- **Command History**: Up/down arrows
- **File Operations**: Built-in diff-based editing
- **GitHub Integration**: Create PRs, push changes

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `OPENAI_API_KEY` | OpenAI API key (required) | - | `sk-1234...` |
| `ANTHROPIC_API_KEY` | Anthropic API key | - | `sk-ant-1234...` |
| `OLLAMA_BASE_URL` | Ollama server URL | `http://localhost:11434` | `http://192.168.1.100:11434` |
| `MENACE_DEFAULT_MODEL` | Default LLM model | `gpt-3.5-turbo` | `gpt-4` |
| `MENACE_MAX_TOKENS` | Max tokens per request | `4096` | `2048` |

For complete configuration options, see [CONFIGURATION.md](docs/CONFIGURATION.md).

### Setting API Keys

**Unix/Linux/macOS:**
```bash
# Add to ~/.bashrc, ~/.zshrc, etc.
export OPENAI_API_KEY="your-key-here"
```

**Windows PowerShell:**
```powershell
# Permanent setting
setx OPENAI_API_KEY "your-key-here"

# Session only
$env:OPENAI_API_KEY = "your-key-here"
```

## 🏗️ Development

### Directory Structure

```
menace-cli/
├── bin/                   # Compiled binaries
├── docs/                  # Documentation
│   ├── api/              # API documentation
│   ├── CONFIGURATION.md   # Configuration guide
│   └── DEVELOPER_GUIDE.md # Development guide
├── scripts/               # Build scripts
└── src/                   # Go source code
    ├── main.go           # Application entry point
    ├── llmServer/        # AI and LLM functionality
    └── ui/               # Terminal user interface
```

### Build Commands

```bash
# Build all platforms
npm run build

# Build specific platforms
npm run build-windows
npm run build-mac

# Development build
cd src && go build .

# Run tests
go test ./...
```

### Architecture

Menace CLI follows a modular architecture:

- **LLMServer Package**: AI integration, command parsing, GitHub operations
- **UI Package**: Terminal interface, text editing, event handling
- **Main**: Application bootstrap and coordination

For detailed architecture information, see [DEVELOPER_GUIDE.md](docs/DEVELOPER_GUIDE.md).

## 📚 Documentation

- [API Documentation](docs/api/) - Detailed API reference
- [Developer Guide](docs/DEVELOPER_GUIDE.md) - Architecture and development setup
- [Configuration Guide](docs/CONFIGURATION.md) - Complete configuration reference
- [Contributing Guidelines](docs/CONTRIBUTING.md) - How to contribute

## 🔒 Security

Menace CLI includes several security features:

- **Command Approval**: Review commands before execution
- **File Access Control**: Restrict access to sensitive directories
- **Sandboxed Execution**: Safe command execution environment
- **API Key Protection**: Never logs or exposes API keys

## 🐛 Troubleshooting

### Common Issues

**API Key Not Working:**
```bash
# Verify the key is set
echo $OPENAI_API_KEY

# Test API access
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
  https://api.openai.com/v1/models
```

**Terminal Display Issues:**
- Update your terminal emulator
- Ensure Unicode/UTF-8 support
- Try a different terminal

**Build Failures:**
```bash
cd src
go mod tidy
go mod download
```

For more troubleshooting help, see [CONFIGURATION.md](docs/CONFIGURATION.md#troubleshooting).

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](docs/CONTRIBUTING.md) for guidelines.

### Quick Start for Contributors

```bash
# Fork and clone
git clone https://github.com/your-username/menace-cli.git
cd menace-cli/menace-cli

# Create a feature branch
git checkout -b feature/your-feature

# Make changes and test
cd src && go run .

# Submit a pull request
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Excellent TUI framework
- [LangChain Go](https://github.com/tmc/langchaingo) - LLM integration library
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling library

## 🔗 Links

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/your-org/menace-cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/menace-cli/discussions)

---

**Made with ❤️ by the Menace CLI team**
