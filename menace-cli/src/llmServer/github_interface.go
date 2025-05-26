package llmServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type PullRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}

func hasChanges() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %v", err)
	}
	return len(output) > 0, nil
}

func createPullRequest(branchName string) (string, error) {
	var logs strings.Builder
	// Get GitHub token from environment
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		logs.WriteString("GITHUB_TOKEN environment variable is not set\n")
		return logs.String(), fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	// Get repository info
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	repoURL, err := cmd.Output()
	if err != nil {
		logs.WriteString(fmt.Sprintf("failed to get repository URL: %v\n", err))
		return logs.String(), fmt.Errorf("failed to get repository URL: %v", err)
	}

	// Parse repository URL to get owner and repo name
	// Example URL: https://github.com/owner/repo.git
	parts := strings.Split(strings.TrimSpace(string(repoURL)), "/")
	if len(parts) < 2 {
		logs.WriteString("invalid repository URL format\n")
		return logs.String(), fmt.Errorf("invalid repository URL format")
	}
	repoName := strings.TrimSuffix(parts[len(parts)-1], ".git")
	owner := parts[len(parts)-2]

	// Create pull request
	pr := PullRequest{
		Title: "Auto PR by test program",
		Body:  "This is an automated pull request created by the test program.",
		Head:  branchName,
		Base:  "main",
	}

	jsonData, err := json.Marshal(pr)
	if err != nil {
		logs.WriteString(fmt.Sprintf("failed to marshal PR data: %v\n", err))
		return logs.String(), fmt.Errorf("failed to marshal PR data: %v", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repoName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logs.WriteString(fmt.Sprintf("failed to create request: %v\n", err))
		return logs.String(), fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.WriteString(fmt.Sprintf("failed to send request: %v\n", err))
		return logs.String(), fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusCreated {
		var errorBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorBody)
		logs.WriteString(fmt.Sprintf("failed to create PR: %s\n", errorBody["message"]))
		return logs.String(), fmt.Errorf("failed to create PR: %s", errorBody["message"])
	}

	logs.WriteString("Pull request created successfully!\n")
	return logs.String(), nil
}

func PushToGitHub() (string, error) {
	var logs strings.Builder
	//are we in a git repository?
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		logs.WriteString(fmt.Sprintf("not in a git repository: %v\n", err))
		return logs.String(), fmt.Errorf("not in a git repository: %v", err)
	}

	// Get the current branch
	cmd = exec.Command("git", "branch", "--show-current")
	branch, err := cmd.Output()
	if err != nil {
		logs.WriteString(fmt.Sprintf("failed to get current branch: %v\n", err))
		return logs.String(), fmt.Errorf("failed to get current branch: %v", err)
	}
	branchName := strings.TrimSpace(string(branch))
	logs.WriteString(fmt.Sprintf("Current branch: %s\n", branchName))

	// Check if there are any changes
	hasChanges, err := hasChanges()
	if err != nil {
		logs.WriteString(fmt.Sprintf("failed to check for changes: %v\n", err))
		return logs.String(), fmt.Errorf("failed to check for changes: %v", err)
	}

	if !hasChanges {
		logs.WriteString("No changes to commit. Proceeding with push...\n")
	} else {
		// Add all changes
		cmd = exec.Command("git", "add", ".")
		if err := cmd.Run(); err != nil {
			logs.WriteString(fmt.Sprintf("failed to add changes: %v\n", err))
			return logs.String(), fmt.Errorf("failed to add changes: %v", err)
		}

		// Commit changes
		logs.WriteString("Committing changes...\n")
		cmd = exec.Command("git", "commit", "-m", "Auto-commit by test program")
		if err := cmd.Run(); err != nil {
			logs.WriteString(fmt.Sprintf("failed to commit changes: %v\n", err))
			return logs.String(), fmt.Errorf("failed to commit changes: %v", err)
		}
		logs.WriteString("Changes added to commit\n")
	}

	// Push to GitHub
	logs.WriteString("Pushing to GitHub...\n")
	cmd = exec.Command("git", "push", "origin", branchName)
	if err := cmd.Run(); err != nil {
		logs.WriteString(fmt.Sprintf("failed to push to GitHub: %v\n", err))
		return logs.String(), fmt.Errorf("failed to push to GitHub: %v", err)
	}

	// Only create pull request if not on main branch
	if branchName != "main" {
		logs.WriteString("Creating pull request...\n")
		prLogs, err := createPullRequest(branchName)
		logs.WriteString(prLogs)
		if err != nil {
			return logs.String(), fmt.Errorf("failed to create pull request: %v", err)
		}
	} else {
		logs.WriteString("On main branch - skipping pull request creation\n")
	}

	return logs.String(), nil
}

// HandlePushCommand handles push command suggestions
func HandlePushCommand(cmd *CommandSuggestion) (string, error) {
	// Check if this is a push command
	if !strings.Contains(cmd.Command, "git push") {
		return "", nil // Not a push command, ignore
	}
	return PushToGitHub()
}
