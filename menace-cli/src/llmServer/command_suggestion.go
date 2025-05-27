package llmServer

import (
	"strings"
)

// CommandSuggestion is a struct that contains the reason and command from the LLM response
// This is used to parse the LLM response and extract the command suggestion
type CommandSuggestion struct {
	Reason       string
	Command      string
	Human_needed string
}

// Run this after every LLM response
//
// Returns a CommandSuggestion if a command is found.
// Returns nil otherwise.
func parseCommandSuggestion(response string) *CommandSuggestion {
	// Find the command suggestion block
	start := strings.Index(response, "[COMMAND_SUGGESTION]")
	end := strings.Index(response, "[/COMMAND_SUGGESTION]")

	if start == -1 || end == -1 {
		return nil // No command suggestion found
	}

	// Extract the content between the tags
	content := response[start:end]

	// Parse reason and command
	reasonStart := strings.Index(content, "Reason: ")
	commandStart := strings.Index(content, "Command: ")
	humanNeededStart := strings.Index(content, "Human_needed: ")

	if reasonStart == -1 || commandStart == -1 || humanNeededStart == -1 {
		return nil
	}

	// Extract reason (from "Reason: " to "Command: ")
	reason := strings.TrimSpace(content[reasonStart+8 : commandStart])

	// Extract command (from "Command: " to end)
	command := strings.TrimSpace(content[commandStart+8 : humanNeededStart])

	// Extract human needed (from "Human_needed: " to end)
	humanNeeded := strings.TrimSpace(content[humanNeededStart+12:])

	return &CommandSuggestion{
		Reason:  reason,
		Command: command,
		Human_needed: humanNeeded,
	}
}
