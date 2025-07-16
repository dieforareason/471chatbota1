// core/bot/handler.go
package bot

import (
	"golang-llm-sqlite-bot/core/db"
	"golang-llm-sqlite-bot/core/llm"
)

func HandleMessage(input string) (string, error) {
	resp, err := llm.SendToGroq(input)
	if err != nil {
		return "", err
	}
	db.LogInteraction(input, resp)
	return resp, nil
}
