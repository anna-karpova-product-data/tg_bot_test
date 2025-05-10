# Telegram Bot

A simple Telegram bot written in Go.

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (optional, for containerized deployment)
- Telegram Bot Token (get it from [@BotFather](https://t.me/BotFather))

## Configuration

1. Copy `.env.dev` to `.env.dev` and `.env` to `.env`
2. Update the `TELEGRAM_TOKEN` in both files with your bot token

## Local Development

### Using Go directly

```bash
# Run in development mode
go run cmd/bot/main.go
```

### Using Docker

```bash
# Build and run using Docker Compose
docker-compose up --build
```

## Building

```bash
# Build the binary
go build -o bot ./cmd/bot

# Build Docker image
docker build -t tg-bot .
```

## Project Structure

```
.
├── cmd/
│   └── bot/
│       └── main.go         # Application entry point
├── internal/
│   └── bot/
│       └── bot.go         # Bot implementation
├── configs/
│   └── config.go          # Configuration management
├── .env                   # Production environment variables
├── .env.dev              # Development environment variables
├── Dockerfile            # Docker build instructions
├── docker-compose.yml    # Docker Compose configuration
└── README.md            # This file
```

## Features

- Basic message handling
- Environment-based configuration
- Docker support
- Development and production modes 