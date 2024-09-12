package dto

// Структура для запроса на вход
type LoginRequest struct {
	Username string `json:"username" validate:"required,email"` // Username должен быть email и обязательным
	Password string `json:"password" validate:"required,min=6"` // Пароль обязательный и минимум 6 символов
}

// Структура для регистрации пользователя
type RegisterUserRequest struct {
	FirstName  string `json:"first_name" validate:"required"`                    // Имя обязательно
	SecondName string `json:"second_name" validate:"required"`                   // Фамилия обязательна
	Birthdate  string `json:"birthdate" validate:"required,datetime=2006-01-02"` // Дата рождения обязательна
	Biography  string `json:"biography" validate:"omitempty"`                    // Биография необязательная
	City       string `json:"city" validate:"required"`                          // Город обязателен
	Username   string `json:"username" validate:"required,email"`                // Username должен быть email и обязательным
	Password   string `json:"password" validate:"required,min=6"`                // Пароль обязателен и минимум 6 символов
}

type LoginResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}
