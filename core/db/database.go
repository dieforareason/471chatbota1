// Package db provides database functionality for the LLM bot
package db

import (
	"context"
	"database/sql"
	"fmt"

	"golang-llm-sqlite-bot/core/config"

	_ "modernc.org/sqlite"
)

// Store defines the interface for database operations
type Store interface {
	LogInteraction(ctx context.Context, prompt, response string) error
	Close() error
}

// SQLiteStore implements the Store interface
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLite database connection with proper connection pooling
func NewSQLiteStore(cfg *config.Config) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	// Configure connection pooling
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	store := &SQLiteStore{db: db}
	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

// initSchema ensures the required database schema exists
func (s *SQLiteStore) initSchema() error {
	createTable := `
	CREATE TABLE IF NOT EXISTS interactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_input TEXT NOT NULL,
		llm_response TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := s.db.Exec(createTable)
	if err != nil {
		return fmt.Errorf("creating interactions table: %w", err)
	}
	return nil
}

// LogInteraction stores a user interaction in the database
func (s *SQLiteStore) LogInteraction(ctx context.Context, prompt, response string) error {
	const query = `INSERT INTO interactions (user_input, llm_response) VALUES (?, ?)`

	result, err := s.db.ExecContext(ctx, query, prompt, response)
	if err != nil {
		return fmt.Errorf("inserting interaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", rows)
	}

	return nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("closing database: %w", err)
	}
	return nil
}
