# BASIC GOLANG LLM Bot

## Prerequisites

1. **Go Installation**
   - Go version 1.21 or higher is required
   - Download Go from [official website](https://go.dev/dl/)
   - Verify installation by running: `go version`

2. **Environment Setup**
   - A Groq API account and API key
   - Text editor or IDE of your choice
   - Git (optional, for version control)
   - SQLite database (can be managed with DBeaver or similar tool)

3. **System Requirements**
   - Any operating system that supports Go (Windows, macOS, or Linux)
   - Internet connection (required for API calls to Groq)

## Setup

1. Create a `.env` file with your Groq API key:
```
GROQ_API_KEY=your_groq_api_key
```

2. Install dependencies:
```bash
go mod download
```

3. Configure database path in `core/db/database.go` to point to your SQLite database file.

## Usage

Run the bot in chat mode:
```bash
go run cmd/bot/main.go
```

The bot will start and prompt you for input. Type your message and press Enter to get a response from MELATI.

Export conversation history to JSONL:
```bash
go run cmd/bot/main.go export
```
This will create a `training_data.jsonl` file containing all conversations in the format:
```jsonl
{"prompt":"user message","completion":"MELATI's response"}
```

## Project Structure

- `cmd/bot/main.go`: Entry point for the bot
- `core/llm/client.go`: Client code for interacting with Groq LLM API
- `core/bot/handler.go`: Bot logic for processing input and sending to LLM
- `core/bot/prompts.go`: System prompt configuration
- `core/db/database.go`: SQLite database initialization
- `core/db/logger.go`: Conversation logging functionality
- `core/db/export.go`: JSONL export functionality
