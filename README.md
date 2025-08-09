# Menace CLI

Menace CLI is a cross-platform, terminal-based intelligent assistant built with Go and Node.js. It provides an interactive TUI (text user interface) powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [LangChainGo](https://github.com/tmc/langchaingo) for LLM integration. Menace can suggest shell commands, help navigate code, read and write files, manage Git operations, and generally make your terminal sessions more productive.

## 📚 Documentation

- **[Complete Documentation](docs/README.md)** - Full documentation index
- **[Architecture Overview](docs/architecture.md)** - System design and components  
- **[API Reference](docs/api-reference.md)** - Internal APIs and function documentation
- **[Troubleshooting](docs/troubleshooting.md)** - Common issues and solutions
- **[Extension Guide](docs/extending.md)** - Adding new providers and features
- **[System Prompts](docs/system-prompts.md)** - LLM prompt engineering

## Features

### Core Capabilities
- **Cross-platform binaries** for Windows, macOS (Intel & Apple Silicon), and Linux
- **Interactive TUI** with mouse support and responsive design
- **Multi-provider LLM support** - OpenAI, Anthropic, and Ollama
- **Intelligent command suggestions** with approval workflow for safety
- **File operations** - Read, write, and diff-based editing without shell commands
- **Git integration** - Stage, commit, push, and create pull requests
- **Context-aware responses** - Shell detection and working directory awareness

### Advanced Features  
- **Structured function calls** for precise file operations
- **Command approval system** to prevent destructive operations
- **Conversation persistence** with message history management
- **Dynamic model switching** between providers at runtime
- **Cross-platform clipboard** integration for copy/paste operations

## Prerequisites

- Go (1.18+ recommended)
- Node.js (v14+) and npm

## Installation

### 1 From npm (recommended and coming soon!)

```bash
npm install --global menace-cli
```

### 2 From source

1. Clone this repo:
   ```bash
   git clone https://github.com/your-org/menace-cli.git
   cd menace-cli
   ```
2. Set your OpenAI API key:
   ```bash
   export OPENAI_API_KEY="your_key_here"
   ```
3. Build all the platform binaries (into `./bin`):
   ```bash
   npm run build
   ```
4. Link the CLI locally:
   ```bash
   npm link
   ```

## Usage

After installation and setting `OPENAI_API_KEY`, simply run:

```bash
menace
```

## Development

### Directory Layout

```
.
├── bin/                   # Compiled Go binaries per platform
├── package.json           # NPM metadata + build scripts
├── scripts/               # Build scripts (build-all, build-windows, build-mac)
└── src/
    ├── main.go            # TUI bootstrap + OpenAI agent init
    ├── llmServer/         # LLM agent factory & prompt management
    └── ui/                # Bubble Tea models & views
```

## Build Scripts

- `npm run build`  Builds Go binaries for all supported OS/architectures into `./bin`.
- `npm run build-windows`  Builds only the Windows executable.
- `npm run build-mac`  Builds only the macOS executables (Intel & ARM).

## Configuration

### API Keys

Menace CLI supports multiple LLM providers. Set up the API keys for the providers you want to use:

**OpenAI (Default):**
```bash
export OPENAI_API_KEY="sk-your-openai-key-here"
```

**Anthropic (Optional):**
```bash
export ANTHROPIC_API_KEY="sk-ant-your-anthropic-key-here"
```

**Ollama (Optional - for local models):**
```bash
# Install and start Ollama server
curl -fsSL https://ollama.ai/install.sh | sh
ollama serve
ollama pull llama2  # or other models
```

### Windows Setup

**Command Prompt:**
```cmd
setx OPENAI_API_KEY "sk-your-key-here"
setx ANTHROPIC_API_KEY "sk-ant-your-key-here"
```

**PowerShell:**
```powershell
[Environment]::SetEnvironmentVariable("OPENAI_API_KEY", "sk-your-key-here", "User")
[Environment]::SetEnvironmentVariable("ANTHROPIC_API_KEY", "sk-ant-your-key-here", "User")
```

### Provider Selection

Once multiple providers are configured, you can switch between them in the UI or via model selection interface. Each provider offers different models with varying capabilities and pricing.

## Getting Help

### Documentation
- **First time?** Start with this README and the [Complete Documentation](docs/README.md)
- **Having issues?** Check the [Troubleshooting Guide](docs/troubleshooting.md)
- **Want to contribute?** Read the [Extension Guide](docs/extending.md) and [Architecture](docs/architecture.md)

### Support
- **Bug reports**: [GitHub Issues](https://github.com/your-org/menace-cli/issues)
- **Feature requests**: [GitHub Discussions](https://github.com/your-org/menace-cli/discussions)
- **Questions**: Search existing issues and discussions first

### Common Issues
- **API key errors**: Verify environment variables are set correctly
- **Command not found**: Ensure npm global directory is in your PATH
- **UI rendering issues**: Try a different terminal emulator or check font support

## Contributing

We welcome contributions! Please see our [Extension Guide](docs/extending.md) for detailed information.

### Quick Start
1. Fork the repo and clone your fork
2. Set up development environment (see [docs](docs/README.md))
3. Create a feature branch (`git checkout -b feature/amazing-feature`)
4. Make your changes with tests and documentation
5. Ensure builds pass (`npm run build`)
6. Submit a pull request

### Development Areas
- **LLM Providers**: Add support for new AI services
- **Functions**: Extend file operations and integrations  
- **UI/UX**: Improve the terminal interface
- **Platform Support**: Enhance cross-platform compatibility
- **Documentation**: Keep guides current and comprehensive

## License

MIT License
