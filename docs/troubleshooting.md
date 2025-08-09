# Troubleshooting Guide

This guide helps you diagnose and resolve common issues with Menace CLI.

## Common Issues

### Installation and Setup

#### Issue: `OPENAI_API_KEY environment variable not set`

**Symptoms:**
```
Error: OPENAI_API_KEY environment variable not set
```

**Solutions:**

**On macOS/Linux:**
```bash
# Temporary (current session only)
export OPENAI_API_KEY="sk-your-api-key-here"

# Permanent (add to shell profile)
echo 'export OPENAI_API_KEY="sk-your-api-key-here"' >> ~/.bashrc
# or for zsh
echo 'export OPENAI_API_KEY="sk-your-api-key-here"' >> ~/.zshrc

# Reload shell configuration
source ~/.bashrc  # or source ~/.zshrc
```

**On Windows (Command Prompt):**
```cmd
setx OPENAI_API_KEY "sk-your-api-key-here"
```

**On Windows (PowerShell):**
```powershell
[Environment]::SetEnvironmentVariable("OPENAI_API_KEY", "sk-your-api-key-here", "User")
```

**Verification:**
```bash
echo $OPENAI_API_KEY  # Should display your API key
menace                # Should start without error
```

---

#### Issue: `npm install -g menace-cli` fails

**Symptoms:**
- Permission denied errors
- Module not found errors
- Installation timeouts

**Solutions:**

**Permission Issues (macOS/Linux):**
```bash
# Use npm's built-in solution
npm config set prefix '~/.npm-global'
export PATH=~/.npm-global/bin:$PATH

# Or use Node Version Manager (recommended)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install node
npm install -g menace-cli
```

**Windows Permission Issues:**
- Run Command Prompt/PowerShell as Administrator
- Or use package managers like Chocolatey or Scoop

**Network Issues:**
```bash
# Configure npm registry and timeout
npm config set registry https://registry.npmjs.org/
npm config set timeout 60000
npm install -g menace-cli
```

---

#### Issue: Build failures when installing from source

**Symptoms:**
```
go build: command not found
npm run build fails
```

**Solutions:**

**Missing Go Installation:**
```bash
# Install Go (version 1.18+ required)
# macOS with Homebrew
brew install go

# Linux (Ubuntu/Debian)
sudo apt update
sudo apt install golang-go

# Verify installation
go version
```

**Missing Node.js/npm:**
```bash
# Install Node.js (v14+ required)
# macOS with Homebrew
brew install node

# Linux (Ubuntu/Debian)
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs

# Verify installation
node --version
npm --version
```

**Build Process:**
```bash
# Clone and build
git clone https://github.com/your-org/menace-cli.git
cd menace-cli
npm run build
npm link
```

---

### Runtime Issues

#### Issue: `menace` command not found after installation

**Symptoms:**
```bash
menace
# zsh: command not found: menace
# or
# bash: menace: command not found
```

**Solutions:**

**Check Installation:**
```bash
# Verify global installation
npm list -g --depth=0 | grep menace-cli

# Check npm global directory
npm config get prefix

# Verify PATH includes npm global binaries
echo $PATH
```

**Fix PATH Issues:**
```bash
# Add npm global binaries to PATH
echo 'export PATH="$(npm config get prefix)/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

**Alternative Installation:**
```bash
# Uninstall and reinstall
npm uninstall -g menace-cli
npm install -g menace-cli

# Or use npx for one-time execution
npx menace-cli
```

---

#### Issue: LLM provider errors

**Symptoms:**
```
Error: failed to create OpenAI client: unauthorized
Error: failed to create Anthropic client: invalid API key
Connection timeout when contacting LLM provider
```

**Solutions:**

**API Key Issues:**
```bash
# Verify API key format
echo $OPENAI_API_KEY | grep -E "^sk-[A-Za-z0-9]{48}$"
echo $ANTHROPIC_API_KEY | grep -E "^sk-ant-[A-Za-z0-9\-_]{10,}$"

# Test API key validity
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
  https://api.openai.com/v1/models | jq '.data[0].id'
```

**Network/Proxy Issues:**
```bash
# Check internet connectivity
curl -I https://api.openai.com/

# Configure proxy if needed
export HTTPS_PROXY=http://your-proxy:port
export HTTP_PROXY=http://your-proxy:port

# Verify DNS resolution
nslookup api.openai.com
```

**Provider-Specific Setup:**

**Anthropic Setup:**
```bash
export ANTHROPIC_API_KEY="sk-ant-your-key-here"
# Test connectivity
curl -H "Authorization: Bearer $ANTHROPIC_API_KEY" \
  -H "Content-Type: application/json" \
  https://api.anthropic.com/v1/messages \
  -d '{"model":"claude-3-sonnet-20240229","max_tokens":10,"messages":[{"role":"user","content":"Hi"}]}'
```

**Ollama Setup:**
```bash
# Install and start Ollama
curl -fsSL https://ollama.ai/install.sh | sh
ollama serve &
ollama pull llama2

# Test local connection
curl http://localhost:11434/api/generate \
  -d '{"model":"llama2","prompt":"Hello","stream":false}'
```

---

#### Issue: UI rendering problems

**Symptoms:**
- Garbled text display
- Missing characters or colors
- Layout issues
- Mouse interactions not working

**Solutions:**

**Terminal Compatibility:**
```bash
# Check terminal capabilities
echo $TERM
tput colors  # Should show 256 or more

