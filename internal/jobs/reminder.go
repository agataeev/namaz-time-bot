package jobs

import (
	"fmt"
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

			// Получаем пользователей с установленными намазами
			users, err := db.GetUsersWithPrayerTimes()
			if err != nil {
				log.Println("Ошибка получения пользователей:", err)
				continue
			}

			for _, user := range users {
				times, err := db.GetPrayerTimes(user.ChatID)
				if err != nil {
					continue
				}

				for prayer, time := range times {
					if time[:5] == currentTime {
						msg := tgbotapi.NewMessage(user.ChatID, fmt.Sprintf("🔔 Время %s! 🙏", prayer))
						botAPI.Send(msg)
					}
				}
			}
		}
	}()
}
