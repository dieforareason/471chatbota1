// Package bot provides the main bot functionality
package bot

import (
	"context"
	"fmt"

	"golang-llm-sqlite-bot/core/db"
	"golang-llm-sqlite-bot/core/llm"
)

// Bot handles chat interactions using the LLM service
type Bot struct {
	llm   llm.Client
	store db.Store
}

// NewBot creates a new bot instance with the provided dependencies
func NewBot(llmClient llm.Client, store db.Store) *Bot {
	return &Bot{
		llm:   llmClient,
		store: store,
	}
}

// HandleMessage processes a user message and returns the LLM's response
func (b *Bot) HandleMessage(ctx context.Context, input string) (string, error) {
	// Send message to LLM
	resp, err := b.llm.SendMessage(ctx, input)
	if err != nil {
		return "", fmt.Errorf("getting LLM response: %w", err)
	}

	// Log the interaction
	if err := b.store.LogInteraction(ctx, input, resp); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to log interaction: %v\n", err)
	}

	return resp, nil
}
