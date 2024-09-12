package domain

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	FirstName    string    `json:"first_name"`
	SecondName   string    `json:"second_name"`
	Birthdate    time.Time `json:"birthdate"`
	Biography    string    `json:"biography"`
	City         string    `json:"city"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Не экспортируем пароль через API
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Token string

func (u *User) SetPassword(password string) {
	u.PasswordHash = HashPassword(password)
}
func (u *User) CheckPassword(password string) bool {
	return u.PasswordHash == HashPassword(password)
}

func HashPassword(password string) string {
	h := md5.New()
	io.WriteString(h, password)
	hash := h.Sum(nil)
	return fmt.Sprintf("%x", hash)
}
