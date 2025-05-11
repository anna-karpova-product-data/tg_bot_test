# Project Architecture

## Overview
This project is a Telegram bot that processes audio messages, transcribes them using the Deepgram API, analyzes the user's mood, and responds accordingly. The bot is written in Go and uses ffmpeg for audio conversion.

## Main Components

- **cmd/bot/main.go**: Entry point for the bot application.
- **internal/bot/bot.go**: Main bot logic, handles user interactions, mood analysis, and message processing.
- **internal/speech/deepgram.go**: Deepgram API client for audio transcription. Configured for Russian language and optimal recognition parameters.
- **internal/speech/audio.go**: Audio conversion utilities (OGG to WAV) using ffmpeg.
- **configs/config.go**: Loads configuration and environment variables.
- **.env / .env.dev**: Environment configuration for dev/prod modes.

## Data Flow
1. User sends a voice message to the bot.
2. Bot downloads and saves the OGG file.
3. Audio is converted to WAV using ffmpeg.
4. WAV file is sent to Deepgram API for transcription (with Russian language and optimal parameters).
5. Transcribed text is analyzed for mood.
6. Bot responds with a message or exercise suggestions based on detected mood.

## Logging
Detailed logging is implemented for each step of audio processing and transcription to aid in debugging and support.

## Dependencies
- Go
- ffmpeg (must be available in PATH)
- Deepgram API
- go-telegram-bot-api

---

For more details on tasks and changes, see `design/tasks.md`. 