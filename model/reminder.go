package model

import (
	"database/sql"
	"fmt"
	"time"
)

type Reminder struct {
	ID       int
	UserID   int
	Message  string
	DueTime  time.Time
	IsActive bool
}

type ReminderManager struct {
	db *sql.DB
}

func NewReminderManager(db *sql.DB) *ReminderManager {
	return &ReminderManager{db: db}
}

func (rm *ReminderManager) SelectAll() ([]Reminder, error) {
	rows, err := rm.db.Query("SELECT id, user_id, message, due_time, is_active FROM reminders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []Reminder
	for rows.Next() {
		var r Reminder
		err := rows.Scan(&r.ID, &r.UserID, &r.Message, &r.DueTime, &r.IsActive)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reminders, nil
}

func (rm *ReminderManager) CheckReminder() error {
	now := time.Now()
	rows, err := rm.db.Query("SELECT id, user_id, message, due_time FROM reminders WHERE is_active = ? AND due_time <= ?", true, now)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var r Reminder
		err := rows.Scan(&r.ID, &r.UserID, &r.Message, &r.DueTime)
		if err != nil {
			return err
		}
		// Implement logic to send reminder to user
		fmt.Printf("Reminder sent to user %d: %s\n", r.UserID, r.Message)
	}

	return nil
}

func (rm *ReminderManager) AddReminder(userID int, message string, dueTime time.Time) (int64, error) {
	result, err := rm.db.Exec("INSERT INTO reminders (user_id, message, due_time, is_active) VALUES (?, ?, ?, ?)", userID, message, dueTime, true)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (rm *ReminderManager) UpdateReminder(reminder Reminder) error {
	_, err := rm.db.Exec("UPDATE reminders SET user_id=?, message=?, due_time=?, is_active=? WHERE id=?", reminder.UserID, reminder.Message, reminder.DueTime, reminder.IsActive, reminder.ID)
	return err
}

func (rm *ReminderManager) DeleteReminder(id int) error {
	_, err := rm.db.Exec("DELETE FROM reminders WHERE id = ?", id)
	return err
}
