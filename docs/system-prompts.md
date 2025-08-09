# System Prompts and LLM Engineering

This document explains how Menace CLI engineers prompts to provide intelligent terminal assistance and the capabilities enabled through structured prompting.

## Overview

Menace CLI uses sophisticated system prompts to transform general-purpose LLMs into specialized terminal assistants. The system prompt defines the agent's capabilities, response formats, and operational constraints.

## Prompt Engineering Architecture

### Dynamic Context Injection

The system prompt is dynamically generated with runtime information:

```go
func getSystemPrompt(shell string) string {
    cwd, err := os.Getwd()
    if err != nil {
        cwd = "unknown directory"
    }
    return fmt.Sprintf(promptTemplate, shell, shell, cwd)
}
```

**Injected Context:**
- **Shell Environment**: Platform-specific shell detection (`windows/CMD`, `darwin/bash`, `linux/bash`)
- **Working Directory**: Current working directory for file operations
- **Platform Constraints**: OS-specific command limitations

### Core Agent Capabilities

The system prompt establishes the following agent capabilities:

#### 1. File System Operations
- **File Reading**: With line-numbered output for precise editing
- **File Writing**: Through structured diff operations
- **Directory Navigation**: Context-aware path resolution

#### 2. Command Execution
- **Shell Commands**: Platform-appropriate command execution
- **Safety Constraints**: Prevention of destructive operations
- **Approval Workflow**: User confirmation for sensitive operations

#### 3. Git Integration
- **Repository Operations**: Staging, committing, pushing
- **Branch Management**: Context-aware branch operations
- **Pull Request Creation**: Automated PR generation with GitHub API

## Response Format Specifications

### Command Suggestion Format

The LLM uses structured blocks for command suggestions:

```
[COMMAND_SUGGESTION]
Reason: <explanation of why this command is needed>
Command: <shell command to execute>
AwaitingCommandApproval: <true/false>
[/COMMAND_SUGGESTION]
```

**Example:**
```
[COMMAND_SUGGESTION]
Reason: To list all files in the current directory
Command: ls -la
AwaitingCommandApproval: false
[/COMMAND_SUGGESTION]
```

### Function Call Format

For file operations and API integrations:

```
[FUNCTION_CALL]
Reason: <explanation of why this function is needed>
AwaitingCommandApproval: <true/false>
Payload:
{
    "name": "FunctionName",
    "args": {
        "parameter": "value"
    }
}
[/FUNCTION_CALL]
```

**Example - File Reading:**
```
[FUNCTION_CALL]
Reason: Read the configuration file to understand current settings
AwaitingCommandApproval: false
Payload:
{
    "name": "ReadFileWithLineNumbers",
    "args": {
        "path": "config.json"
    }
}
[/FUNCTION_CALL]
```

**Example - File Editing:**
```
[FUNCTION_CALL]
Reason: Update the version number in package.json
AwaitingCommandApproval: true
Payload:
{
    "name": "CreateAndApplyDiffs",
    "args": {
        "path": "package.json",
        "diffs": [
            {
                "Type": 2,
                "LineIndex": 3,
                "OldContent": "  \"version\": \"1.0.0\",",
                "NewContent": "  \"version\": \"1.1.0\","
            }
        ]
    }
}
[/FUNCTION_CALL]
```

## Diff System Engineering

### Line-Level Editing

The diff system enables precise file modifications without shell redirection:

**Operation Types:**
- `Type: 0` - **Add**: Insert new content at specified line
- `Type: 1` - **Delete**: Remove existing line content
- `Type: 2` - **Modify**: Replace existing line with new content

**Diff Structure:**
```json
{
    "Type": 0,          // Operation type
    "LineIndex": 5,     // 1-based line number
    "OldContent": "",   // Previous content (for Delete/Modify)
    "NewContent": "new line content"  // New content (for Add/Modify)
}
```

### Multi-Line Operations

Complex edits can be performed with multiple diff operations in sequence:

```json
{
    "diffs": [
        {
            "Type": 1,
            "LineIndex": 10,
            "OldContent": "old function definition",
            "NewContent": ""
        },
        {
            "Type": 0,
            "LineIndex": 10,
            "OldContent": "",
            "NewContent": "def new_function():"
        },
        {
            "Type": 0,
            "LineIndex": 11,
            "OldContent": "",
            "NewContent": "    return True"
        }
    ]
}
```

