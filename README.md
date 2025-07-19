# LLM Chatbot with WhatsApp Integration

A Golang-based chatbot that uses Groq LLM API for processing messages and can interact through both CLI and WhatsApp interfaces.

## Features

- CLI-based chat interface
- WhatsApp integration via webhook
- Groq LLM API integration
- SQLite message history storage
- Multi-interface support (CLI and WhatsApp)

## Prerequisites

- Go 1.19 or higher
- SQLite
- Groq API Key
- WhatsApp API Server (go-whatsapp-web-multidevice) ( https://github.com/aldinokemal/go-whatsapp-web-multidevice )

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd 471chatbota1
```

2. Copy the environment file and configure it:
```bash
vi .env
```

3. Set your Groq API key in the .env file:
```
GROQ_API_KEY=your-api-key-here
```

## Usage

### CLI Mode

Run the CLI bot:
```bash
go run cmd/bot/main.go
```

### WhatsApp Mode

1. First, set up and run the WhatsApp API server (go-whatsapp-web-multidevice):
```bash
# In WSL or your preferred environment
./whatsapp-api --webhook="http://YOUR_WINDOWS_IP:4444/webhook/wa"
```

2. Run the webhook server:
```bash
go run cmd/wabot/main.go
```

3. Send a message to your connected WhatsApp number to interact with the bot.

### Configuration

The bot can be configured through environment variables:

- `GROQ_API_KEY`: Your Groq API key
- `DEFAULT_PROMPT`: Custom system prompt for the LLM (optional)
- `MODEL_NAME`: LLM model name
- `REQUEST_TIMEOUT`: API request timeout
- `DB_PATH`: SQLite database path

#### Prompt Customization

The bot uses a default system prompt defined in `core/config/prompts.go`. You can customize the prompt in two ways:

1. Set the `DEFAULT_PROMPT` environment variable:
```bash
DEFAULT_PROMPT="You are a helpful assistant that..."
```

2. Modify the `DefaultPrompt` constant in `core/config/prompts.go` for a permanent change.

The default prompt includes multiple traits to define the assistant's behavior, making it easy to adjust the bot's personality and communication style.

## Architecture

The project follows a clean architecture pattern:

```
.
├── cmd/
│   ├── bot/      # CLI bot entry point
│   └── wabot/    # WhatsApp bot entry point
├── core/
│   ├── bot/      # Bot logic and handlers
│   ├── config/   # Configuration management
│   ├── db/       # Database operations
│   └── llm/      # LLM client implementation
```

### WhatsApp Integration

The bot integrates with WhatsApp through a webhook server that:
1. Receives messages from WhatsApp API
2. Processes them using Groq LLM
3. Sends responses back to WhatsApp

Webhook payload format:
```json
{
    "chat_id": "1234567890",
    "from": "1234567890@s.whatsapp.net",
    "message": {
        "text": "Hello bot!",
        "id": "message-id",
        "replied_id": "",
        "quoted_message": ""
    },
    "pushname": "User",
    "sender_id": "1234567890",
    "timestamp": "2025-07-17T10:03:34Z"
}
```

## Development

### Running Tests

```bash
go test ./...
```

### Building

Build both CLI and WhatsApp bots:
```bash
go build -o bin/bot cmd/bot/main.go
go build -o bin/wabot cmd/wabot/main.go
```
