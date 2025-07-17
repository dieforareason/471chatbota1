// Package db provides database functionality for the LLM bot
package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type Interaction struct {
	Prompt     string `json:"prompt"`
	Completion string `json:"completion"`
}

// ExportAsJSONL exports all interactions to a JSONL file
func (s *SQLiteStore) ExportAsJSONL(ctx context.Context, filename string) error {
	rows, err := s.db.QueryContext(ctx, "SELECT user_input, llm_response FROM interactions")
	if err != nil {
		return fmt.Errorf("querying interactions: %w", err)
	}
	defer rows.Close()

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	for rows.Next() {
		var entry Interaction
		if err := rows.Scan(&entry.Prompt, &entry.Completion); err != nil {
			return fmt.Errorf("scanning row: %w", err)
		}
		line, _ := json.Marshal(entry)
		fmt.Fprintln(file, string(line))
	}
	return rows.Err()
}
