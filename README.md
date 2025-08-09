# Menace CLI

> An intelligent terminal assistant powered by AI

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-1.24.1-blue.svg)](https://golang.org)
[![Node Version](https://img.shields.io/badge/node-18%2B-green.svg)](https://nodejs.org)

Menace CLI is a cross-platform, terminal-based intelligent assistant that combines the power of large language models with a rich text user interface. Built with Go and Node.js, it provides an interactive TUI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea) and integrates with multiple LLM providers including OpenAI, Anthropic, and Ollama.

Transform your terminal experience with AI-powered command suggestions, code assistance, file operations, and GitHub integration—all through natural language conversations.

## ✨ Key Features

- 🤖 **AI-Powered Assistance** - Natural language command generation and code help
- 🖥️ **Cross-Platform** - Works seamlessly on Linux, macOS, and Windows
- ⌨️ **Rich Terminal UI** - Interactive interface with mouse support and keyboard shortcuts
- 📁 **Safe File Operations** - Read, write, and edit files with preview and approval workflows
- 🔧 **GitHub Integration** - Create pull requests, manage repositories, and commit changes
- 🔌 **Multi-Provider Support** - OpenAI, Anthropic, and local models via Ollama
- 🚀 **Easy Installation** - Single npm command installs the right binary for your platform
- 🔒 **Security First** - Command approval system and safe execution practices

## 🚀 Quick Start

### Installation

```bash
# Install globally via npm (recommended)
npm install -g menace-cli

# Set your API key
export OPENAI_API_KEY="your_openai_api_key_here"

# Launch Menace CLI
menace
```

### Alternative: From Source

```bash
git clone https://github.com/your-org/menace-cli.git
cd menace-cli
npm run build
npm link
```

## 🎯 Usage Examples

### Basic Conversations
```
> Help me find all Python files in this directory
> Create a new React component called UserProfile
> Explain what this Git command does: git rebase -i HEAD~3
```

### File Operations
```
> Read the package.json file and show me the dependencies
> Edit the README file to add installation instructions  
> Create a new configuration file for the database connection
```

### GitHub Integration
```
> Stage my changes and create a commit with a good message
> Create a pull request for this feature branch
> Push my changes to the remote repository
```

## 📚 Documentation

| Document | Description |
|----------|-------------|
| [User Guide](docs/USER_GUIDE.md) | Complete installation and usage guide |
| [API Reference](docs/API.md) | Technical documentation and architecture |
| [Developer Guide](docs/DEVELOPER_GUIDE.md) | Development setup and contribution guide |
| [Build System](docs/BUILD_SYSTEM.md) | Build process and cross-platform compilation |
| [Contributing](CONTRIBUTING.md) | How to contribute to the project |

## 🏗️ Architecture

Menace CLI uses a hybrid architecture combining the best of both worlds:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Terminal UI   │◄──►│  Agent Logic    │◄──►│   LLM APIs      │
│  (Bubble Tea)   │    │  (Go Backend)   │    │ (Multi-provider)│
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

- **Frontend**: Rich TUI built with Bubble Tea (Go)
- **Backend**: LLM integration and command processing (Go)  
- **Distribution**: Cross-platform packaging (Node.js)
- **AI Providers**: OpenAI, Anthropic, and Ollama support

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `OPENAI_API_KEY` | OpenAI API key for GPT models | Yes (default provider) |
| `ANTHROPIC_API_KEY` | Anthropic API key for Claude models | No |

### Supported Models

**OpenAI:**
- GPT-4, GPT-4 Turbo
- GPT-3.5 Turbo
- o4-mini (default)

**Anthropic:**
- Claude-3, Claude-3.5-sonnet
- Claude-3-haiku

**Ollama (Local):**
- Any model supported by Ollama
- Requires local Ollama installation

## 🎮 Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Ctrl+C`, `Ctrl+Q` | Exit application |
| `Ctrl+L`, `Esc` | Clear input |
| `Ctrl+M` | Open model selection |
| `Ctrl+A` | Select all text |
| `Ctrl+C` | Copy selection |
| `Ctrl+V` | Paste from clipboard |
| `Page Up/Down` | Scroll chat history |

## 🔧 Development

### Prerequisites

- Go 1.18+ 
- Node.js 14+
- npm

### Local Development

```bash
# Clone the repository
git clone https://github.com/your-org/menace-cli.git
cd menace-cli

# Install Go dependencies
cd src && go mod download

# Run from source (fastest for development)
go run .

# Or build all platforms
npm run build
```

### Project Structure

```
menace-cli/
├── bin/                   # Compiled binaries and Node.js wrapper
├── docs/                  # Documentation
├── scripts/               # Build automation scripts  
├── src/                   # Go source code
│   ├── main.go           # Application entry point
│   ├── llmServer/        # AI and command processing
│   └── ui/               # Terminal user interface
├── CONTRIBUTING.md       # Contribution guidelines
└── package.json          # npm configuration
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- Setting up the development environment
- Code style and standards
- Testing requirements
- Pull request process

### Quick Contribution Steps

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Commit: `git commit -m 'Add amazing feature'`
5. Push: `git push origin feature/amazing-feature`
6. Open a Pull Request

## 📋 Roadmap

- [ ] Plugin system for extensibility
- [ ] Custom model configuration
- [ ] Cloud synchronization
- [ ] Enhanced GitHub integration
- [ ] Docker support
- [ ] Package manager integrations (Homebrew, Chocolatey)

## 🐛 Troubleshooting

### Common Issues

**"OPENAI_API_KEY not set"**
```bash
export OPENAI_API_KEY="sk-your-key-here"
```

**Build fails on M1 Mac**
```bash
# Ensure you have the right Go version
go version  # Should be 1.18+
npm run build-mac
```

**Command execution denied**
- Check file permissions
- Review command before approval
- Ensure you're in the correct directory

See our [User Guide](docs/USER_GUIDE.md) for detailed troubleshooting.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Elegant TUI framework
- [LangChain Go](https://github.com/tmc/langchaingo) - LLM integration library
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

---

**Made with ❤️ for developers who love the terminal**

*Star ⭐ this repository if you find it helpful!*