# Changelog

All notable changes to Menace CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive documentation suite including API docs, architecture guide, and user manual
- LLM integration documentation with provider comparison and setup guides
- Developer contribution guidelines and coding standards
- Cross-platform build system with automated binary generation

### Changed
- Improved project structure with organized documentation
- Enhanced README with clearer installation and usage instructions

### Fixed
- Documentation structure and organization

## [1.0.0] - 2024-XX-XX

### Added
- Initial release of Menace CLI
- Cross-platform terminal user interface built with Bubble Tea
- Multi-LLM provider support (OpenAI, Anthropic, Ollama)
- Intelligent command suggestion system with safety approval
- File operations with structured function calling
- GitHub integration for automated pull request creation
- Real-time model switching and provider selection
- Advanced text editing with clipboard support
- Context-aware system prompts adapted to shell environment
- Shell detection for bash, zsh, PowerShell, and cmd
- Mouse support and responsive UI design
- Streaming responses for better user experience
- Local model support via Ollama integration

### Features
- **Multi-Provider LLM Support**: OpenAI GPT models, Anthropic Claude, and local Ollama models
- **Smart Command Execution**: AI suggests shell commands with safety approval system  
- **File Management**: Read, write, and modify files through natural language
- **Git Integration**: Automated pull request creation and repository management
- **Cross-Platform**: Native binaries for Windows, macOS (Intel/ARM), and Linux
- **Advanced UI**: Rich terminal interface with text selection, scrolling, and mouse support
- **Model Selection**: Runtime switching between different AI models and providers
- **Context Management**: Maintains conversation history and project awareness
- **Safety Features**: Human approval required for potentially dangerous operations

### Technical Details
- Built with Go for high performance and cross-platform compatibility
- Node.js wrapper for easy npm distribution and platform detection  
- Bubble Tea framework for sophisticated terminal UI
- LangChain Go for standardized LLM provider integration
- Thread-safe concurrent operations
- Efficient memory management and context handling

---

## Release Notes Template

Use this template for future releases:

## [X.Y.Z] - YYYY-MM-DD

### Added
- New features that were added in this release

### Changed  
- Changes to existing functionality

### Deprecated
- Features that are being phased out

### Removed
- Features that were removed in this release

### Fixed
- Bug fixes and corrections

### Security
- Security improvements and vulnerability fixes

---

## Release Categories

### Major Releases (X.0.0)
Breaking changes that require user action:
- API changes that break existing functionality
- Removal of deprecated features  
- Major architecture changes
- Significant behavior modifications

### Minor Releases (X.Y.0)
New features and enhancements (backward compatible):
- New LLM provider integrations
- Additional command capabilities
- UI improvements and new features
- Performance enhancements
- New configuration options

### Patch Releases (X.Y.Z)  
Bug fixes and small improvements (backward compatible):
- Bug fixes
- Security patches
- Documentation improvements
- Minor UI tweaks
- Dependency updates

---

## Contribution Guidelines for Changelog

### When to Update
- **Every Pull Request**: Update the [Unreleased] section
- **Before Release**: Move [Unreleased] items to the new version section
- **Breaking Changes**: Always document in the changelog and README

### How to Write Entries
- **Be Descriptive**: Explain what changed and why it matters to users
- **Use Present Tense**: "Add feature" not "Added feature"  
- **Reference Issues**: Link to GitHub issues when relevant
- **Group Related Changes**: Combine related changes into single entries when appropriate

### Good Examples
```markdown
### Added
- Multi-provider model selection with runtime switching between OpenAI, Anthropic, and Ollama models
- Advanced file operations with structured function calling for reading, writing, and modifying files
- GitHub integration enabling automated pull request creation with branch management

### Changed
- Improved command suggestion safety system with clearer approval prompts and risk categorization
- Enhanced UI responsiveness with better text wrapping and scroll performance for long conversations

### Fixed
- Resolved memory leak in conversation history management for long-running sessions  
- Fixed cursor positioning issues with Unicode characters in multi-line input
- Corrected shell detection on Windows systems with PowerShell Core
```

### Poor Examples
```markdown
### Added
- New stuff
- Fixed things
- Better performance

### Changed
- Updated code
- Made improvements
```

---

## Versioning Strategy

### Version Number Guidelines

**1.0.0 → 2.0.0** (Major)
- Remove support for deprecated LLM provider
- Change command-line interface significantly
- Modify configuration file format incompatibly

**1.0.0 → 1.1.0** (Minor)  
- Add support for new LLM provider (Azure OpenAI)
- Introduce new function calling capabilities
- Add configuration options for model parameters

**1.0.0 → 1.0.1** (Patch)
- Fix bug in command parsing
- Update dependencies for security
- Correct documentation errors

### Pre-release Versions
- **Alpha**: `1.1.0-alpha.1` - Early development, unstable
- **Beta**: `1.1.0-beta.1` - Feature complete, testing phase
- **Release Candidate**: `1.1.0-rc.1` - Final testing before release

### Release Schedule
- **Major Releases**: Annually or when breaking changes accumulate
- **Minor Releases**: Quarterly or when significant features are ready
- **Patch Releases**: As needed for important bug fixes

---

## Migration Guides

When creating releases with breaking changes, include migration guides:

### Example Migration Guide (v2.0.0)

#### Configuration Changes
The configuration file format has changed in v2.0.0:

**Before (v1.x.x)**:
```json
{
  "default_model": "gpt-3.5-turbo",
  "provider": "openai"
}
```

**After (v2.0.0)**:
```json
{
  "providers": {
    "openai": {
      "default_model": "gpt-3.5-turbo",
      "api_key_env": "OPENAI_API_KEY"
    }
  },
  "default_provider": "openai"
}
```

#### API Changes
- `Agent.SetModel(model)` → `Agent.SetModel(provider, model)`
- `ProcessQuery(query)` → `ProcessQuery(query, options)`

#### Migration Steps
1. Update configuration file format
2. Update any custom integrations to use new API
3. Test functionality with new provider selection system

---

This changelog provides transparency into the development process and helps users understand what's changing in each release. It serves as both a historical record and a communication tool for the community.