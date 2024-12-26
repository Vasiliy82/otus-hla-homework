package apperrors

import "fmt"

type ErrorType string

const (
	ClientError        ErrorType = "client_error"
	ValidationError    ErrorType = "validation_error"
	AuthorizationError ErrorType = "authorization_error"
	ServiceError       ErrorType = "service_error"
	RemoteServiceError ErrorType = "remote_service_error"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}
