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
	location := time.FixedZone("UTC+5", 5*60*60)

	go func() {
		for range ticker.C {
			currentTime := time.Now().In(location).Format("15:04")
			log.Printf("Проверка напоминаний на время: %s", currentTime)

			// Получаем пользователей с установленными намазами
			users, err := db.GetUsersWithPrayerTimes()
			if err != nil {
				log.Println("Ошибка получения пользователей:", err)
				continue
			}

			log.Printf("Найдено пользователей с установленными намазами: %d", len(users))

			for _, user := range users {
				times, err := db.GetPrayerTimes(user.ChatID)
				if err != nil {
					log.Printf("Ошибка получения времени намазов для пользователя %d: %v", user.ChatID, err)
					continue
				}

				for prayer, prayerTime := range times {
					if prayerTime == currentTime {
						log.Printf("Отправка напоминания пользователю %d о намазе %s в %s", user.ChatID, prayer, prayerTime)

						// Создаем сообщение с кнопкой для отметки
						msg := tgbotapi.NewMessage(user.ChatID, fmt.Sprintf("🔔 Время %s! 🙏\n\n⏰ %s\n\n✅ Не забудьте отметить выполнение намаза", prayer, prayerTime))

						// Добавляем кнопку для быстрой отметки
						var buttonData string
						switch prayer {
						case "Фаджр":
							buttonData = "mark_fajr"
						case "Зухр":
							buttonData = "mark_dhuhr"
						case "Аср":
							buttonData = "mark_asr"
						case "Магриб":
							buttonData = "mark_maghrib"
						case "Иша":
							buttonData = "mark_isha"
						}

						if buttonData != "" {
							button := tgbotapi.NewInlineKeyboardMarkup(
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData("✅ Отметить выполнение", buttonData),
								),
							)
							msg.ReplyMarkup = button
						}

						_, err := botAPI.Send(msg)
						if err != nil {
							log.Printf("Ошибка отправки напоминания пользователю %d: %v", user.ChatID, err)
						}
					}
				}
			}
		}
	}()
}
