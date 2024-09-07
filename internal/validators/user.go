package validators

import (
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/internal/apperrors"
	"github.com/Vasiliy82/otus-hla-homework/internal/dto"
)

func ValidateRegisterUserRequest(request dto.RegisterUserRequest) error {
	if err := request.Validate(); err != nil {
		return apperrors.NewBadRequestError(fmt.Sprintf("validation failed: %s", err.Error()))
	}
	return nil
}

func ValidateLoginRequest(request dto.LoginRequest) error {
	if err := request.Validate(); err != nil {
		return apperrors.NewBadRequestError(fmt.Sprintf("validation failed: %s", err.Error()))
	}
	return nil
}

func ValidateUserId(id string) error {
	if id == "" {
		return apperrors.NewBadRequestError("id cannot be empty")
	}
	return nil
}
