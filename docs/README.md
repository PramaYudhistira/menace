# Menace CLI Documentation

Welcome to the comprehensive documentation for Menace CLI, a cross-platform terminal-based intelligent assistant.

## 📚 Documentation Index

### Getting Started
- [Main README](../README.md) - Project overview, installation, and basic usage
- [Troubleshooting Guide](troubleshooting.md) - Common issues and solutions

### Technical Documentation
- [Architecture](architecture.md) - System design and component overview
- [API Reference](api-reference.md) - Internal Go APIs and function documentation
- [System Prompts](system-prompts.md) - LLM prompt engineering and capabilities

### Development
- [Extending Menace CLI](extending.md) - Adding new providers, functions, and features

## 🚀 Quick Navigation

### For Users
- **New to Menace CLI?** Start with the [Main README](../README.md)
- **Having issues?** Check the [Troubleshooting Guide](troubleshooting.md)
- **Want to understand how it works?** Read the [Architecture](architecture.md)

### For Developers
- **Contributing code?** Review [API Reference](api-reference.md) and [Architecture](architecture.md)
- **Adding new features?** Follow the [Extension Guide](extending.md)
- **Working with prompts?** Study [System Prompts](system-prompts.md)

### For AI/ML Engineers
- **Understanding LLM integration?** See [System Prompts](system-prompts.md)
- **Optimizing responses?** Review prompt engineering patterns
- **Adding new providers?** Follow [Extending Guide](extending.md)

## 🏗️ Architecture Overview

Menace CLI is built with a layered architecture:

```
┌─────────────────────────────────────────┐
│          User Interface (TUI)           │  ← Bubble Tea framework
├─────────────────────────────────────────┤
│        Application Logic & State        │  ← Event handling & UI state
├─────────────────────────────────────────┤
│         LLM Server & Agent              │  ← Conversation management
├─────────────────────────────────────────┤
│    Platform Layer (Shell, File, Git)   │  ← System integrations
└─────────────────────────────────────────┘
```

## 🛠️ Key Technologies

- **Language**: Go 1.18+ (core application)
- **UI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) (TUI)
- **LLM Integration**: [LangChainGo](https://github.com/tmc/langchaingo)
- **Packaging**: Node.js wrapper for cross-platform distribution
- **Supported Providers**: OpenAI, Anthropic, Ollama

## 📖 Feature Documentation

### Core Features

| Feature | Description | Documentation |
|---------|-------------|---------------|
| **Multi-Provider LLM** | OpenAI, Anthropic, Ollama support | [API Reference](api-reference.md#llmserver-package) |
| **Structured Responses** | Command suggestions and function calls | [System Prompts](system-prompts.md#response-format-specifications) |
| **File Operations** | Read, write, diff-based editing | [API Reference](api-reference.md#function-call-system) |
| **Git Integration** | Staging, commits, pull requests | [System Prompts](system-prompts.md#github-integration) |
| **Cross-Platform** | Windows, macOS, Linux support | [Architecture](architecture.md#cross-platform-architecture) |
| **Interactive TUI** | Terminal user interface with mouse support | [Architecture](architecture.md#user-interface-layer-ui) |

### Advanced Capabilities

| Capability | Use Case | Learn More |
|------------|----------|------------|
| **Command Approval** | Safety workflow for destructive operations | [System Prompts](system-prompts.md#safety-and-approval-mechanisms) |
| **Context Injection** | Shell-aware and directory-aware prompts | [System Prompts](system-prompts.md#dynamic-context-injection) |
| **Diff System** | Precise line-level file editing | [System Prompts](system-prompts.md#diff-system-engineering) |
| **Function Calls** | Structured operations beyond text responses | [API Reference](api-reference.md#function-call-system) |

## 🔧 Development Workflow

### Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/your-org/menace-cli.git
cd menace-cli

# Install dependencies and build
npm run build

# Link for development
npm link

# Set up API keys
export OPENAI_API_KEY="your-key-here"
export ANTHROPIC_API_KEY="your-key-here"  # Optional

# Run in debug mode
export MENACE_DEBUG=true
menace
```

### Project Structure

```
menace-cli/
├── docs/                    # Documentation (you are here)
├── menace-cli/
│   ├── bin/                # Distribution wrapper
│   ├── scripts/            # Build scripts
│   └── src/               # Go source code
│       ├── main.go        # Application entry point
│       ├── llmServer/     # LLM integration layer
│       └── ui/            # Terminal user interface
├── LICENSE                # MIT License
└── README.md             # Main project documentation
```

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Read the documentation** - Especially [Architecture](architecture.md) and [API Reference](api-reference.md)
2. **Check existing issues** - See what's already being worked on
3. **Follow the extension guide** - For adding new features, see [Extending](extending.md)
4. **Write tests** - Include unit tests for new functionality
5. **Update documentation** - Keep docs in sync with code changes

### Contribution Areas

- **LLM Providers**: Add support for new AI providers
- **Functions**: Extend the function call system
- **UI Components**: Improve the terminal interface
- **Platform Support**: Enhance cross-platform compatibility
- **Documentation**: Improve guides and examples

## 📊 Performance Considerations

### Resource Usage
- **Memory**: ~50MB typical usage (varies with conversation length)
- **CPU**: Minimal during idle, peaks during LLM requests
- **Network**: Depends on LLM provider and request frequency

### Optimization Tips
- Use local models (Ollama) for better latency
- Clear conversation history periodically
- Choose appropriate models for your use case

## 🔐 Security Considerations

### API Key Management
- Never commit API keys to version control
- Use environment variables for configuration
- Rotate keys regularly

### Command Execution
- All destructive operations require approval
- Input validation prevents common attacks
- Shell escaping protects against injection

## 🆘 Getting Help

### Documentation Order
1. **Start here**: [Main README](../README.md)
2. **Issues**: [Troubleshooting Guide](troubleshooting.md)
3. **Deep dive**: [Architecture](architecture.md)
4. **Development**: [API Reference](api-reference.md) and [Extending](extending.md)

### Support Channels
- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: Questions and community support
- **Documentation**: Search these docs first

### Reporting Issues
Include:
- Menace CLI version
- Operating system
- Full error messages
- Steps to reproduce
- Debug logs (redact sensitive data)

---

## 📄 License

Menace CLI is released under the [MIT License](../LICENSE).

---

**Happy coding with Menace CLI!** 🚀