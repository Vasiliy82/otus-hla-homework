package dto

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ruTranslations "github.com/go-playground/validator/v10/translations/ru"
)

// Переменные для валидатора и переводчика
var (
	validate *validator.Validate
	trans    ut.Translator
)

// Инициализация валидатора и переводчика
func init() {
	validate = validator.New()

	// Установка переводчика для русского языка
	ruLocale := ru.New()
	uni := ut.New(ruLocale, ruLocale)

	// Получение переводчика для русского языка
	trans, _ = uni.GetTranslator("ru")

	// Регистрация русских переводов ошибок валидации
	ruTranslations.RegisterDefaultTranslations(validate, trans)
}

// Структура для запроса на вход
type LoginRequest struct {
	Username string `json:"username" validate:"required,email" label:"Имя пользователя"` // Username должен быть email и обязательным
	Password string `json:"password" validate:"required,min=6" label:"Пароль"`           // Пароль обязательный и минимум 6 символов
}

// Структура для регистрации пользователя
type RegisterUserRequest struct {
	FirstName  string `json:"first_name" validate:"required" label:"Имя"`                  // Имя обязательно
	SecondName string `json:"second_name" validate:"required" label:"Фамилия"`             // Фамилия обязательна
	Birthdate  string `json:"birthdate" validate:"required" label:"Дата рождения"`         // Дата рождения обязательна
	Biography  string `json:"biography" validate:"omitempty" label:"Биография"`            // Биография необязательная
	City       string `json:"city" validate:"required" label:"Город"`                      // Город обязателен
	Username   string `json:"username" validate:"required,email" label:"Имя пользователя"` // Username должен быть email и обязательным
	Password   string `json:"password" validate:"required,min=6" label:"Пароль"`           // Пароль обязателен и минимум 6 символов
}

// CustomValidationError представляет переведенную ошибку валидации с дружественным именем поля
type CustomValidationError struct {
	Field   string
	Message string
}

// translateValidationErrors создает список кастомных ошибок валидации с переводом и заменой технических имен полей на дружественные
func translateValidationErrors(err error, s interface{}) []CustomValidationError {
	var customErrors []CustomValidationError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, validationErr := range validationErrs {
			// Получаем отражение структуры, чтобы извлечь тег label
			val := reflect.ValueOf(s)
			field, _ := val.Type().FieldByName(validationErr.StructField())

			// Извлекаем значение тега label
			label := field.Tag.Get("label")
			if label == "" {
				label = validationErr.Field() // если не найден тег label, используем имя поля
			}

			// Получаем текущий перевод ошибки
			translatedError := validationErr.Translate(trans)

			// Заменяем техническое название поля на дружественное имя в сообщении об ошибке
			translatedError = strings.Replace(translatedError, validationErr.Field(), label, 1)

			// Добавляем в список кастомную ошибку с переведенным сообщением
			customErrors = append(customErrors, CustomValidationError{
				Field:   validationErr.Field(),
				Message: translatedError,
			})
		}
	}

	return customErrors
}

// Валидация структуры LoginRequest с возвратом ошибок
func (l *LoginRequest) Validate() error {
	err := validate.Struct(l)
	if err != nil {
		translatedErrors := translateValidationErrors(err, l)
		if len(translatedErrors) > 0 {
			// Возвращаем ошибку в виде списка переведенных ошибок
			return &ValidationErrors{Errors: translatedErrors}
		}
	}
	return nil
}

// Валидация структуры RegisterUserRequest с возвратом ошибок
func (r *RegisterUserRequest) Validate() error {
	err := validate.Struct(r)
	if err != nil {
		translatedErrors := translateValidationErrors(err, r)
		if len(translatedErrors) > 0 {
			// Возвращаем ошибку в виде списка переведенных ошибок
			return &ValidationErrors{Errors: translatedErrors}
		}
	}
	return nil
}

// ValidationErrors представляет список переведенных ошибок
type ValidationErrors struct {
	Errors []CustomValidationError
}

// Error реализует интерфейс error для структуры ValidationErrors
func (ve *ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, ", ")
}
