# Golang LLM Bot

A simple Go bot that reads user input, sends it to Groq LLM API, and prints the response.

## Setup

1. Create a `.env` file with your Groq API key:
```
GROQ_API_KEY=your_groq_api_key
```

2. Initialize the Go module:
```bash
go mod init golang-llm-bot
```

3. Install dependencies:
```bash
go get
```

## Usage

Run the bot:
```bash
go run cmd/bot/main.go
```

The bot will start and prompt you for input. Type your message and press Enter to get a response from the LLM.

## Project Structure

- `cmd/bot/main.go`: Entry point for the bot
- `internal/llm/client.go`: Client code for interacting with Groq LLM API
- `internal/bot/handler.go`: Bot logic for processing input and sending to LLM