# Try different terminal emulators
# - iTerm2 (macOS)
# - Windows Terminal (Windows)
# - GNOME Terminal (Linux)
```

**Font and Encoding Issues:**
```bash
# Verify UTF-8 support
locale | grep UTF-8

# Set locale if needed
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8
```

**Terminal Size Issues:**
```bash
# Check terminal dimensions
stty size

# Resize terminal or adjust font size
# Minimum recommended: 80x24 characters
```

**Alternative UI Mode:**
```bash
# Try without mouse support
MENACE_NO_MOUSE=true menace

# Try without colors
NO_COLOR=1 menace
```

---

#### Issue: File operation failures

**Symptoms:**
```
Error: permission denied reading file
Error: no such file or directory
Error: failed to apply diffs to file
```

**Solutions:**

**Permission Issues:**
```bash
# Check file permissions
ls -la filename.txt

# Fix permissions if needed
chmod 644 filename.txt  # Read/write for owner, read for others
chmod 755 directory/    # Execute permissions for directories
```

**Path Resolution Issues:**
```bash
# Use absolute paths
/full/path/to/file.txt

# Check current working directory
pwd

# Verify file exists
ls -la /path/to/file.txt
```

**File Locking Issues:**
```bash
# Check for file locks
lsof filename.txt

# Close other applications using the file
# Or wait for operations to complete
```

---

#### Issue: Git integration problems

**Symptoms:**
```
Error: git command not found
Error: not a git repository
Error: failed to create pull request
```

**Solutions:**

**Git Installation:**
```bash
# Install Git
# macOS
brew install git

# Linux (Ubuntu/Debian)
sudo apt update
sudo apt install git

# Verify installation
git --version
```

**Repository Setup:**
```bash
# Initialize repository if needed
git init
git remote add origin https://github.com/user/repo.git

# Configure Git user
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

**GitHub Authentication:**
```bash
# Generate personal access token at https://github.com/settings/tokens
# With repo permissions

# Configure authentication
git config --global credential.helper store
# Then authenticate on first push

# Or use SSH keys
ssh-keygen -t ed25519 -C "your.email@example.com"
# Add public key to GitHub account
```

**Pull Request Issues:**
```bash
# Verify repository permissions
git remote -v
git branch -a

# Check GitHub API access
curl -H "Authorization: token YOUR_GITHUB_TOKEN" \
  https://api.github.com/user
```

---

### Performance Issues

#### Issue: Slow LLM responses

**Symptoms:**
- Long wait times for responses
- Timeouts
- Degraded performance over time

**Solutions:**

**Model Selection:**
- Switch to faster models (e.g., GPT-3.5-turbo instead of GPT-4)
- Use local models with Ollama for better latency
- Consider model-specific optimization settings

**Network Optimization:**
```bash
# Check network latency
ping api.openai.com

# Use regional endpoints if available
# Configure connection pooling and keep-alive
```

**Memory Management:**
```bash
# Clear conversation history
# In Menace CLI, use Ctrl+L or restart the application

# Monitor memory usage
top -p $(pgrep menace)
```

---

#### Issue: High memory usage

**Symptoms:**
- System slowdown
- Out of memory errors
- Swap usage increase

**Solutions:**

**Configuration Adjustments:**
```bash
# Limit conversation history
export MENACE_MAX_HISTORY=100

# Reduce UI refresh rate
export MENACE_REFRESH_RATE=30
```

**Resource Monitoring:**
```bash
# Monitor memory usage
ps aux | grep menace
htop

# Restart application periodically for long sessions
```

---

## Debug Mode

### Enable Debug Logging

```bash
# Enable verbose logging
export MENACE_DEBUG=true
export MENACE_LOG_LEVEL=debug

# Run with debug output
menace 2>&1 | tee menace-debug.log
```

### Log Analysis

**Log Locations:**
- Linux/macOS: `~/.menace/logs/`
- Windows: `%USERPROFILE%\.menace\logs\`

**Common Log Patterns:**
```bash
# Search for errors
grep -i error menace-debug.log

# Find LLM communication issues
grep -i "llm\|openai\|anthropic" menace-debug.log

# Check command execution
grep -i "command\|execute" menace-debug.log
```

## Getting Help

### Information Gathering

When reporting issues, include:

**System Information:**
```bash
# System details
uname -a
echo $SHELL
echo $TERM

# Go and Node versions
go version
node --version
npm --version

# Menace CLI version
menace --version  # If working
npm list -g menace-cli  # Installation info
```

**Configuration:**
```bash
# Environment variables (redact sensitive data)
env | grep -E "(OPENAI|ANTHROPIC|MENACE)"

# Git configuration
git config --list --show-origin
```

**Error Context:**
- Full error messages
- Steps to reproduce
- Expected vs. actual behavior
- Recent changes to system/configuration

### Support Channels

1. **GitHub Issues**: [Repository Issues](https://github.com/your-org/menace-cli/issues)
2. **Documentation**: Check existing docs for solutions
3. **Community**: Search for similar issues in discussions

### Bug Reports

**Include in Bug Reports:**
1. Menace CLI version
2. Operating system and version
3. Terminal emulator and version
4. Full error messages
5. Steps to reproduce
6. Debug logs (with sensitive data redacted)

**Security Issues:**
- Report security vulnerabilities privately
- Don't include API keys or sensitive data in public issues
- Use secure communication channels for sensitive reports

This troubleshooting guide covers the most common issues. For additional help, consult the project documentation or open an issue with detailed information about your problem.