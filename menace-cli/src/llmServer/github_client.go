package llmServer

import (
	"fmt"
	"os/exec"

	"github.com/tmc/langchaingo/llms"
)

func GithubStart(messages []llms.MessageContent, a *Agent) string {
	// 1. check if the message chain needs to be added/staged
	recent_messages := messages
	if len(messages) > 3 {
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
		return "error"
	} else {
		// fmt.Println(resp)
		var add_checker Add_check
		err = convert_str_to_json(resp, &add_checker)
		if err != nil {
			return err.Error()
		}
		fmt.Println(add_checker)
		if add_checker.Add == "true" {
			a.AddToMessageChain(
				"I think we should stage the changes",
				llms.ChatMessageTypeSystem,
			)
			// 2. if staged, assess if enough changes have been made, then commit and push
			hasChanges, adds, err := hasChanges()
			if err != nil {
				return err.Error()
			}
			if hasChanges {
				a.messages = append(a.messages, llms.MessageContent{
					Role:  llms.ChatMessageTypeSystem,
					Parts: []llms.ContentPart{llms.TextContent{Text: "It seems like there are changes that need to be committed, let's check them out."}},
				})
			}
			// -------------------- GIT COMMIT --------------------
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
				return err.Error()
			}
			// fmt.Println(commit_resp)
			var commit_checker Commit_check
			err = convert_str_to_json(commit_resp, &commit_checker)
			if err != nil {
				return err.Error()
			}
			if commit_checker.Is_commit_needed == "true" {
				a.AddToMessageChain(
					"Do you want to commit the code as is? (y/n): ",
					llms.ChatMessageTypeSystem,
				)
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
						return err.Error()
					}
					// fmt.Println(pull_request_resp)
					var pull_request_checker Pull_request_check
					err = convert_str_to_json(pull_request_resp, &pull_request_checker)
					if err != nil {
						return err.Error()
					}
					if pull_request_checker.Is_pull_request_needed == "true" {
						a.AddToMessageChain(
							"Do you want to create a pull request? (y/n): ",
							llms.ChatMessageTypeAI,
						)
						var response string
						fmt.Scanln(&response)
						if response == "y" {
							branch := exec.Command("git", "branch", "--show-current")
							branch_name, err := branch.Output()
							if err != nil {
								return err.Error()
							}
							createPullRequest(string(branch_name))
						}
					}
				}
			}
		}
	}
	return "continue with the conversation"
}
