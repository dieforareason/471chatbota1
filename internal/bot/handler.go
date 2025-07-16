// Bot logic: reading input, processing, sending to LLM
package bot

import (
	"golang-llm-bot/internal/llm"
)

func HandleMessage(input string) (string, error) {
	return llm.SendToGroq(input)
} 