// cmd/bot/main.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang-llm-sqlite-bot/core/bot"
	"golang-llm-sqlite-bot/core/db"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db.InitDB()

	if len(os.Args) > 1 && os.Args[1] == "export" {
		err := db.ExportAsJSONL("training_data.jsonl")
		if err != nil {
			log.Fatal("Export failed:", err)
		}
		fmt.Println("✅ Exported to training_data.jsonl")
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("🤖 Bot is ready. Type a message:")
	for scanner.Scan() {
		input := scanner.Text()
		if input == "exit" {
			fmt.Println("👋 Goodbye!")
			break
		}
		resp, err := bot.HandleMessage(input)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Println("Bot:", resp)
		fmt.Println("\nType another message or 'exit' to quit:")
	}
}
