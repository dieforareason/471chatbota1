// Package main provides the WhatsApp bot entry point
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start webhook server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := chatBot.StartWebhookServer(8080); err != nil {
			errChan <- err
		}
	}()

	// Wait for signal or error
	select {
	case <-sigChan:
		log.Println("Shutting down gracefully...")
	case err := <-errChan:
		log.Printf("Server error: %v", err)
	}
}
