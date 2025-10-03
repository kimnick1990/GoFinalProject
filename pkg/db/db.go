package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // Подключаем драйвер SQLite
)

var db *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255),
    comment TEXT,
    repeat VARCHAR(128)
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
`

func Init(dbFile string) error {
	var err error
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	_, err = db.Exec(schema)
	if err != nil {
		log.Println("Ошибка при создании таблицы:", err)
		return err
	}
	return nil
}

func Close() error {
	return db.Close()
}
