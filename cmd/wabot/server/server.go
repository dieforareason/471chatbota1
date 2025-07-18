// Package server provides the WhatsApp bot server entry point
package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang-llm-sqlite-bot/core/bot"
	"golang-llm-sqlite-bot/core/config"
	"golang-llm-sqlite-bot/core/db"
	"golang-llm-sqlite-bot/core/llm"

	"github.com/joho/godotenv"
)

func StartServer() {
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

	// Set up HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook/wa", chatBot.HandleWhatsAppWebhook)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for errors coming from the server
	serverErrors := make(chan error, 1)

	// Start server
	go func() {
		log.Printf("WhatsApp bot server starting on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	// Blocking main and waiting for shutdown
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)

	case <-shutdown:
		log.Println("Starting shutdown...")

		// Give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Asking listener to shut down
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("Graceful shutdown did not complete in %v: %v", 10*time.Second, err)
			err = server.Close()
		}

		if err != nil {
			log.Fatalf("Could not stop server gracefully: %v", err)
		}
	}
}
