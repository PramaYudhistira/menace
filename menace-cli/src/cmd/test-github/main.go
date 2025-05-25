package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func pushToGitHub() error {
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
	fmt.Printf("Current branch: %s\n", branchName)

	// Add all changes
	fmt.Println("Adding changes...")
	cmd = exec.Command("git", "add", ".")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add changes: %v", err)
	}

	// Commit changes
	fmt.Println("Committing changes...")
	cmd = exec.Command("git", "commit", "-m", "Auto-commit by test program")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit changes: %v", err)
	}

	// Push to GitHub
	fmt.Println("Pushing to GitHub...")
	cmd = exec.Command("git", "push", "origin", branchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push to GitHub: %v", err)
	}

	return nil
}

func main() {
	fmt.Println("Attempting to push to GitHub...")
	err := pushToGitHub()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully pushed to GitHub!")
}
