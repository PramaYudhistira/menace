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
}

func main(messages []llms.MessageContent) {
	// 1. check if the message chain needs to be added/staged
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
	`, messages)

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
			adds := exec.Command("git", "status", ".")
			fmt.Println(adds)
			commit_prompt := fmt.Sprintf(`
				Read the following file changes and determine if we should make a commit.

				File changes:
				%s

				RETURN JSON IN THE FOLLOWING FORMAT:
				{
					"reason": "reason for the decision"
					"is_commit_needed": "true/false"
				}
			`, adds)

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
					fmt.Println("Code committed")
				}
			} else {
				fmt.Println("Code does not need to be committed")
			}
		} else {
			fmt.Println("Code does not need to be staged")
		}
	}

	// 3. assess if enough pushes have been made, then create a pull request

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



const example_messages = []llms.MessageContent{
	{
		Role: llms.ChatMessageTypeSystem,
		Parts: []llms.ContentPart{llms.TextContent{Text: "You are a helpful assistant."}},
	},
	{
		Role: llms.ChatMessageTypeUser,
		Parts: []llms.ContentPart{llms.TextContent{Text: "What is the capital of France?"}},
	},
	{
		Role: llms.ChatMessageTypeAssistant,
		Parts: []llms.ContentPart{llms.TextContent{Text: "The capital of France is Paris."}},
	},
}

