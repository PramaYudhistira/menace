package llmServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Message represents a message in the LLM conversation

type PullRequestReview struct {
	Body  string `json:"body"`
	Event string `json:"event"` // APPROVE, REQUEST_CHANGES, or COMMENT
}

// ReviewPullRequest adds a review comment to a pull request
func ReviewPullRequest(owner, repoName string, prNumber int, comment string, event string) error {
	// Get GitHub token from environment
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	// Create review
	review := PullRequestReview{
		Body:  comment,
		Event: event,
	}

	jsonData, err := json.Marshal(review)
	if err != nil {
		return fmt.Errorf("failed to marshal review data: %v", err)
	}
	fmt.Println("jsonData: ", jsonData)
	// Submit review
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/reviews", owner, repoName, prNumber)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create review request: %v", err)
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to submit review: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorBody)
		return fmt.Errorf("failed to submit review: %s", errorBody["message"])
	}

	fmt.Println("comment: ", comment)
	return nil
}

func main() {
	// Test parameters
	owner := "Serhan-Asad"
	repoName := "Beginner"
	prNumber := 1
	comment := "This is a test review comment. The code looks good!"
	event := "COMMENT" // Can be "APPROVE", "REQUEST_CHANGES", or "COMMENT"

	// Call the review function
	err := ReviewPullRequest(owner, repoName, prNumber, comment, event)
	if err != nil {
		fmt.Printf("Error reviewing PR: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("PR review completed successfully!")
}
