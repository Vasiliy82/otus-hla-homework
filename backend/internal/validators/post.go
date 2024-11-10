package validators

import (
	"github.com/Vasiliy82/otus-hla-homework/internal/apperrors"
	"github.com/Vasiliy82/otus-hla-homework/internal/dto"
)

// Валидация запроса на создание поста
func ValidateCreateOrUpdatePostRequest(request dto.CreateOrUpdatePostRequest) error {
	if err := validate.Struct(request); err != nil {
		return apperrors.NewValidationError("Validation error", err, tran)
	}
	return nil
}

// Валидация идентификатора поста
func ValidatePostId(id string) error {
	if id == "" {
		return apperrors.NewBadRequestError("id cannot be empty")
	}
	return nil
}
