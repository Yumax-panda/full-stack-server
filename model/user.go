package model

import (
	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID      `json:"id"`
	Admin    bool           `json:"admin"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Comment  sql.NullString `json:"comment"`
	Skills   []*Skill       `json:"skills"`
}

func GetUsers() ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}

func CreateUser(name string, email string, password string, comment string) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	user := User{
		ID:       uuid.New(),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Comment:  sql.NullString{String: comment, Valid: true},
	}
	err = db.Create(&user).Error
	return user, err
}

func GetUserByEmail(email string) (User, error) {
	var user User
	err := db.Where("email = ?", email).First(&user).Error
	return user, err
}

func GetUserByID(id string) (User, error) {
	var user User
	err := db.Where("id = ?", id).First(&user).Error
	return user, err
}
