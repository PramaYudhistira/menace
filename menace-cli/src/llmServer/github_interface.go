package llmServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type PullRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}

func hasChanges() (bool, string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, "", fmt.Errorf("failed to check git status: %v", err)
	}
	return len(output) > 0, string(output), nil
}

func CreatePullRequest(branchName string, title string, summary string) error {
	// Get GitHub token from environment
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	// Get repository info
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	repoURL, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get repository URL: %v", err)
	}

	// Parse repository URL to get owner and repo name
	// Example URL: https://github.com/owner/repo.git
	parts := strings.Split(strings.TrimSpace(string(repoURL)), "/")
	if len(parts) < 2 {
		return fmt.Errorf("invalid repository URL format")
	}
	repoName := strings.TrimSuffix(parts[len(parts)-1], ".git")
	owner := parts[len(parts)-2]

	// Create pull request
	pr := PullRequest{
		Title: title,
		Body:  summary,
		Head:  branchName,
		Base:  "main",
	}

	jsonData, err := json.Marshal(pr)
	if err != nil {
		return fmt.Errorf("failed to marshal PR data: %v", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repoName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusCreated {
		var errorBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorBody)
		return fmt.Errorf("failed to create PR: %s", errorBody["message"])
	}

	return nil
}

func PushToGitHub(commit_message string) error {
	//are we in a git repository?
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

	// Check if there are any changes
	hasChanges, _, err := hasChanges()
	if err != nil {
		return fmt.Errorf("failed to check for changes: %v", err)
	}

	if !hasChanges {
		fmt.Println("No changes to commit. Skipping push.")
		return nil
	} else {
		// Add all changes
		cmd = exec.Command("git", "add", ".")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add changes: %v", err)
		}

		// Commit changes
		fmt.Println("Committing changes...")
		cmd = exec.Command("git", "commit", "-m", commit_message)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to commit changes: %v", err)
		}
		fmt.Println("Changes added to commit")
	}

	// Push to GitHub
	fmt.Println("Pushing to GitHub...")
	cmd = exec.Command("git", "push", "origin", branchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push to GitHub: %v", err)
	}

	// Only create pull request if not on main branch
	if branchName != "main" {
		fmt.Println("Creating pull request...")
		if err := CreatePullRequest(branchName, "Auto PR by test program", "This is an automated pull request created by the test program."); err != nil {
			return fmt.Errorf("failed to create pull request: %v", err)
		}
	} else {
		fmt.Println("On main branch - skipping pull request creation")
	}

	return nil
}

type Add_check struct {
	Is_conversation_relevant_to_github string `json:"is_conversation_relevant_to_github"`
	Reason                             string `json:"reason"`
	Add                                string `json:"add"`
}

type Commit_check struct {
	Is_commit_needed string `json:"is_commit_needed"`
	Reason           string `json:"reason"`
	Commit_message   string `json:"commit_message"`
}

type Pull_request_check struct {
	Is_pull_request_needed string `json:"is_pull_request_needed"`
	Reason                  string `json:"reason"`
}

func isolated_single_message_to_ai(message string) (string, error) {
	llm, err := openai.New(
		openai.WithToken(os.Getenv("OPENAI_API_KEY")),
		openai.WithModel("o4-mini-2025-04-16"),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create OpenAI client: %v", err)
	}

	resp, err := llm.GenerateContent(context.Background(), []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{Text: message}},
		},
	}, llms.WithTemperature(1))

	if err != nil {
		return "", err
	}
	return resp.Choices[0].Content, nil
}

func convert_str_to_json(str string, json_format interface{}) error {
	err := json.Unmarshal([]byte(str), json_format)
	if err != nil {
		return err
	}
	return err
}
