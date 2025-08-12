<div align="center">
  <h1>🚀 Menace CLI</h1>
  <p><em>Your intelligent terminal companion</em></p>
  
  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
  [![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
  [![Node Version](https://img.shields.io/badge/Node-14+-green.svg)](https://nodejs.org/)
  [![Built with Bubble Tea](https://img.shields.io/badge/Built%20with-Bubble%20Tea-purple.svg)](https://github.com/charmbracelet/bubbletea)
</div>

---

**Menace CLI** is a powerful, cross-platform terminal-based AI assistant that transforms your command-line experience. Built with Go and powered by modern LLM APIs, it provides an intuitive TUI (Text User Interface) for intelligent code assistance, command suggestions, and productivity enhancement.

## ✨ Features

- 🖥️ **Cross-Platform**: Native binaries for Windows, macOS (Intel & Apple Silicon), and Linux
- 🎨 **Beautiful TUI**: Interactive interface with mouse support powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- 🤖 **AI-Powered**: Context-aware assistance backed by OpenAI and Anthropic models
- ⚡ **Fast & Lightweight**: Efficient Go backend with minimal Node.js wrapper
- 🔧 **Easy Setup**: Simple installation via npm or build from source
- 🎯 **Smart Suggestions**: Intelligent shell command recommendations
- 📝 **Code Assistance**: Read, write, and navigate code with AI guidance

## 🛠️ Prerequisites

Before getting started, ensure you have:

- **Go** 1.24+ ([Download](https://golang.org/dl/))
- **Node.js** v14+ and npm ([Download](https://nodejs.org/))
- **OpenAI API Key** ([Get one here](https://platform.openai.com/api-keys))

## 📦 Installation

### Option 1: From npm (Recommended)

> 🚧 Coming soon! This will be the easiest way to install Menace CLI.

```bash
npm install --global menace-cli
```

### Option 2: From Source

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-org/menace-cli.git
   cd menace-cli
   ```

2. **Set up your API key**
   ```bash
   # Unix/macOS/Linux
   export OPENAI_API_KEY="your_api_key_here"
   
   # Windows PowerShell
   $env:OPENAI_API_KEY="your_api_key_here"
   ```

3. **Build the application**
   ```bash
   cd menace-cli
   npm run build
   ```

4. **Link globally (optional)**
   ```bash
   npm link
   ```

## 🚀 Usage

Once installed and configured, launch Menace CLI:

```bash
menace
```

The interactive TUI will start, providing you with:
- Command suggestions based on your current context
- Code analysis and writing assistance  
- Intelligent terminal session enhancement

## 🏗️ Development

### Project Structure

```
menace-cli/
├── 📁 bin/                    # Compiled Go binaries per platform
├── 📄 package.json            # NPM metadata & build scripts
├── 📁 scripts/                # Build automation scripts
│   ├── build-all.js           # Cross-platform build
│   ├── build-mac.js           # macOS specific build
│   └── build-windows.js       # Windows specific build
└── 📁 src/                    # Go source code
    ├── 📄 main.go             # Application entry point
    ├── 📁 llmServer/           # LLM integration & AI logic
    │   ├── agent.go           # AI agent implementation
    │   ├── llmFactory.go      # LLM provider factory
    │   ├── llmService.go      # LLM service interface
    │   └── system_prompt.go   # AI prompt management
    └── 📁 ui/                  # Bubble Tea UI components
        ├── model.go           # UI state management
        ├── view.go            # UI rendering
        ├── update.go          # Event handling
        └── executor.go        # Command execution
```

### Build Commands

| Command | Description |
|---------|-------------|
| `npm run build` | 🔨 Build binaries for all platforms |
| `npm run build-windows` | 🪟 Build Windows executable only |
| `npm run build-mac` | 🍎 Build macOS executables (Intel & ARM) |

### Technology Stack

- **Backend**: Go 1.24+ with [LangChain Go](https://github.com/tmc/langchaingo)
- **UI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **AI Integration**: OpenAI API, Anthropic Claude
- **Build System**: Node.js + npm scripts
- **Cross-platform**: Native Go compilation

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `OPENAI_API_KEY` | Your OpenAI API key | ✅ Yes |

### Setting up your API key

**Unix/macOS/Linux:**
```bash
echo 'export OPENAI_API_KEY="sk-your-key-here"' >> ~/.bashrc
source ~/.bashrc
```

**Windows Command Prompt:**
```cmd
setx OPENAI_API_KEY "sk-your-key-here"
```

**Windows PowerShell:**
```powershell
[Environment]::SetEnvironmentVariable("OPENAI_API_KEY", "sk-your-key-here", "User")
```

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feat/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Test** your changes (`npm run build`)
5. **Push** to your branch (`git push origin feat/amazing-feature`)
6. **Open** a Pull Request

### Development Setup

```bash
# Clone your fork
git clone https://github.com/your-username/menace-cli.git
cd menace-cli/menace-cli

# Install dependencies and build
npm install
npm run build

# Test your changes
./bin/menace-darwin-amd64  # or appropriate binary for your platform
```

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

<div align="center">
  <p>Made with ❤️ by the Menace team</p>
  <p>
    <a href="#-menace-cli">Back to top</a> •
    <a href="https://github.com/your-org/menace-cli/issues">Report Bug</a> •
    <a href="https://github.com/your-org/menace-cli/issues">Request Feature</a>
  </p>
</div>
