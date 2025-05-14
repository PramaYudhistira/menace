package llmServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Available models
const (
	GPT4Mini     = "gpt-4-0125-preview" // Latest GPT-4 model
	GPT4MiniHigh = "gpt-4-1106-preview" // Higher quality GPT-4 model
	GPT35Turbo   = "gpt-3.5-turbo"      // Fallback model
)

// LLMService represents the singleton LLM service
type LLMService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	model      string
	messages   []Message //this stores context for now...
	mu         sync.Mutex
}

var (
	instance *LLMService
	once     sync.Once
)

// GetInstance returns the singleton instance of LLMService
func GetInstance() *LLMService {
	once.Do(func() {
		instance = &LLMService{
			httpClient: &http.Client{},
			baseURL:    "https://api.openai.com/v1/chat/completions",
			model:      GPT4MiniHigh, // Default to the higher quality GPT-4 model
			messages:   make([]Message, 0),
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

// ChatRequest represents the request structure for chat completion
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
		Messages: s.messages,
	}
	s.mu.Unlock()

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
