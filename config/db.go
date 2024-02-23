package config

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	// Execute DDL statements
	err = createTables(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTableUser(db *sql.DB) (err error) {
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT,
		first_name TEXT,
		last_name TEXT,
		join_at DATETIME,
		is_active BOOLEAN
	)`)
	return
}
func createTableReminder(db *sql.DB) (err error) {
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS reminders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		message TEXT,
		due_time DATETIME,
		is_active BOOLEAN
	)`)
	return
}

func createTables(db *sql.DB) error {
	if err := createTableUser(db); err != nil {
		return err
	}
	if err := createTableReminder(db); err != nil {
		return err
	}
	return nil
}
