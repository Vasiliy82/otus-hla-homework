package validators

import (
	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/dto"
)

func ValidateRegisterUserRequest(request dto.RegisterUserRequest) error {
	if request.FirstName == "" || request.SecondName == "" || request.Birthdate == "" || request.Password == "" || request.Username == "" {
		return domain.ErrBadParamInput
	}
	return nil
}

func ValidateLoginRequest(request dto.LoginRequest) error {
	if request.Username == "" || request.Password == "" {
		return domain.ErrBadParamInput
	}
	return nil
}

func ValidateUserId(id string) error {
	if id == "" {
		return domain.ErrBadParamInput
	}
	return nil
}
