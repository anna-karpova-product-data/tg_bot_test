package main

import (
	"log"

	"tg_bot/configs"
	"tg_bot/internal/bot"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	b, err := bot.New(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	if err := b.Run(); err != nil {
		log.Fatalf("Error running bot: %v", err)
	}
}
