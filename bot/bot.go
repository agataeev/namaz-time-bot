package bot

import (
	"log"
	"namaz-time-bot/config"
	"namaz-time-bot/internal/db"
	"namaz-time-bot/internal/handlers"
	"namaz-time-bot/internal/jobs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func StartBot() {
	var err error
	Bot, err = tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	db.InitDB()                // Подключаем БД
	handlers.InitHandlers(Bot) // Передаём бота в обработчики
	jobs.InitJobs(Bot)         // Передаём бота в фоновые задачи
	jobs.StartReminderJob()    // Запускаем напоминания

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := Bot.GetUpdatesChan(u)

	//add log
	log.Printf("Authorized on account %s", Bot.Self.UserName)
	log.Printf("Bot started...🚀🚀🚀")
	for update := range updates {
		if update.Message != nil {
			handlers.HandleCommand(update.Message)
		} else if update.CallbackQuery != nil {
			handlers.HandleCallback(update.CallbackQuery)
		}
	}
}
