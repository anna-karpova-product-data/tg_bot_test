package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	IsDev         bool
}

func LoadConfig() (*Config, error) {
	envFile := ".env"
	if os.Getenv("ENV") == "dev" {
		envFile = ".env.dev"
	}

	if err := godotenv.Load(envFile); err != nil {
		return nil, err
	}

	isDev, _ := strconv.ParseBool(os.Getenv("DEV"))

	return &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		IsDev:         isDev,
	}, nil
}
