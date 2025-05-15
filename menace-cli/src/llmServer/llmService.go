package llmServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Available models
const (
	GPT4Mini     = "gpt-4-0125-preview" // Latest GPT-4 model
	GPT4MiniHigh = "gpt-4-1106-preview" // Higher quality GPT-4 model
	GPT35Turbo   = "gpt-3.5-turbo"      // Fallback model
)

var (
	shell         = ModelFactory{}.DetectShell()
	shellType     = strings.Split(shell, "/")[1] // Get just the shell type (bash, powershell, etc.)
	system_prompt = fmt.Sprintf(`You are operating as and within the Menace CLI built by Prama Yudhistira. You must be safe, precise and helpful.
	Menace-CLI is a Go-based CLI tool that uses large language models to provide intelligent terminal assistance.
	You have access to the local file system and can execute commands in %s.

	When you need to execute a command, format it like this:
	`+"```"+shellType+"\n"+`your_command_here
	`+"```\n"+`

	For example, to list files:
	`+"```"+shellType+"\n"+`ls
	`+"```\n"+`

	You should respond as if you are part of this real application, not a fictional tool.
	Always explain what you're doing before executing commands, and explain the results after.`, shell)
)

// LLMService represents the singleton LLM service
type LLMService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	model      string
	messages   []Message //this stores context for now, we will eventually run out of stack space
	mu         sync.Mutex
}

var (
	instance *LLMService
	once     sync.Once
)

// GetInstance returns the singleton instance of LLMService
/*
TODO: messages field will eventually become large and will cause stack overflow
      make sure to implement directory specific context
*/
func GetInstance() *LLMService {
	once.Do(func() {
		instance = &LLMService{
			httpClient: &http.Client{},
			baseURL:    "https://api.openai.com/v1/chat/completions",
			model:      GPT4MiniHigh, // Default to the higher quality GPT-4 model
			messages: []Message{
				{
					Role:    "system",
					Content: system_prompt,
				},
			},
		}
	})
	return instance
}

// Configure sets up the LLM service with the provided API key
func (s *LLMService) Configure(apiKey string) {
	s.apiKey = apiKey
}

// SetModel allows changing the model
func (s *LLMService) SetModel(model string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.model = model
}

// GetModel returns the current model
func (s *LLMService) GetModel() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.model
}

// ClearHistory clears the conversation history
func (s *LLMService) ClearHistory() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = make([]Message, 0)
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the request structure for the LLM, the entire context is stored in Messages field
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse represents the response from the LLM
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// SendMessage sends a message to the LLM and returns the response
func (s *LLMService) SendMessage(input string) (*ChatResponse, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("API key not configured")
	}

	s.mu.Lock()
	// Add user message to history
	s.messages = append(s.messages, Message{
		Role:    "user",
		Content: input,
	})

	reqBody := ChatRequest{
		Model:    s.model,
		Messages: s.messages, //entire history sent
	}
	s.mu.Unlock()

	//convert ChatRequest to json
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Add assistant's response to history
	if len(chatResp.Choices) > 0 {
		s.mu.Lock()
		s.messages = append(s.messages, chatResp.Choices[0].Message)
		s.mu.Unlock()
	}

	return &chatResp, nil
}
