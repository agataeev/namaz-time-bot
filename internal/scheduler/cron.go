package scheduler

import (
	"log"
	"namaz-time-bot/internal/api"
	"namaz-time-bot/internal/db"

	"github.com/robfig/cron/v3"
)

func StartScheduler() {
	c := cron.New()
	//once a day at 00:00
	_, err := c.AddFunc("0 0 * * *", func() { UpdatePrayerTimesJob() })
	if err != nil {
		log.Println("Ошибка добавления задачи в cron:", err)
		return
	}
	c.Start()
}

// UpdatePrayerTimesJob обновляет время намазов для всех пользователей
func UpdatePrayerTimesJob() {
	// Получаем пользователей с установленными намазами
	users, err := db.GetUsersWithPrayerTimes()
	if err != nil {
		log.Println("Ошибка получения пользователей:", err)
		return
	}
	for _, user := range users {
		city, err := db.GetUserCity(user.ChatID)
		if err != nil {
			log.Println(err)
			continue
		}

		times, err := api.GetPrayerTimes(city, "Kazakhstan")
		if err != nil {
			log.Println(err)
			continue
		}

		// Сохраняем в БД
		err = db.SavePrayerTimes(user.ChatID, city, times.Fajr, times.Dhuhr, times.Asr, times.Maghrib, times.Isha)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
