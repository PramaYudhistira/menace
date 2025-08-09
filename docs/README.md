# Documentation Index

Welcome to the Menace CLI documentation! This directory contains comprehensive guides for users, developers, and contributors.

## 📚 Documentation Overview

### For Users
- **[User Guide](USER_GUIDE.md)** - Complete guide for installing and using Menace CLI
  - Installation methods
  - Basic and advanced usage
  - Keyboard shortcuts
  - Troubleshooting
  - Best practices

### For Developers  
- **[API Documentation](API.md)** - Technical API reference and architecture overview
  - Core components and data structures
  - Function signatures and usage
  - Message types and workflows
  - Integration patterns

- **[Developer Guide](DEVELOPER_GUIDE.md)** - In-depth development documentation
  - Architecture deep dive
  - Code organization and patterns
  - State management
  - Performance considerations
  - Security best practices

- **[Build System](BUILD_SYSTEM.md)** - Complete build and deployment documentation
  - Cross-platform compilation
  - Build scripts and automation
  - CI/CD integration
  - Troubleshooting builds

### For Contributors
- **[Contributing Guidelines](../CONTRIBUTING.md)** - How to contribute to the project
  - Development setup
  - Code standards
  - Testing requirements
  - Pull request process

## 🚀 Quick Start

New to Menace CLI? Start here:

1. **Install**: `npm install -g menace-cli`
2. **Configure**: `export OPENAI_API_KEY="your_key"`
3. **Run**: `menace`
4. **Learn**: Check out the [User Guide](USER_GUIDE.md)

## 🏗️ Architecture at a Glance

Menace CLI combines:
- **Go** - Core application with Bubble Tea TUI
- **Node.js** - Cross-platform packaging and distribution
- **LLM APIs** - OpenAI, Anthropic, and Ollama support
- **GitHub Integration** - Built-in Git and PR workflows

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Terminal UI   │◄──►│  Agent Logic    │◄──►│   LLM APIs      │
│  (Bubble Tea)   │    │  (Go Backend)   │    │ (Multi-provider)│
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 📖 Key Features

- **🤖 AI-Powered**: Intelligent command suggestions and code assistance
- **🖥️ Cross-Platform**: Works on Linux, macOS, and Windows
- **⌨️ Rich TUI**: Interactive terminal interface with mouse support
- **📁 File Operations**: Read, write, and edit files safely
- **🔧 GitHub Integration**: Create PRs and manage repositories
- **🔌 Multi-Provider**: Support for OpenAI, Anthropic, and local models

## 🔗 External Resources

- **Project Repository**: [GitHub](https://github.com/your-org/menace-cli)
- **Issue Tracker**: [GitHub Issues](https://github.com/your-org/menace-cli/issues)
- **NPM Package**: [menace-cli](https://npmjs.com/package/menace-cli)

## 📋 Documentation Status

| Document | Status | Last Updated |
|----------|--------|--------------|
| User Guide | ✅ Complete | Current |
| API Documentation | ✅ Complete | Current |
| Developer Guide | ✅ Complete | Current |
| Build System | ✅ Complete | Current |
| Contributing | ✅ Complete | Current |

## 🤝 Contributing to Documentation

Documentation improvements are always welcome! Please see our [Contributing Guidelines](../CONTRIBUTING.md) for details on:

- Documentation standards
- How to submit improvements
- Review process

## 💡 Need Help?

- **General Questions**: Check the [User Guide](USER_GUIDE.md)
- **Development Issues**: See the [Developer Guide](DEVELOPER_GUIDE.md)  
- **Bug Reports**: Open an [issue](https://github.com/your-org/menace-cli/issues)
- **Feature Requests**: Start a [discussion](https://github.com/your-org/menace-cli/discussions)

---

*This documentation is maintained by the Menace CLI community. Last updated: 2025*