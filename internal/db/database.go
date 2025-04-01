package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// InitDB инициализирует подключение к базе данных
func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Не задан DATABASE_URL")
	}

	var err error
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
}

// SaveUser сохраняет пользователя и его город
func SaveUser(chatID int64, city string) error {
	_, err := DB.Exec(context.Background(),
		"INSERT INTO users (chat_id, city) VALUES ($1, $2)",
		chatID, city)
	return err
}

// GetUserCity получает город пользователя
func GetUserCity(chatID int64) (string, error) {
	var city string
	err := DB.QueryRow(context.Background(),
		"SELECT city FROM users WHERE chat_id = $1", chatID).Scan(&city)
	return city, err
}

// SavePrayer отмечает выполнение намаза
func SavePrayer(chatID int64, prayer string) error {
	_, err := DB.Exec(context.Background(),
		"INSERT INTO prayers (chat_id, prayer_name) VALUES ($1, $2)", chatID, prayer)
	return err
}

// SaveReminder сохраняет напоминание
func SaveReminder(chatID int64, prayer string, time string) error {
	_, err := DB.Exec(context.Background(),
		"INSERT INTO reminders (chat_id, prayer_name, reminder_time) VALUES ($1, $2, $3)",
		chatID, prayer, time)
	return err
}

// GetReminders получает напоминания, которые нужно отправить
func GetReminders(currentTime string) ([]struct {
	ChatID     int64
	PrayerName string
}, error) {
	rows, err := DB.Query(context.Background(),
		"SELECT chat_id, prayer_name FROM reminders WHERE reminder_time = $1", currentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []struct {
		ChatID     int64
		PrayerName string
	}
	for rows.Next() {
		var r struct {
			ChatID     int64
			PrayerName string
		}
		if err := rows.Scan(&r.ChatID, &r.PrayerName); err != nil {
			return nil, err
		}
		reminders = append(reminders, r)
	}
	return reminders, nil
}

// SavePrayerTimes сохраняет расписание намазов для пользователя
func SavePrayerTimes(chatID int64, city string, fajr, dhuhr, asr, maghrib, isha string) error {
	_, err := DB.Exec(context.Background(),
		`INSERT INTO prayer_times (chat_id, city, fajr, dhuhr, asr, maghrib, isha, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
         ON CONFLICT (chat_id) 
         DO UPDATE SET city = $2, fajr = $3, dhuhr = $4, asr = $5, maghrib = $6, isha = $7, updated_at = NOW()`,
		chatID, city, fajr, dhuhr, asr, maghrib, isha)
	return err
}

// GetPrayerTimes возвращает время намазов для конкретного пользователя
func GetPrayerTimes(chatID int64) (map[string]string, error) {
	row := DB.QueryRow(context.Background(),
		"SELECT fajr, dhuhr, asr, maghrib, isha FROM prayer_times WHERE chat_id = $1", chatID)

	var fajr, dhuhr, asr, maghrib, isha string
	err := row.Scan(&fajr, &dhuhr, &asr, &maghrib, &isha)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"Фаджр":  fajr,
		"Зухр":   dhuhr,
		"Аср":    asr,
		"Магриб": maghrib,
		"Иша":    isha,
	}, nil
}

// GetUsersWithPrayerTimes получает пользователей с установленными намазами
func GetUsersWithPrayerTimes() ([]struct {
	ChatID int64
	City   string
}, error) {
	rows, err := DB.Query(context.Background(),
		"SELECT chat_id, city FROM users WHERE city IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []struct {
		ChatID int64
		City   string
	}
	for rows.Next() {
		var user struct {
			ChatID int64
			City   string
		}
		if err := rows.Scan(&user.ChatID, &user.City); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
