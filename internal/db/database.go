package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
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

	// Тестируем подключение
	if err := TestConnection(); err != nil {
		log.Fatal("Ошибка тестирования подключения к БД:", err)
	}

	log.Println("✅ Подключение к базе данных успешно установлено")
}

// TestConnection тестирует подключение к базе данных
func TestConnection() error {
	if DB == nil {
		return fmt.Errorf("подключение к БД не инициализировано")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Тестируем подключение простым запросом
	var result int
	err := DB.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("ошибка выполнения тестового запроса: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("неожиданный результат тестового запроса: %d", result)
	}

	return nil
}

// GetDatabaseInfo возвращает информацию о базе данных
func GetDatabaseInfo() (map[string]interface{}, error) {
	if DB == nil {
		return nil, fmt.Errorf("подключение к БД не инициализировано")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	info := make(map[string]interface{})

	// Получаем версию PostgreSQL
	var version string
	err := DB.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения версии БД: %w", err)
	}
	info["version"] = version

	// Получаем текущую базу данных
	var dbName string
	err = DB.QueryRow(ctx, "SELECT current_database()").Scan(&dbName)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения имени БД: %w", err)
	}
	info["database"] = dbName

	// Получаем текущего пользователя
	var user string
	err = DB.QueryRow(ctx, "SELECT current_user").Scan(&user)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя БД: %w", err)
	}
	info["user"] = user

	// Получаем статистику подключений
	stats := DB.Stat()
	info["connections"] = map[string]interface{}{
		"total":        stats.TotalConns(),
		"idle":         stats.IdleConns(),
		"acquired":     stats.AcquiredConns(),
		"constructing": stats.ConstructingConns(),
	}

	return info, nil
}

// CheckTablesExist проверяет существование необходимых таблиц
func CheckTablesExist() (map[string]bool, error) {
	if DB == nil {
		return nil, fmt.Errorf("подключение к БД не инициализировано")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tables := []string{"users", "prayer_times", "prayers", "reminders", "schema_migrations"}
	result := make(map[string]bool)

	for _, table := range tables {
		var exists bool
		query := `SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)`

		err := DB.QueryRow(ctx, query, table).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("ошибка проверки таблицы %s: %w", table, err)
		}
		result[table] = exists
	}

	return result, nil
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
	location := time.FixedZone("UTC+5", 5*60*60)
	completedAt := time.Now().In(location).Format("2006-01-02 15:04:05.000000")
	_, err := DB.Exec(context.Background(),
		"INSERT INTO prayers (chat_id, prayer_name, completed_at) VALUES ($1, $2, $3)", chatID, prayer, completedAt)
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
	// Сохраняем расписание намазов в БД

	var exists bool
	err := DB.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM prayer_times WHERE chat_id = $1)`,
		chatID,
	).Scan(&exists)
	if err != nil {
		// обработка ошибки
		log.Println("ошибка при проверке существования:", err)
		return err
	}

	if exists {
		// Обновляем
		_, err = DB.Exec(context.Background(),
			`UPDATE prayer_times 
		 SET city = $2, fajr = $3, dhuhr = $4, asr = $5, maghrib = $6, isha = $7, updated_at = NOW()
		 WHERE chat_id = $1`,
			chatID, city, fajr, dhuhr, asr, maghrib, isha,
		)
	} else {
		// Вставляем
		_, err = DB.Exec(context.Background(),
			`INSERT INTO prayer_times (chat_id, city, fajr, dhuhr, asr, maghrib, isha, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`,
			chatID, city, fajr, dhuhr, asr, maghrib, isha,
		)
	}

	if err != nil {
		log.Println("ошибка при вставке/обновлении:", err)
		return err
	}

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
