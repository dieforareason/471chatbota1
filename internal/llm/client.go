// Client code for interacting with Groq LLM API
package llm

import (
	"encoding/json"
	"os"

	"github.com/go-resty/resty/v2"
)

func SendToGroq(systemPrompt string, userInput string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"messages": []map[string]string{
				{"role": "system", "content": systemPrompt},
				{"role": "user", "content": userInput},
			},
			"model": "llama3-8b-8192",
		}).
		Post("https://api.groq.com/openai/v1/chat/completions")

	if err != nil {
		return "", err
	}

	// Parse JSON response
	type choice struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	var result struct {
		Choices []choice `json:"choices"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "No response", nil
	}
	return result.Choices[0].Message.Content, nil
}
