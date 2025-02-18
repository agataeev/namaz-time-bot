package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"namaz-time-bot/internal/db"
	"strings"
)

var userCities = make(map[int64]string) // Хранение выбранного города

func HandleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	if strings.HasPrefix(data, "mark_") {
		prayer := strings.TrimPrefix(data, "mark_")
		err := db.SavePrayer(chatID, prayer)
		if err != nil {
			return
		}
		sendMessage(chatID, fmt.Sprintf("✅ Вы отметили выполнение намаза: %s", prayer))
	} else if strings.HasPrefix(data, "city_") {
		city := strings.TrimPrefix(data, "city_")
		err := db.SaveUser(chatID, city)
		if err != nil {
			return
		}
		sendMessage(chatID, fmt.Sprintf("🌍 Вы выбрали город: %s", city))
	} else {
		sendMessage(chatID, "Неизвестное действие")
	}

	botAPI.Request(tgbotapi.NewCallback(callback.ID, "Действие выполнено"))
}
