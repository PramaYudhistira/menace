package github

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

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

func GithubStart(messages []llms.MessageContent) {
	// 1. check if the message chain needs to be added/staged
	recent_messages := messages
	if (len(messages) > 3) {
		recent_messages = messages[len(messages)-3:]
	}

	// -------------------- GIT ADD --------------------
	add_prompt := fmt.Sprintf(`
	Read the following messages between a user and an assistant, and determine whether 
	the code needs to be staged in version control. Does the user want code to be staged?

	Messages:
	%s

	RETURN JSON IN THE FOLLOWING FORMAT:
	{
	    "is_conversation_relevant_to_github": "true/false",
		"reason": "reason for the decision"
		"add": "true/false"
	}
	`, recent_messages)

	resp, err := isolated_single_message_to_ai(add_prompt)

	if err != nil {
		return
	} else {
		fmt.Println(resp)
		var add_checker Add_check
		err = convert_str_to_json(resp, &add_checker)
		if err != nil {
			return
		}
		if add_checker.Add == "true" {
			fmt.Println("Code needs to be staged")
			// 2. if staged, assess if enough changes have been made, then commit and push
			hasChanges, adds, err := hasChanges()
			if err != nil {
				return
			}
			if hasChanges {
				fmt.Println("It seems like there are changes that need to be committed, let's check them out.")
			}
			// -------------------- GIT COMMIT --------------------
			fmt.Println(adds)
			commit_prompt := fmt.Sprintf(`
				Read the following file changes and determine if we should make a commit. Does the user want to commit the code as is?

				File changes:
				%s

				Last 3 messages:
				%s

				RETURN JSON IN THE FOLLOWING FORMAT:
				{
					"reason": "reason for the decision"
					"is_commit_needed": "true/false"
					"commit_message": "message for the commit" | nil
				}
			`, adds, recent_messages)

			commit_resp, err := isolated_single_message_to_ai(commit_prompt)
			if err != nil {
				return
			}
			fmt.Println(commit_resp)
			var commit_checker Commit_check
			err = convert_str_to_json(commit_resp, &commit_checker)
			if err != nil {
				return
			}
			if commit_checker.Is_commit_needed == "true" {
				fmt.Print("Do you want to commit the code as is? (y/n): ")
				var response string
				fmt.Scanln(&response)
				if response == "y" {
					PushToGitHub(commit_checker.Commit_message)
					// -------------------- PULL REQUEST --------------------
					pushes := "blah blah blah"
					pull_request_prompt := fmt.Sprintf(`
						Read the following file changes and determine if we should make a pull request.

						Recent pushes into this branch:
						%s

						Latest messages:
						%s

						RETURN JSON IN THE FOLLOWING FORMAT:
						{
							"reason": "reason for the decision"
							"is_pull_request_needed": "true/false"
						}
					`, pushes, recent_messages)

					pull_request_resp, err := isolated_single_message_to_ai(pull_request_prompt)
					if err != nil {
						return
					}
					fmt.Println(pull_request_resp)
					var pull_request_checker Pull_request_check
					err = convert_str_to_json(pull_request_resp, &pull_request_checker)
					if err != nil {
						return
					}
					if pull_request_checker.Is_pull_request_needed == "true" {
						fmt.Print("Do you want to create a pull request? (y/n): ")
						var response string
						fmt.Scanln(&response)
						if response == "y" {
							branch := exec.Command("git", "branch", "--show-current")
							branch_name, err := branch.Output()
							if err != nil {
								return
							}
							createPullRequest(string(branch_name))
						}
					}
				}
			} else {
				fmt.Println("Code does not need to be committed")
			}
		} else {
			fmt.Println("Code does not need to be staged")
		}
	}
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