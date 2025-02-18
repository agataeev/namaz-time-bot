package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "prayers.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable()
}

func createTable() {
	query := `CREATE TABLE IF NOT EXISTS prayers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		prayer_time TEXT,
		marked BOOLEAN
	)`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
