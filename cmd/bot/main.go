// Package main provides the entry point for the LLM bot
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"golang-llm-sqlite-bot/core/bot"
	"golang-llm-sqlite-bot/core/config"
	"golang-llm-sqlite-bot/core/db"
	"golang-llm-sqlite-bot/core/llm"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Validate required config
	if cfg.GroqAPIKey == "" {
		log.Fatal("GROQ_API_KEY environment variable is required")
	}

	// Create LLM client
	llmClient := llm.NewGroqClient(cfg)

	// Initialize database
	store, err := db.NewSQLiteStore(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer store.Close()

	// Create bot instance
	chatBot := bot.NewBot(llmClient, store)

	// Set up context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		fmt.Println("\nGoodbye!")
		cancel()
	}()

	// Start chat loop
	fmt.Println("Chat started. Press Ctrl+C to exit.")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Print("> ")
			if !scanner.Scan() {
				return
			}

			input := strings.TrimSpace(scanner.Text())
			if input == "" {
				continue
			}

			// Create context with timeout for each request
			reqCtx, reqCancel := context.WithTimeout(ctx, 30*time.Second)
			response, err := chatBot.HandleMessage(reqCtx, input)
			reqCancel()

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			fmt.Printf("Bot: %s\n\n", response)
		}
	}
}
