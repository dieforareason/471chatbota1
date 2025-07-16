// Entry point for the bot
package main

import (
	"bufio"
	"fmt"
	"golang-llm-bot/core/bot"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("GOLANG Bot is running. Type a message:")

	for scanner.Scan() {
		input := scanner.Text()
		response, err := bot.HandleMessage(input)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Println("Bot:", response)
		fmt.Println("\nType another message:")
	}
}
