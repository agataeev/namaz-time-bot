package jobs

import (
	"log"
	"namaz-time-bot/internal/db"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var botAPI *tgbotapi.BotAPI

// InitJobs инициализирует фоновые задачи с ботом
func InitJobs(bot *tgbotapi.BotAPI) {
	botAPI = bot
}

// StartReminderJob запускает проверку напоминаний каждую минуту
func StartReminderJob() {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for range ticker.C {
			currentTime := time.Now().Format("15:04")
			reminders, err := db.GetReminders(currentTime)
			if err != nil {
				log.Println("Ошибка получения напоминаний:", err)
				continue
			}

			for _, r := range reminders {
				msg := tgbotapi.NewMessage(r.ChatID, "🔔 Время намаза: "+r.PrayerName)
				botAPI.Send(msg)
			}
		}
	}()
}
