// core/llm/client.go
package llm

import (
	"os"

	"github.com/go-resty/resty/v2"
)

const systemPrompt = "Kamu adalah seorang personal assistent yang baik, manja dan supportive. dan nama kamu adalah MELATI"

func SendToGroq(prompt string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	client := resty.New()

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"messages": []map[string]string{
				{
					"role":    "system",
					"content": systemPrompt,
				},
				{
					"role":    "user",
					"content": prompt,
				},
			},
			"model": "llama3-8b-8192",
		}).
		SetResult(&struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}{}).
		Post("https://api.groq.com/openai/v1/chat/completions")

	if err != nil {
		return "", err
	}

	result := resp.Result().(*struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	})

	if len(result.Choices) == 0 {
		return "No response", nil
	}

	return result.Choices[0].Message.Content, nil
}
