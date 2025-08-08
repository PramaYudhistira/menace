# Menace CLI

<div align="center">
  <img src="https://via.placeholder.com/200x200?text=Menace+CLI" alt="Menace CLI Logo" width="200"/>
  <p><strong>Your Intelligent Terminal Assistant</strong></p>
  <p>AI-powered command line interface built with Go and Node.js</p>
</div>

---

**Menace CLI** is a sophisticated terminal-based intelligent assistant that transforms your command-line experience. Built with Go for performance and Node.js for distribution, it leverages multiple Large Language Model providers to offer intelligent command suggestions, code assistance, file operations, and automated Git workflows.

## ✨ Key Features

🤖 **Multi-LLM Support** - OpenAI GPT, Anthropic Claude, and local Ollama models  
⚡ **Smart Command Execution** - AI suggests commands with safety approval system  
📁 **File Operations** - Read, write, and modify files through natural language  
🔀 **Git Integration** - Automated pull request creation and repository management  
🖥️ **Cross-Platform** - Native binaries for Windows, macOS (Intel/ARM), and Linux  
🎨 **Rich Terminal UI** - Advanced text editing, mouse support, and responsive design  
🔄 **Runtime Model Switching** - Change AI providers and models on the fly  
🧠 **Context Awareness** - Maintains conversation history and understands your project  
🛡️ **Safety First** - Human approval required for potentially dangerous operations

## 🚀 Quick Start

### Installation

**Option 1: NPM (Recommended)**
```bash
npm install -g menace-cli
```

**Option 2: From Source**
```bash
git clone https://github.com/your-org/menace-cli.git
cd menace-cli
npm run build
npm link
```

### Configuration

Set your API keys:
```bash
# Required
export OPENAI_API_KEY="sk-your-openai-key"

# Optional (for additional providers)
export ANTHROPIC_API_KEY="sk-ant-your-anthropic-key"
export GITHUB_TOKEN="ghp_your-github-token"
```

### Launch Menace

```bash
menace
```

That's it! You now have an AI-powered terminal assistant ready to help.

## 💡 Example Usage

**Get command suggestions:**
```
> How do I list all files recursively?
AI: You can use the find command: find . -type f
[Execute] [Cancel]
```

**File operations:**
```
> Create a new React component for user authentication
AI: I'll create a React component for you...
[Function Call: create_file] auth.jsx created
```

**Git integration:**
```
> Create a pull request for my changes
AI: I'll create a PR with your current changes...
[Pull Request Created] https://github.com/user/repo/pull/123
```

## 🛠️ Supported LLM Providers

| Provider | Models | Features |
|----------|--------|----------|
| **OpenAI** | GPT-4, GPT-4 Turbo, GPT-3.5 Turbo | Function calling, large context |
| **Anthropic** | Claude-3 Opus, Sonnet, Haiku | Superior reasoning, safety-focused |
| **Ollama** | Llama 2, Code Llama, Mistral | Local models, complete privacy |

## 📚 Documentation

- **[User Guide](docs/USER_GUIDE.md)** - Complete usage instructions and tips
- **[LLM Integration](docs/LLM_INTEGRATION.md)** - Provider setup and model selection
- **[Architecture](docs/ARCHITECTURE.md)** - System design and technical details
- **[API Reference](docs/API.md)** - Developer API documentation  
- **[Contributing](docs/CONTRIBUTING.md)** - Development setup and guidelines

## 🏗️ Architecture

```
┌─────────────────────────────────────────┐
│           Terminal User                 │
├─────────────────────────────────────────┤
│  Node.js Wrapper (Cross-platform)      │
├─────────────────────────────────────────┤
│  Go Backend                             │
│  ├─ Bubble Tea UI                      │
│  ├─ LLM Agent System                   │
│  └─ Command Processing                 │
├─────────────────────────────────────────┤
│  LLM Providers                         │
│  ├─ OpenAI API                         │
│  ├─ Anthropic API                      │
│  └─ Ollama (Local)                     │
└─────────────────────────────────────────┘
```

**Tech Stack:**
- **Go** - High-performance backend with Bubble Tea TUI
- **Node.js** - Cross-platform distribution and binary management  
- **LangChain Go** - Unified LLM provider interface
- **Bubble Tea** - Rich terminal user interface framework

## 🔧 Development

### Prerequisites
- Go 1.18+
- Node.js 14+
- Git

### Setup
```bash
git clone https://github.com/your-org/menace-cli.git
cd menace-cli
go mod download
npm install
npm run build
```

### Project Structure
```
├── bin/                   # Cross-platform binaries
├── scripts/               # Build automation  
├── src/                   # Go source code
│   ├── main.go           # Application entry point
│   ├── llmServer/        # LLM integration layer
│   └── ui/               # Terminal UI components
├── docs/                  # Comprehensive documentation
└── tests/                 # Test suites
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](docs/CONTRIBUTING.md) for details.

### Quick Contribution Steps
1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Make your changes
4. Add tests and ensure they pass
5. Commit your changes (`git commit -m 'feat: add amazing feature'`)
6. Push to your branch (`git push origin feat/amazing-feature`)
7. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Excellent TUI framework
- [LangChain Go](https://github.com/tmc/langchaingo) - LLM integration library
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- OpenAI, Anthropic, and Ollama teams for excellent AI models

---

<div align="center">
  <p>Made with ❤️ by the Menace CLI team</p>
  <p>
    <a href="#menace-cli">Back to Top</a> •
    <a href="docs/USER_GUIDE.md">User Guide</a> •
    <a href="docs/CONTRIBUTING.md">Contributing</a>
  </p>
</div>
