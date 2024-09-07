package dto

import (
	"github.com/go-playground/validator/v10"
)

// Структура для запроса на вход
type LoginRequest struct {
	Username string `json:"username" validate:"required,email"` // Username должен быть email и обязательным
	Password string `json:"password" validate:"required,min=6"` // Пароль обязательный и минимум 6 символов
}

// Структура для регистрации пользователя
type RegisterUserRequest struct {
	FirstName  string `json:"first_name" validate:"required"`     // Имя обязательно
	SecondName string `json:"second_name" validate:"required"`    // Фамилия обязательна
	Birthdate  string `json:"birthdate" validate:"required"`      // Дата рождения обязательна
	Biography  string `json:"biography" validate:"omitempty"`     // Биография необязательная
	City       string `json:"city" validate:"required"`           // Город обязателен
	Username   string `json:"username" validate:"required,email"` // Username должен быть email и обязательным
	Password   string `json:"password" validate:"required,min=6"` // Пароль обязателен и минимум 6 символов
}

// Валидация структуры LoginRequest
func (l *LoginRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(l)
}

// Валидация структуры RegisterUserRequest
func (r *RegisterUserRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
