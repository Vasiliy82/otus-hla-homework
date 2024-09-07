package validators

import (
	"github.com/Vasiliy82/otus-hla-homework/internal/apperrors"
	"github.com/Vasiliy82/otus-hla-homework/internal/dto"
)

func ValidateRegisterUserRequest(request dto.RegisterUserRequest) error {
	if request.FirstName == "" {
		return apperrors.NewBadRequestError("first name cannot be empty")
	} else if request.SecondName == "" {
		return apperrors.NewBadRequestError("second name cannot be empty")
	} else if request.Birthdate == "" {
		return apperrors.NewBadRequestError("birth date cannot be empty")
	} else if request.Password == "" {
		return apperrors.NewBadRequestError("password cannot be empty")
	} else if request.Username == "" {
		return apperrors.NewBadRequestError("username cannot be empty")
	}
	return nil
}

func ValidateLoginRequest(request dto.LoginRequest) error {
	if request.Username == "" {
		return apperrors.NewBadRequestError("username cannot be empty")
	} else if request.Password == "" {
		return apperrors.NewBadRequestError("password cannot be empty")
	}
	return nil
}

func ValidateUserId(id string) error {
	if id == "" {
		return apperrors.NewBadRequestError("id cannot be empty")
	}
	return nil
}
