package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var BotToken string

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	BotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
}
