package bot

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func New(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api: api,
	}, nil
}

func (b *Bot) Run() error {
	log.Printf("Authorized on account %s", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			text := strings.ToLower(strings.TrimSpace(update.Message.Text))

			var response string
			switch text {
			case "привет":
				response = "Привет! 👋"
			default:
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			if _, err := b.api.Send(msg); err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}

	return nil
}
