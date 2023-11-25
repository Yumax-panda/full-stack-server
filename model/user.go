package model

import (
	"database/sql"

	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string
	Comment  sql.NullString
}

func GetUser() ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}
