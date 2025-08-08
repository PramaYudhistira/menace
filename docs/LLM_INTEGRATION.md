# LLM Integration Guide

This document provides comprehensive information about Large Language Model (LLM) integration in Menace CLI, including supported providers, model selection, configuration, and advanced usage patterns.

## Table of Contents

1. [Overview](#overview)
2. [Supported Providers](#supported-providers)
3. [Model Configuration](#model-configuration)
4. [Provider-Specific Setup](#provider-specific-setup)
5. [Model Selection Strategies](#model-selection-strategies)
6. [Advanced Features](#advanced-features)
7. [Performance & Cost Optimization](#performance--cost-optimization)
8. [Custom Provider Integration](#custom-provider-integration)
9. [Troubleshooting](#troubleshooting)

## Overview

Menace CLI leverages [LangChain Go](https://github.com/tmc/langchaingo) to provide a unified interface for multiple LLM providers. This abstraction allows seamless switching between different models and providers while maintaining consistent functionality.

### Architecture

```
Menace CLI
    ↓
LangChain Go Interface
    ↓
┌─────────────┬──────────────┬─────────────────┐
│   OpenAI    │  Anthropic   │     Ollama      │
│    API      │     API      │ (Local Models)  │
└─────────────┴──────────────┴─────────────────┘
```

### Key Benefits

- **Provider Flexibility**: Switch between cloud and local models
- **Cost Optimization**: Choose models based on task complexity and budget
- **Privacy Options**: Use local models for sensitive data
- **Performance Tuning**: Select models optimized for specific use cases

## Supported Providers

### 1. OpenAI

**Primary Provider**: Most tested and feature-complete integration.

**Available Models**:
- `gpt-4`: Most capable model for complex reasoning
- `gpt-4-turbo`: Faster GPT-4 with improved efficiency  
- `gpt-3.5-turbo`: Fast and cost-effective for most tasks

**Features**:
- Function calling support
- Large context windows (up to 128k tokens)
- Excellent code generation and analysis
- Strong reasoning capabilities

**Best For**:
- Complex code analysis and debugging
- Detailed explanations and tutorials
- System administration tasks
- Creative problem solving

### 2. Anthropic (Claude)

**Advanced Reasoning**: Exceptional analytical capabilities.

**Available Models**:
- `claude-3-opus-20240229`: Most capable model for complex tasks
- `claude-3-sonnet-20240229`: Balanced performance and cost
- `claude-3-haiku-20240307`: Fast and efficient for simple tasks

**Features**:
- Superior reasoning and analysis
- Excellent code review capabilities
- Strong safety and alignment
- Long context understanding

**Best For**:
- Code review and architecture analysis
- Complex problem decomposition
- Security-focused tasks
- Educational explanations

### 3. Ollama (Local Models)

**Privacy-First**: Run models locally without sending data to external services.

**Popular Models**:
- `llama2`: General-purpose conversational model
- `codellama`: Specialized for code generation
- `mistral`: Efficient performance model
- `phi`: Compact model for resource-constrained environments

**Features**:
- Complete privacy - no data leaves your machine
- No API costs or rate limits
- Works offline
- Customizable model parameters

**Best For**:
- Sensitive or proprietary code
- Offline development
- Cost-sensitive applications
- Learning and experimentation

## Model Configuration

### Environment Variables

Configure providers using environment variables:

```bash
# Required for OpenAI models
export OPENAI_API_KEY="sk-your-openai-key-here"

# Optional for Anthropic models  
export ANTHROPIC_API_KEY="sk-ant-your-anthropic-key-here"

# Optional for GitHub integration
export GITHUB_TOKEN="ghp_your-github-token-here"
```

### Runtime Model Selection

Switch models dynamically within Menace CLI:

1. **Via Configuration Panel**: Press `Tab` to open model selection
2. **Via Commands**: Ask Menace to switch models
   ```
   > Switch to GPT-4 for this complex task
   > Use Claude for code review
   > Switch to local Llama model for privacy
   ```

### Default Model Configuration

Set default preferences by modifying the agent initialization:

```go
// In llmServer/agent.go
func NewAgent(apiKey string) (*Agent, error) {
    agent := &Agent{
        provider: "openai",        // Default provider
        Model:    "gpt-3.5-turbo", // Default model
        // ... other initialization
    }
}
```

## Provider-Specific Setup

### OpenAI Setup

1. **Get API Key**: Visit [OpenAI API Platform](https://platform.openai.com/api-keys)
2. **Set Environment Variable**:
   ```bash
   export OPENAI_API_KEY="sk-proj-your-key-here"
   ```
3. **Verify Access**: Test with a simple query
4. **Check Usage**: Monitor usage at [OpenAI Usage Dashboard](https://platform.openai.com/usage)

**Rate Limits** (as of 2024):
- GPT-3.5 Turbo: 3,500 requests/min
- GPT-4: 500 requests/min  
- Varies by account tier

### Anthropic Setup

1. **Get API Key**: Visit [Anthropic Console](https://console.anthropic.com/)
2. **Set Environment Variable**:
   ```bash
   export ANTHROPIC_API_KEY="sk-ant-your-key-here"
   ```
3. **Configure Credits**: Anthropic uses a credit-based system
4. **Test Connection**: Try a Claude model in Menace

**Rate Limits**:
- Claude-3 Opus: 50 requests/min
- Claude-3 Sonnet: 50 requests/min
- Claude-3 Haiku: 50 requests/min

### Ollama Setup

1. **Install Ollama**:
   ```bash
   # macOS/Linux
   curl -fsSL https://ollama.ai/install.sh | sh
   
   # Windows: Download from https://ollama.ai/download
   ```

2. **Download Models**:
   ```bash
   # General purpose model
   ollama pull llama2
   
   # Code-specialized model
   ollama pull codellama
   
   # Efficient model
   ollama pull mistral
   ```

3. **Start Ollama Service**:
   ```bash
   ollama serve
   ```

4. **Verify Installation**:
   ```bash
   ollama list  # Show installed models
   ```

**System Requirements**:
- **RAM**: 8GB minimum (16GB+ recommended for larger models)
- **Storage**: 4-7GB per model
- **CPU**: Modern x86_64 or ARM64 processor
- **GPU**: Optional but significantly improves performance

## Model Selection Strategies

### Task-Based Selection

#### Code Generation & Analysis

**Simple Tasks** (syntax help, basic functions):
- Primary: `gpt-3.5-turbo`
- Alternative: `claude-3-haiku`, `codellama` (local)

**Complex Tasks** (architecture, debugging):
- Primary: `gpt-4` or `claude-3-opus`
- Alternative: `claude-3-sonnet`

**Code Review**:
- Best: `claude-3-opus` (superior analysis)
- Good: `gpt-4`
- Budget: `claude-3-sonnet`

#### System Administration

**Command Help & Scripting**:
- Primary: `gpt-3.5-turbo`
- Alternative: `llama2` (local)

**Complex Configuration**:
- Primary: `gpt-4`
- Alternative: `claude-3-sonnet`

**Security Analysis**:
- Best: `claude-3-opus`
- Good: `gpt-4`

#### Learning & Education

**Quick Questions**:
- Primary: `gpt-3.5-turbo`
- Local: `mistral`

**Deep Explanations**:
- Best: `claude-3-opus`
- Good: `gpt-4`

### Cost Optimization

**Price Comparison** (approximate, as of 2024):

| Model | Input Cost | Output Cost | Context | Best Use |
|-------|------------|-------------|---------|----------|
| GPT-3.5 Turbo | $0.50/1M | $1.50/1M | 16k | General tasks |
| GPT-4 | $10/1M | $30/1M | 128k | Complex reasoning |
| Claude-3 Haiku | $0.25/1M | $1.25/1M | 200k | Fast tasks |
| Claude-3 Sonnet | $3/1M | $15/1M | 200k | Balanced |
| Claude-3 Opus | $15/1M | $75/1M | 200k | Highest quality |
| Ollama Models | Free | Free | Varies | Local/Private |

**Cost-Saving Tips**:
1. Start with cheaper models (GPT-3.5, Claude Haiku)
2. Use local models for experimentation
3. Clear conversation context frequently (`Ctrl+L`)
4. Be specific in queries to avoid back-and-forth

### Privacy Considerations

**High Privacy Requirements**:
- **Recommended**: Ollama local models only
- **Avoid**: Cloud-based providers

**Moderate Privacy**:
- **Preferred**: Anthropic Claude (strong privacy stance)
- **Acceptable**: OpenAI with data opt-out

**Standard Privacy**:
- Any provider acceptable
- Consider data sensitivity before sharing code

## Advanced Features

### Function Calling

Menace CLI supports structured function calling for tool use:

**Supported Functions**:
```go
// File operations
read_file(path: string) → string
write_file(path: string, content: string, line?: number) → success

// Git operations  
create_pull_request(title: string, description: string, branch?: string) → pr_url

// System operations
execute_command(command: string, safe: boolean) → output
```

**Provider Support**:
- ✅ OpenAI: Full function calling support
- ✅ Anthropic: XML-based tool use  
- ⚠️ Ollama: Limited support (model dependent)

### Context Management

**Context Window Sizes**:
- GPT-3.5 Turbo: 16,384 tokens
- GPT-4: 128,000 tokens
- Claude-3: 200,000 tokens
- Ollama: Model-specific (typically 2k-32k)

**Context Strategies**:
- **Sliding Window**: Automatically manage context size
- **Manual Clear**: Use `Ctrl+L` to reset conversation
- **Context Compression**: Summarize older messages

### Custom System Prompts

System prompts are automatically adapted based on:

1. **Detected Shell**: bash, zsh, PowerShell, cmd
2. **Operating System**: Linux, macOS, Windows
3. **Current Directory**: Project type detection
4. **Available Tools**: Function calling capabilities

Example system prompt adaptation:
```
# For bash on Linux
You are an intelligent terminal assistant for bash on Linux.
Available commands: ls, cd, mkdir, git, npm, python, etc.

# For PowerShell on Windows  
You are an intelligent terminal assistant for PowerShell on Windows.
Available commands: Get-ChildItem, Set-Location, New-Item, git, npm, etc.
```

### Streaming Responses

For better user experience, responses can stream in real-time:

```go
// Enable streaming for supported providers
func (a *Agent) ProcessQueryStream(query string) (<-chan string, error) {
    stream, err := a.llm.GenerateContentStream(a.ctx, messages)
    // Handle streaming response
}
```

**Provider Support**:
- ✅ OpenAI: Full streaming support
- ✅ Anthropic: Full streaming support
- ✅ Ollama: Supported

## Performance & Cost Optimization

### Response Time Optimization

**Fastest Models**:
1. Local Ollama models (no network latency)
2. GPT-3.5 Turbo
3. Claude-3 Haiku
4. GPT-4 Turbo
5. Claude-3 Sonnet
6. GPT-4 / Claude-3 Opus (slowest but most capable)

**Optimization Strategies**:
- **Model Switching**: Use fast models for simple tasks
- **Local First**: Try local models when possible  
- **Concise Queries**: Shorter inputs = faster responses
- **Context Management**: Smaller context = faster processing

### Memory Usage

**Client-Side Memory**:
- Conversation history: ~1MB per 1000 messages
- Model metadata: ~50KB per provider
- UI state: ~10MB maximum

**Local Models (Ollama)**:
- 7B parameter models: ~4GB RAM
- 13B parameter models: ~8GB RAM  
- 70B parameter models: ~40GB RAM

### Network Optimization

**Bandwidth Usage**:
- Typical query: 1-10KB
- Typical response: 1-50KB  
- Function calls: Additional 5-20KB

**Offline Capabilities**:
- Local models work completely offline
- Cloud models require internet connection
- Cached responses available for repeated queries

## Custom Provider Integration

### Adding New Providers

To integrate a new LLM provider:

1. **Implement LangChain Interface**:
```go
type CustomProvider struct {
    apiKey   string
    baseURL  string  
    client   *http.Client
}

func (p *CustomProvider) GenerateContent(
    ctx context.Context,
    messages []llms.MessageContent,
    options ...llms.CallOption,
) (*llms.ContentResponse, error) {
    // Implement provider-specific API calls
    return &llms.ContentResponse{
        Choices: []llms.ContentChoice{{
            Content: response,
        }},
    }, nil
}
```

2. **Add Provider Detection**:
```go
// In llmServer/llmFactory.go
func createLLM(provider, model string) (llms.Model, error) {
    switch provider {
    case "custom":
        return NewCustomProvider(apiKey), nil
    // ... existing cases
    }
}
```

3. **Update UI Integration**:
```go
// Add to available providers list
var availableProviders = []string{
    "openai", "anthropic", "ollama", "custom",
}
```

### Provider Requirements

**Minimum Implementation**:
- Text generation endpoint
- API key authentication
- Error handling
- Rate limit respect

**Advanced Features**:
- Streaming responses
- Function calling
- Token usage tracking
- Model switching

**Testing Checklist**:
- ✅ Basic text generation
- ✅ Error handling  
- ✅ Long conversations
- ✅ Special characters/Unicode
- ✅ Function calling (if supported)
- ✅ Streaming (if supported)

## Troubleshooting

### Common Issues

#### Authentication Errors

**Symptoms**:
```
Error: Invalid API key
Error: Authentication failed
```

**Solutions**:
1. Verify API key format and validity
2. Check environment variable is set correctly
3. Ensure sufficient credits/quota
4. Test API key with provider's CLI tools

#### Rate Limiting

**Symptoms**:
```
Error: Rate limit exceeded
Error: Too many requests
```

**Solutions**:
1. Wait for rate limit reset
2. Switch to a model with higher limits
3. Implement request queuing
4. Upgrade API plan if available

#### Model Unavailability

**Symptoms**:
```
Error: Model not found
Error: Access denied to model
```

**Solutions**:
1. Check if model name is spelled correctly
2. Verify account has access to the model
3. Try alternative model
4. Check provider documentation for model changes

#### Local Model Issues (Ollama)

**Symptoms**:
```
Error: Connection refused
Error: Model not found locally
```

**Solutions**:
1. Ensure Ollama service is running: `ollama serve`
2. Verify model is downloaded: `ollama list`
3. Download missing model: `ollama pull model-name`
4. Check system resources (RAM, disk space)

### Debug Mode

Enable detailed logging:

```bash
# Debug LLM interactions
DEBUG_LLM=1 menace

# Debug all components  
DEBUG=* menace
```

### Performance Monitoring

Track model performance:

```go
// Add metrics collection
type ModelMetrics struct {
    Provider     string
    Model        string
    RequestTime  time.Duration
    TokensUsed   int
    Success      bool
    Error        string
}
```

### Health Checks

Verify provider connectivity:

```bash
# Test OpenAI
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
     "https://api.openai.com/v1/models"

# Test Anthropic  
curl -H "x-api-key: $ANTHROPIC_API_KEY" \
     "https://api.anthropic.com/v1/messages"

# Test Ollama
curl http://localhost:11434/api/tags
```

---

## Best Practices Summary

1. **Start Simple**: Begin with GPT-3.5 Turbo for most tasks
2. **Escalate When Needed**: Switch to more capable models for complex tasks
3. **Consider Privacy**: Use local models for sensitive data
4. **Monitor Costs**: Track API usage and optimize model selection
5. **Test Providers**: Try different providers for the same task to compare results
6. **Keep Updated**: Stay informed about new models and capabilities
7. **Backup Plans**: Have fallback options when primary providers are unavailable

The LLM integration in Menace CLI provides flexibility and choice while maintaining consistent functionality. By understanding the strengths and characteristics of different providers and models, you can optimize for your specific needs regarding performance, cost, privacy, and capabilities.