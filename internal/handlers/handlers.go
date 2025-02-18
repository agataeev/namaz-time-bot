package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var botAPI *tgbotapi.BotAPI

// InitHandlers инициализирует обработчики с ботом
func InitHandlers(bot *tgbotapi.BotAPI) {
	botAPI = bot
}

// sendMessage отправляет сообщение через бота
func sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	botAPI.Send(msg)
}
