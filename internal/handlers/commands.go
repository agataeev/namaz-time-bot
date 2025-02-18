package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"namaz-time-bot/internal/api"
	"namaz-time-bot/internal/db"
)

func HandleCommand(msg *tgbotapi.Message) {
	switch msg.Text {
	case "/start":
		sendMessage(msg.Chat.ID, "Привет! Я бот для напоминаний о намазе.")
	case "/mark":
		sendPrayerButtons(msg.Chat.ID)
	case "/prayer_times":
		city, err := db.GetUserCity(msg.Chat.ID)
		if err != nil {
			sendMessage(msg.Chat.ID, "🌍 Вы не выбрали город! Используйте /set_city.")
			return
		}

		times, err := api.GetPrayerTimes(city, "Russia")
		if err != nil {
			sendMessage(msg.Chat.ID, "Ошибка при получении времени намаза.")
			return
		}

		response := fmt.Sprintf(
			"Время намазов в %s:\n🌅 Фаджр: %s\n🏙️ Зухр: %s\n🌇 Аср: %s\n🌆 Магриб: %s\n🌃 Иша: %s",
			city, times.Fajr, times.Dhuhr, times.Asr, times.Maghrib, times.Isha,
		)
		sendMessage(msg.Chat.ID, response)
	case "/set_reminders":
		setReminders(msg.Chat.ID)
	default:
		sendMessage(msg.Chat.ID, "Неизвестная команда. Используйте /start.")
	}
}

func setReminders(chatID int64) {
	city, err := db.GetUserCity(chatID)
	if err != nil {
		sendMessage(chatID, "🌍 Вы не выбрали город! Используйте /set_city.")
		return
	}

	times, err := api.GetPrayerTimes(city, "Russia")
	if err != nil {
		sendMessage(chatID, "Ошибка при получении времени намаза.")
		return
	}

	// Сохраняем напоминания в БД
	err = db.SaveReminder(chatID, "Фаджр", times.Fajr)
	if err != nil {
		return
	}
	err = db.SaveReminder(chatID, "Зухр", times.Dhuhr)
	if err != nil {
		return
	}
	err = db.SaveReminder(chatID, "Аср", times.Asr)
	if err != nil {
		return
	}
	err = db.SaveReminder(chatID, "Магриб", times.Maghrib)
	if err != nil {
		return
	}
	err = db.SaveReminder(chatID, "Иша", times.Isha)
	if err != nil {
		return
	}

	sendMessage(chatID, "🔔 Напоминания о намазах установлены!")
}

func sendPrayerButtons(chatID int64) {
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Фаджр", "mark_fajr"),
			tgbotapi.NewInlineKeyboardButtonData("✅ Зухр", "mark_dhuhr"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Аср", "mark_asr"),
			tgbotapi.NewInlineKeyboardButtonData("✅ Магриб", "mark_maghrib"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Иша", "mark_isha"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, "Выберите намаз, который вы выполнили:")
	msg.ReplyMarkup = buttons
	_, err := botAPI.Send(msg)
	if err != nil {
		return
	}
}

func sendCityButtons(chatID int64) {
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏙️ Москва", "city_Moscow"),
			tgbotapi.NewInlineKeyboardButtonData("🏙️ Казань", "city_Kazan"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏙️ Мекка", "city_Mecca"),
			tgbotapi.NewInlineKeyboardButtonData("🏙️ Медина", "city_Medina"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, "Выберите ваш город:")
	msg.ReplyMarkup = buttons
	botAPI.Send(msg)
}

//func sendMessage(chatID int64, text string) {
//	msg := tgbotapi.NewMessage(chatID, text)
//	bot.Bot.Send(msg)
//}
