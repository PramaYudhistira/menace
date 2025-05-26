package main

import (
	"github.com/tmc/langchaingo/llms"
	"menace-go/github"
	"fmt"
)

func main() {
	// TEST 1: make a python file and user asks for push
	fmt.Println("-------- TEST 1: make a python file and user asks for push")
	example_messages := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeGeneric,
			Parts: []llms.ContentPart{llms.TextContent{
			Text: "I want to make a python file that multiplies two numbers"}},
		},
		{
			Role:  llms.ChatMessageTypeGeneric,
			Parts: []llms.ContentPart{llms.TextContent{
			Text: "The file is made and is viewable"}},
		},
		{
			Role:  llms.ChatMessageTypeGeneric,
			Parts: []llms.ContentPart{llms.TextContent{
			Text: "Thank you, now lets go ahead and stage"}},
		},
	}
	github.GithubStart(example_messages)

	// TEST 2: user wants to stage only
	fmt.Println("-------- TEST 2: user wants to stage only")
	example_messages2 := []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeGeneric,
			Parts: []llms.ContentPart{llms.TextContent{
				Text: "Please stage the changes"}},
		},
		{
			Role: llms.ChatMessageTypeGeneric,
			Parts: []llms.ContentPart{llms.TextContent{
				Text: "Just prepare it for commit"}},
		},
	}
	github.GithubStart(example_messages2)

	// TEST 3: user wants to stage and push
	fmt.Println("-------- TEST 3: user wants to stage and push")
	example_messages3 := []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeGeneric,
			Parts: []llms.ContentPart{llms.TextContent{
				Text: "Go ahead and stage and push the file"}},
		},
		{
			Role: llms.ChatMessageTypeGeneric,
			Parts: []llms.ContentPart{llms.TextContent{
				Text: "Save and upload this to the repository"}},
		},
	}
	github.GithubStart(example_messages3)

	// TEST 4: user wants to stage, push and make a PR
	fmt.Println("-------- TEST 4: user wants to stage, push and make a PR")
	example_messages4 := []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeGeneric,
			Parts: []llms.ContentPart{llms.TextContent{
				Text: "Please stage these updates, push them, and open a pull request"}},
		},
		{
			Role: llms.ChatMessageTypeGeneric,
			Parts: []llms.ContentPart{llms.TextContent{
				Text: "Let's create a PR with the latest changes"}},
		},
	}
	github.GithubStart(example_messages4)

	// TEST 5: user just discussing, no action needed
	fmt.Println("-------- TEST 5: user just discussing, no action needed")
	example_messages5 := []llms.MessageContent{
	{
		Role: llms.ChatMessageTypeGeneric,
		Parts: []llms.ContentPart{llms.TextContent{
			Text: "Can you explain what the code does?"}},
	},
	{
		Role: llms.ChatMessageTypeGeneric,
		Parts: []llms.ContentPart{llms.TextContent{
			Text: "Interesting, I’ll think about modifying it"}},
	},
	}
	github.GithubStart(example_messages5)

}