## Safety and Approval Mechanisms

### Approval Workflow

The `AwaitingCommandApproval` flag determines execution behavior:

- `false`: **Automatic execution** for safe operations (reading files, listing directories)
- `true`: **User approval required** for potentially destructive operations (file modifications, system commands)

### Command Safety Classification

**Automatic Execution (AwaitingCommandApproval: false):**
- File reading operations
- Directory listing
- Git status checks
- Safe query commands

**Requires Approval (AwaitingCommandApproval: true):**
- File modifications
- Package installations
- System configuration changes
- Git commits and pushes
- Pull request creation

## GitHub Integration

### Pull Request Automation

The system prompt enables intelligent PR creation:

```
[FUNCTION_CALL]
Reason: Create a pull request for the current branch
AwaitingCommandApproval: true
Payload:
{
    "name": "createPullRequest",
    "args": {
        "branch_name": "feature/new-functionality",
        "title": "Add New Functionality",
        "summary": "This PR adds new functionality including:\n- Feature A\n- Bug fix for issue #123\n- Updated documentation"
    }
}
[/FUNCTION_CALL]
```

### Context-Aware Git Operations

The agent can perform intelligent git workflows:

1. **Staging Changes**: `git add .` with file awareness
2. **Commit Generation**: Descriptive commit messages based on changes
3. **Branch Management**: Smart branch naming and switching
4. **PR Creation**: Auto-generated titles and descriptions

## Advanced Prompt Techniques

### Context Preservation

The system maintains conversation context through:

- **Message History**: Persistent conversation chain
- **Working Directory Awareness**: Current location tracking
- **File State Tracking**: Modified files and their states

### Error Recovery

The prompt includes error handling instructions:

- **Command Failures**: Explanation and alternative approaches
- **Permission Issues**: Guidance for permission resolution
- **File Not Found**: Intelligent path suggestion and creation

### Iterative Command Execution

The agent follows a one-command-at-a-time pattern:

1. Suggest single command
2. Wait for execution and output
3. Analyze results
4. Suggest next command if needed

This prevents cascading failures and allows for user intervention.

## Prompt Version Management

### Current Implementation

The codebase includes two prompt versions:

1. **`getSystemPrompt()`**: Current production prompt with structured blocks
2. **`getSystemPromptNew()`**: Experimental JSON-based response format

### JSON Response Format (Experimental)

The new format uses structured JSON responses:

```json
{
    "type": "command_suggestion",
    "reason": "Explanation for the action",
    "awaiting_command_approval": false,
    "command": "ls -la",
    "function_call": null,
    "textual_response": null
}
```

**Response Types:**
- `command_suggestion`: Shell command execution
- `function_call`: File operations or API calls
- `textual_response`: Pure text response
- `github_action`: Git/GitHub operations

## Best Practices for Prompt Engineering

### 1. Specificity
- Use exact command formats and examples
- Specify platform constraints clearly
- Define precise response structures

### 2. Safety First
- Always require approval for destructive operations
- Provide clear reasoning for all actions
- Include rollback instructions where applicable

### 3. Context Awareness
- Inject dynamic context (CWD, shell, platform)
- Maintain conversation state
- Reference previous operations

### 4. Error Handling
- Anticipate common failure modes
- Provide alternative approaches
- Include diagnostic commands

### 5. User Experience
- Clear explanations for all actions
- Consistent response formatting
- Minimal cognitive load for approvals

## Implementation Details

### Prompt Parsing

The `parseCommandSuggestion()` function extracts structured responses:

```go
func parseCommandSuggestion(response string) *CommandSuggestion {
    // Parse [COMMAND_SUGGESTION] blocks
    // Extract Reason, Command, AwaitingCommandApproval
    // Return structured CommandSuggestion
}
```

### Function Call Parsing

Similar parsing for function calls:

```go
func parseFunctionCall(response string) *FunctionCall {
    // Parse [FUNCTION_CALL] blocks
    // Extract JSON payload
    // Validate function name and arguments
}
```

This prompt engineering approach enables Menace CLI to provide intelligent, safe, and context-aware terminal assistance while maintaining user control over potentially destructive operations.