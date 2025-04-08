package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"namaz-time-bot/internal/api"
	"namaz-time-bot/internal/db"
)

func HandleCommand(msg *tgbotapi.Message) {
	switch msg.Text {
	case "/start":
		sendWelcomeMessage(msg.Chat.ID)
	case "/mark":
		sendPrayerButtons(msg.Chat.ID)
	case "/prayer_times":
		city, err := db.GetUserCity(msg.Chat.ID)
		if err != nil {
			sendMessage(msg.Chat.ID, "🌍 Вы не выбрали город! Используйте /set_city.")
			return
		}

		times, err := api.GetPrayerTimes(city, "Kazakhstan")
		if err != nil {
			sendMessage(msg.Chat.ID, "Ошибка при получении времени намаза.")
			return
		}

		response := fmt.Sprintf(
			"Время намазов в %s:\n🌅 Фаджр: %s\n🏙️ Зухр: %s\n🌇 Аср: %s\n🌆 Магриб: %s\n🌃 Иша: %s",
			city, times.Fajr, times.Dhuhr, times.Asr, times.Maghrib, times.Isha,
		)
		sendMessage(msg.Chat.ID, response)
	case "/set_city":
		sendCityButtons(msg.Chat.ID)
	case "/set_reminders":
		setReminders(msg.Chat.ID)
	case "/help":
		sendHelpMessage(msg.Chat.ID)
	case "/set_prayer_times":
		setPrayerTimes(msg.Chat.ID)
	default:
		sendMessage(msg.Chat.ID, "Неизвестная команда. Используйте /start.")
	}
}

// setPrayerTimes устанавливает время намазов для пользователя
func setPrayerTimes(chatID int64) {
	city, err := db.GetUserCity(chatID)
	if err != nil {
		log.Println(err)
		sendMessage(chatID, "🌍 Вы не выбрали город! Используйте /set_city.")
		return
	}

	times, err := api.GetPrayerTimes(city, "Kazakhstan")
	if err != nil {
		log.Println(err)
		sendMessage(chatID, "Ошибка при получении времени намаза.")
		return
	}

	// Сохраняем в БД
	err = db.SavePrayerTimes(chatID, city, times.Fajr, times.Dhuhr, times.Asr, times.Maghrib, times.Isha)
	if err != nil {
		log.Println(err)
		sendMessage(chatID, "❌ Ошибка сохранения времени намаза.")
		return
	}

	sendMessage(chatID, "✅ Время намазов обновлено!")
}

// sendWelcomeMessage отправляет приветственное сообщение
func sendWelcomeMessage(chatID int64) {
	text := "🌟 Добро пожаловать в бота для напоминаний о намазе!\n\n" +
		"Используйте /help, чтобы увидеть список доступных команд."

	sendMessage(chatID, text)
}

// sendHelpMessage отправляет список команд
func sendHelpMessage(chatID int64) {
	text := "📌 *Список команд:*\n\n" +
		"🔹 /start - Начать работу с ботом\n" +
		"🔹 /help - Показать список команд\n" +
		"🔹 /set_city - Установить город для определения времени намаза\n" +
		//"🔹 /set_reminders - Включить напоминания о намазах\n" +
		//"🔹 /disable_reminders - Отключить напоминания\n" +
		"🔹 /prayer_times - Показать текущее расписание намазов\n" +
		"🔹 /mark - Отметить выполнение намаза\n" +
		"🔹 /set_prayer_times - Установить время намазов\n"

	sendMessage(chatID, text)
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
			tgbotapi.NewInlineKeyboardButtonData("🏙️ Алматы", "city_Almaty"),
			tgbotapi.NewInlineKeyboardButtonData("🏙️ Астана", "city_Astana"),
		),
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
