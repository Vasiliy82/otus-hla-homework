package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError базовая структура для всех ошибок приложения.
type AppError struct {
	Code    int    // HTTP код ошибки
	Message string // Сообщение для клиента
	Err     error  // Техническая ошибка (для логирования)
}

// Error реализует интерфейс error.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap позволяет использовать стандартные функции errors.Unwrap для работы с вложенными ошибками.
func (e *AppError) Unwrap() error {
	return e.Err
}

// New создаёт новый объект AppError с заданным кодом, сообщением и вложенной ошибкой.
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
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
		Err:     err,                     // Это сообщение будет в логах
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
