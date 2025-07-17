// Package llm provides functionality for interacting with the Groq LLM API
package llm

import (
	"context"
	"fmt"
	"time"

	"golang-llm-sqlite-bot/core/config"

	"github.com/go-resty/resty/v2"
)

// Client defines the interface for LLM interactions
type Client interface {
	SendMessage(ctx context.Context, prompt string) (string, error)
}

// GroqClient implements the LLM Client interface for Groq's API
type GroqClient struct {
	client *resty.Client
	config *config.Config
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents the API response structure
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewGroqClient creates a new Groq API client with retry middleware
func NewGroqClient(cfg *config.Config) *GroqClient {
	client := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(1*time.Second).
		SetRetryMaxWaitTime(5*time.Second).
		SetTimeout(cfg.RequestTimeout).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", cfg.GroqAPIKey)).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return err != nil || r.StatusCode() >= 500
		})

	return &GroqClient{
		client: client,
		config: cfg,
	}
}

// SendMessage sends a message to the Groq API and returns the response
func (c *GroqClient) SendMessage(ctx context.Context, prompt string) (string, error) {
	messages := []Message{
		{Role: "system", Content: c.config.SystemPrompt},
		{Role: "user", Content: prompt},
	}

	var result ChatResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"messages": messages,
			"model":    c.config.ModelName,
		}).
		SetResult(&result).
		Post("https://api.groq.com/openai/v1/chat/completions")

	if err != nil {
		return "", fmt.Errorf("failed to send message to Groq: %w", err)
	}

	if !resp.IsSuccess() {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned from API")
	}

	return result.Choices[0].Message.Content, nil
}
