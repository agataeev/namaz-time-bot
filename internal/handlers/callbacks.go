package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"namaz-time-bot/intern
	"namaz-time-bot/internal/api"
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
			sendMessage(chatID, "❌ Ошибка сохранения города")
			return

		
		// Автоматически устанавливаем время намазов и напоминания

		
		sendMessage(chatID, fmt.Sprintf("🌍 Город установлен: %s\n🔔 Напоминания о намазах активированы!\n\n📋 Используйте /prayer_times чтобы увидеть расписание", city))
	} else {
		sendMessage(chatID, "Неизвестное действие")
	}

	botAPI.Request(tgbotapi.NewCallback(callback.ID, "Действие выполнено"))
}

// setupPrayerTimesAndReminders автоматически устанавливает время намазов и напоминания для пользователя
func setupPrayerTimesAndReminders(chatID int64, city string) {
	// Получаем время намазов из API
	times, err := api.GetPrayerTimes(city, "Kazakhstan")
	if err != nil {
		log.Printf("Ошибка получения времени намазов для пользователя %d: %v", chatID, err)
		sendMessage(chatID, "⚠️ Не удалось получить время намазов. Попробуйте позже.")
		return
	}

	// Сохраняем время намазов в БД
	err = db.SavePrayerTimes(chatID, city, times.Fajr, times.Dhuhr, times.Asr, times.Maghrib, times.Isha)
	if err != nil {
		log.Printf("Ошибка сохранения времени намазов для пользователя %d: %v", chatID, err)
		sendMessage(chatID, "⚠️ Не удалось сохранить время намазов.")
		return
	}

	log.Printf("Время намазов успешно установлено для пользователя %d в городе %s", chatID, city)
}
