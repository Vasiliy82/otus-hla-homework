package apperrors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// AppError базовая структура для всех ошибок приложения.
type AppError struct {
	Code    int    `json:"-"`     // HTTP код ошибки
	Message string `json:"error"` // Сообщение для клиента
	err     error  // Техническая ошибка (для логирования)
}

// Error реализует интерфейс error.
func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.err)
	}
	return e.Message
}

// Unwrap позволяет использовать стандартные функции errors.Unwrap для работы с вложенными ошибками.
func (e *AppError) Unwrap() error {
	return e.err
}

// New создаёт новый объект AppError с заданным кодом, сообщением и вложенной ошибкой.
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		err:     err,
	}
}

// NewBadRequestError создаёт ошибку с кодом 400 (Bad Request) и заданным сообщением.
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

// NewUnauthorizedError создаёт ошибку с кодом 401 (Unauthorized) и заданным сообщением.
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

// NewNotFoundError создаёт ошибку с кодом 404 (Not Found) и заданным сообщением.
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

// NewConflictError создаёт ошибку с кодом 409 (Conflict) и заданным сообщением.
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
	}
}

// NewInternalServerError создаёт ошибку с кодом 500 (Internal Server Error) и вложенной технической ошибкой.
func NewInternalServerError(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error", // Это сообщение увидит пользователь
		err:     err,                     // Это сообщение будет в логах
	}
}

// Is позволяет использовать errors.Is для проверки типа ошибки.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As позволяет использовать errors.As для проверки типа ошибки.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// ValidationErrorField представляет ошибку валидации для конкретного поля
type ValidationErrorField struct {
	Field   string `json:"field"`   // Название поля (имя переменной в структуре)
	Message string `json:"message"` // Сообщение об ошибке для этого поля
}

// ValidationError представляет кастомную ошибку валидации
type ValidationError struct {
	Message string                 `json:"error"`
	Errors  []ValidationErrorField `json:"validation_errors"` // Список ошибок для каждого поля
}

// Error реализует интерфейс error
func (v *ValidationError) Error() string {
	var sb strings.Builder
	sb.WriteString(v.Message)
	for _, errField := range v.Errors {
		sb.WriteString(fmt.Sprintf(", %s: %s", errField.Field, errField.Message))
	}
	return sb.String()
}

// NewValidationError создает новую кастомную ошибку валидации на основе ошибок валидатора
func NewValidationError(message string, err error, tran ut.Translator) *ValidationError {

	var validationErrors []ValidationErrorField
	var validatorErr validator.ValidationErrors

	// Проверяем, что ошибка является типом validator.ValidationErrors
	if errors.As(err, &validatorErr) {
		for _, fieldErr := range validatorErr {
			// Добавляем каждую ошибку валидации поля в список кастомных ошибок
			validationErrors = append(validationErrors, ValidationErrorField{
				Field:   fieldErr.Field(),         // Имя поля
				Message: fieldErr.Translate(tran), // Сообщение об ошибке, можно локализовать через переводчик
			})
		}
	}

	return &ValidationError{
		Message: message,
		Errors:  validationErrors,
	}
}
