# Menace CLI - Deployment & Configuration Guide

## Table of Contents

1. [Production Deployment](#production-deployment)
2. [Configuration Management](#configuration-management)
3. [Environment Setup](#environment-setup)
4. [Advanced Configuration](#advanced-configuration)
5. [Security Considerations](#security-considerations)
6. [Performance Tuning](#performance-tuning)
7. [Monitoring and Logging](#monitoring-and-logging)
8. [Troubleshooting](#troubleshooting)

## Production Deployment

### Binary Distribution

Menace CLI is distributed as pre-compiled binaries for multiple platforms through npm.

#### Supported Platforms

- **Linux**: `x86_64` (AMD64)
- **macOS**: `x86_64` (Intel) and `arm64` (Apple Silicon) 
- **Windows**: `x86_64` (AMD64)

#### Global Installation

```bash
# Install globally via npm
npm install -g menace-cli

# Verify installation
menace --version
which menace
```

#### Local Installation

```bash
# Install in project
npm install menace-cli

# Run with npx
npx menace
```

### From Source Deployment

For custom builds or unsupported platforms:

```bash
# Clone repository
git clone https://github.com/PramaYudhistira/terminal-chat.git
cd terminal-chat/menace-cli

# Install dependencies
cd src && go mod tidy

# Build for current platform
go build -o menace-go .

# Or build for all platforms
npm run build

# Install globally
npm link
```

### Docker Deployment

Create a containerized deployment:

```dockerfile
# Dockerfile
FROM node:18-alpine

# Install Go for building
RUN apk add --no-cache go git

# Set working directory
WORKDIR /app

# Copy source
COPY . .

# Build the application
WORKDIR /app/menace-cli
RUN npm install
RUN npm run build

# Set up runtime
FROM node:18-alpine
WORKDIR /app
COPY --from=0 /app/menace-cli .

# Install globally
RUN npm link

# Default command
CMD ["menace"]
```

Build and run:

```bash
docker build -t menace-cli .
docker run -it -e OPENAI_API_KEY="your-key" menace-cli
```

## Configuration Management

### Environment Variables

#### Required Variables

```bash
# OpenAI API Key (required for OpenAI models)
export OPENAI_API_KEY="sk-your-openai-key-here"
```

#### Optional Variables

```bash
# Anthropic API Key (for Claude models)
export ANTHROPIC_API_KEY="your-anthropic-key"

# Custom OpenAI endpoint
export OPENAI_BASE_URL="https://your-custom-endpoint.com"

# Debug mode
export DEBUG="1"

# Custom config directory
export MENACE_CONFIG_DIR="$HOME/.config/menace"

# Default model
export MENACE_DEFAULT_MODEL="gpt-4"

# Default provider
export MENACE_DEFAULT_PROVIDER="openai"
```

### Configuration Files

Create a configuration directory:

```bash
mkdir -p ~/.config/menace
```

#### Main Configuration (`~/.config/menace/config.json`)

```json
{
  "defaultProvider": "openai",
  "defaultModel": "gpt-4o-mini",
  "apiKeys": {
    "openai": "sk-your-key",
    "anthropic": "your-anthropic-key"
  },
  "ui": {
    "theme": "dark",
    "autoExecuteCommands": false,
    "showLineNumbers": true,
    "maxHistoryLength": 1000
  },
  "features": {
    "githubIntegration": true,
    "fileOperations": true,
    "commandSuggestions": true
  },
  "security": {
    "confirmDestructiveCommands": true,
    "allowFileWrites": true,
    "allowNetworkAccess": true
  }
}
```

#### Model Configuration (`~/.config/menace/models.json`)

```json
{
  "providers": {
    "openai": {
      "baseURL": "https://api.openai.com/v1",
      "models": [
        {
          "id": "gpt-4",
          "name": "GPT-4",
          "description": "Most capable model",
          "maxTokens": 8192,
          "costPer1kTokens": 0.03
        },
        {
          "id": "gpt-3.5-turbo",
          "name": "GPT-3.5 Turbo", 
          "description": "Fast and efficient",
          "maxTokens": 4096,
          "costPer1kTokens": 0.002
        }
      ]
    },
    "anthropic": {
      "baseURL": "https://api.anthropic.com",
      "models": [
        {
          "id": "claude-3-sonnet-20240229",
          "name": "Claude 3 Sonnet",
          "description": "Balanced performance",
          "maxTokens": 200000
        }
      ]
    },
    "ollama": {
      "baseURL": "http://localhost:11434",
      "models": [
        {
          "id": "codellama",
          "name": "Code Llama",
          "description": "Local code assistant",
          "local": true
        }
      ]
    }
  }
}
```

## Environment Setup

### Development Environment

```bash
#!/bin/bash
# setup-dev.sh

# Set development API keys
export OPENAI_API_KEY="sk-dev-key-here"
export ANTHROPIC_API_KEY="dev-anthropic-key"

# Enable debug mode
export DEBUG="1"

# Use development endpoints if needed
# export OPENAI_BASE_URL="https://dev-api.openai.com/v1"

# Set development config directory
export MENACE_CONFIG_DIR="./config-dev"

echo "Development environment configured"
```

### Production Environment

```bash
#!/bin/bash
# setup-prod.sh

# Production API keys (consider using secrets management)
export OPENAI_API_KEY="$(cat /etc/secrets/openai-key)"
export ANTHROPIC_API_KEY="$(cat /etc/secrets/anthropic-key)"

# Production configuration
export MENACE_CONFIG_DIR="/etc/menace"

# Disable debug mode
unset DEBUG

# Set resource limits
ulimit -v 1048576  # 1GB virtual memory limit
ulimit -n 1024     # File descriptor limit

echo "Production environment configured"
```

### CI/CD Environment

```yaml
# .github/workflows/deploy.yml
name: Deploy Menace CLI

on:
  release:
    types: [created]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          registry-url: 'https://registry.npmjs.org'
      
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.20'
      
      - name: Install dependencies
        run: |
          cd menace-cli
          npm install
          cd src && go mod tidy
      
      - name: Build all platforms
        run: |
          cd menace-cli
          npm run build
      
      - name: Run tests
        run: |
          # Add your tests here
          cd menace-cli/src
          go test ./...
      
      - name: Publish to npm
        run: |
          cd menace-cli
          npm publish
        env:
          NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}
```

## Advanced Configuration

### Custom LLM Endpoints

#### OpenAI-Compatible APIs

```bash
# Use with OpenAI-compatible APIs
export OPENAI_BASE_URL="https://api.your-provider.com/v1"
export OPENAI_API_KEY="your-provider-key"
```

#### Ollama Local Models

```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Start Ollama service
ollama serve

# Pull models
ollama pull codellama
ollama pull mistral

# Configure Menace to use local models
# (Model selection will automatically detect running Ollama)
```

### Custom System Prompts

Create custom system prompts for specific use cases:

```bash
# Create custom prompts directory
mkdir -p ~/.config/menace/prompts

# Create specialized prompt
cat > ~/.config/menace/prompts/code-review.txt << EOF
You are a senior code reviewer. Focus on:
- Code quality and best practices
- Security vulnerabilities
- Performance optimizations
- Documentation quality
Always provide specific, actionable feedback.
EOF
```

### Plugin System (Future)

Menace CLI is designed to support plugins:

```json
{
  "plugins": [
    {
      "name": "git-assistant",
      "path": "~/.config/menace/plugins/git-assistant.js",
      "enabled": true
    },
    {
      "name": "docker-helper", 
      "path": "~/.config/menace/plugins/docker-helper.js",
      "enabled": false
    }
  ]
}
```

## Security Considerations

### API Key Management

#### Best Practices

1. **Use environment variables** instead of hardcoded keys
2. **Rotate keys regularly**
3. **Use different keys** for different environments
4. **Monitor usage** and set billing alerts

#### Secrets Management

```bash
# Using HashiCorp Vault
vault kv put secret/menace-cli \
  openai_key="sk-your-key" \
  anthropic_key="your-key"

# Retrieve in startup script
export OPENAI_API_KEY=$(vault kv get -field=openai_key secret/menace-cli)
```

#### Docker Secrets

```yaml
# docker-compose.yml
version: '3.8'
services:
  menace-cli:
    image: menace-cli:latest
    environment:
      - OPENAI_API_KEY_FILE=/run/secrets/openai_key
    secrets:
      - openai_key

secrets:
  openai_key:
    external: true
```

### File System Security

#### Restrict File Access

```json
{
  "security": {
    "allowedDirectories": [
      "/home/user/projects",
      "/tmp/menace-workspace"
    ],
    "blockedDirectories": [
      "/etc",
      "/root",
      "/.ssh"
    ],
    "maxFileSize": "10MB",
    "allowedExtensions": [".py", ".js", ".go", ".md", ".txt"]
  }
}
```

#### Sandboxing

```bash
# Run in restricted environment
firejail --private --net=none menace
```

### Network Security

#### Firewall Configuration

```bash
# Allow only necessary outbound connections
iptables -A OUTPUT -d api.openai.com -p tcp --dport 443 -j ACCEPT
iptables -A OUTPUT -d api.anthropic.com -p tcp --dport 443 -j ACCEPT
iptables -A OUTPUT -p tcp --dport 443 -j DROP
```

## Performance Tuning

### Memory Management

```bash
# Set Go garbage collection target
export GOGC=100

# Limit memory usage
ulimit -v 2097152  # 2GB limit
```

### Response Optimization

```json
{
  "performance": {
    "maxTokens": 2048,
    "temperature": 0.7,
    "timeoutSeconds": 30,
    "retryAttempts": 3,
    "cacheResponses": true,
    "maxCacheSize": "100MB"
  }
}
```

### Concurrent Operations

```go
// Limit concurrent API calls
var semaphore = make(chan struct{}, 5) // Max 5 concurrent calls

func (a *Agent) SendMessage(ctx context.Context, input string) (string, *CommandSuggestion, error) {
    semaphore <- struct{}{} // Acquire
    defer func() { <-semaphore }() // Release
    
    // Rest of implementation
}
```

## Monitoring and Logging

### Structured Logging

```go
import (
    "log/slog"
    "os"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    
    slog.SetDefault(logger)
    
    slog.Info("Menace CLI starting", "version", version, "platform", runtime.GOOS)
}
```

### Metrics Collection

```bash
# Log API usage metrics
cat > ~/.config/menace/metrics.json << EOF
{
  "apiCalls": {
    "openai": {
      "count": 142,
      "tokensUsed": 12847,
      "cost": 0.42
    },
    "anthropic": {
      "count": 23,
      "tokensUsed": 5630,
      "cost": 0.15
    }
  },
  "features": {
    "commandsExecuted": 67,
    "filesModified": 12,
    "pullRequestsCreated": 3
  }
}
EOF
```

### Health Checks

```bash
#!/bin/bash
# health-check.sh

# Test basic functionality
echo "Hello" | timeout 10 menace > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "Menace CLI: OK"
    exit 0
else
    echo "Menace CLI: FAILED"
    exit 1
fi
```

## Troubleshooting

### Common Issues

#### High Memory Usage

```bash
# Monitor memory usage
ps aux | grep menace
top -p $(pgrep menace)

# Solutions:
# 1. Reduce conversation history length
# 2. Use smaller models
# 3. Restart periodically
```

#### API Rate Limiting

```bash
# Check rate limits
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
     "https://api.openai.com/v1/usage" | jq .

# Solutions:
# 1. Implement exponential backoff
# 2. Use multiple API keys
# 3. Switch to different providers
```

#### Build Issues

```bash
# Clean and rebuild
cd menace-cli
rm -rf node_modules bin
npm install
npm run build

# Check Go dependencies
cd src
go mod tidy
go mod verify
```

### Log Analysis

```bash
# Enable verbose logging
DEBUG=1 menace 2>&1 | tee menace.log

# Analyze logs
grep ERROR menace.log
grep "API call" menace.log | wc -l
```

### Performance Profiling

```go
// Add profiling endpoint
import _ "net/http/pprof"
import "net/http"

func init() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}
```

Then profile:

```bash
go tool pprof http://localhost:6060/debug/pprof/profile
go tool pprof http://localhost:6060/debug/pprof/heap
```

---

This deployment guide provides a comprehensive foundation for running Menace CLI in production environments. Adjust configurations based on your specific requirements and security policies.