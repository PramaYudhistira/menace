package llmServer

import (
	"fmt"
	"os/exec"
	"strings"
)

// PushToGitHub handles pushing changes to GitHub when prompted by the LLM
func (a *Agent) PushToGitHub() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if we're in a git repository
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not in a git repository: %v", err)
	}

	// Get the current branch
	cmd = exec.Command("git", "branch", "--show-current")
	branch, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}
	branchName := strings.TrimSpace(string(branch))

	// Add all changes
	cmd = exec.Command("git", "add", ".")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add changes: %v", err)
	}

	// Commit changes
	cmd = exec.Command("git", "commit", "-m", "Auto-commit by Menace CLI")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit changes: %v", err)
	}

	// Push to GitHub
	cmd = exec.Command("git", "push", "origin", branchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push to GitHub: %v", err)
	}

	return nil
}