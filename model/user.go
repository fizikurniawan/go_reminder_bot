package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID       int
	UserID   int
	Name     string
	JoinAt   time.Time
	IsActive bool
}

type UserManager struct {
	db *sql.DB
}

func NewUserManager(db *sql.DB) *UserManager {
	return &UserManager{db: db}
}

func (um *UserManager) AddUser(userID string, name string, jointAt time.Time) (int64, error) {
	result, err := um.db.Exec("INSERT INTO users (user_id, name, join_at, is_active) VALUES (?, ?, ?, ?)", userID, name, jointAt, true)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (um *UserManager) GetByUserID(userId string) (User, error) {
	var user User
	result, err := um.db.Query("SELECT * FROM users WHERE user_id = ?", userId)
	if err != nil {
		return user, err
	}
	defer result.Close()
	result.Scan(&user.ID, &user.UserID, &user.Name)

	return user, nil
}
