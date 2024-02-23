package model

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int
	UserID    int
	FirstName string
	LastName  string
	JoinAt    time.Time
	IsActive  bool
}

type UserManager struct {
	db *sql.DB
}

func NewUserManager(db *sql.DB) *UserManager {
	return &UserManager{db: db}
}

func (um *UserManager) AddUser(user User) (int64, error) {
	result, err := um.db.Exec("INSERT INTO users (user_id, first_name, last_name, join_at, is_active) VALUES (?, ?, ?, ?, ?)",
		user.UserID,
		user.FirstName,
		user.LastName,
		user.JoinAt,
		true,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (um *UserManager) GetByUserID(userId int64) (User, error) {
	var user User
	sql := `SELECT 
		id, user_id, first_name, last_name, is_active 
	FROM 
		users 
	WHERE
		user_id = ? 
	LIMIT 1`
	result, err := um.db.Query(sql, userId)
	if err != nil {
		return user, err
	}
	defer result.Close()

	// Periksa apakah ada baris yang tersedia
	if result.Next() {
		// Jika ada baris tersedia, pindahkan nilai-nilai ke variabel user
		if err := result.Scan(&user.ID, &user.UserID, &user.FirstName, &user.LastName, &user.IsActive); err != nil {
			println(err.Error())
			return user, err
		}

		// Cetak nilai-nilai untuk memastikan bahwa mereka telah diisi
		println(user.ID, user.UserID, user.FirstName, user.LastName)
	} else {
		// Jika tidak ada baris yang tersedia, kembalikan error
		return user, errors.New("no user found with the given user ID")
	}

	return user, nil
}
